package main

import (
	"errors"
	"os"
	"strconv"

	"github.com/gocql/gocql"
)

type Trans struct {
	Bank          string
	Id            gocql.UUID
	PromizeBank   string
	PromizeId     gocql.UUID
	PromizeAmount string
	PromizeBlob   string
	FromZaddress  string
	FromBank      string
	FromAccount   string
	ToZaddress    string
	ToBank        string
	ToAccount     string
	Timestamp     int64
	Digsig        string
	Type          string
}

type User struct {
	Zaddress  string
	Bank      string
	Account   string
	Salt      string // random debit amount
	PublicKey string
	Verified  bool
	Zode      string
	Active    bool
	Timestamp int64
}

var Session *gocql.Session

func initCStarSession() {
	cluster := gocql.NewCluster(cassandraConfig.host)
	cluster.Port = port(cassandraConfig.port)
	cluster.Keyspace = cassandraConfig.keyspace
	cluster.Consistency = consistancy(cassandraConfig.consistancy)

	s, err := cluster.CreateSession()
	if err != nil {
		println("Error cassandra session:", err.Error())
		os.Exit(1)
	}
	Session = s
}

func clearCStarSession() {
	Session.Close()
}

func port(p string) int {
	i, err := strconv.Atoi(p)
	if err != nil {
		return 9042
	}

	return i
}

func consistancy(c string) gocql.Consistency {
	gc, err := gocql.MustParseConsistency(c)
	if err != nil {
		return gocql.All
	}

	return gc
}

func createTrans(trans *Trans) error {
	insert := func(table string) error {
		q := "INSERT INTO " + table + ` (
                bank,
                id,
                promize_bank,
                promize_id,
                promize_amount,
                promize_blob,
                from_zaddress,
                from_bank,
                from_account,
                to_zaddress,
                to_bank,
                to_account,
                timestamp,
                digsig,
                type
            )
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		err := Session.Query(q,
			trans.Bank,
			trans.Id,
			trans.PromizeBank,
			trans.PromizeId,
			trans.PromizeAmount,
			trans.PromizeBlob,
			trans.FromZaddress,
			trans.FromBank,
			trans.FromAccount,
			trans.ToZaddress,
			trans.ToBank,
			trans.ToAccount,
			trans.Timestamp,
			trans.Digsig,
			trans.Type).Exec()
		if err != nil {
			println(err.Error())
		}

		return err
	}

	// insert to both trans and transactions
	insert("trans")
	insert("transactions")

	return nil
}

func updateTrans(state string, bank string, id string) error {
	update := func(table string) error {
		q := "UPDATE " + table + " " +
			`
              SET state = ?
              WHERE
                bank = ?
                AND id = ?
             `
		err := Session.Query(q, state, bank, id).Exec()
		if err != nil {
			println(err.Error())
		}

		return err
	}

	// insert to both trans and transactions
	update("trans")
	update("transactions")

	return nil
}

func isDoubleSpend(from string, cid string) bool {
	// parse cid and get uuid
	uuid, err := gocql.ParseUUID(cid)
	if err != nil {
		println(err.Error)
		return true
	}

	m := map[string]interface{}{}
	q := `
        SELECT id FROM trans
            WHERE from_zaddress=?
            AND promize_id=?
        LIMIT 1
        ALLOW FILTERING
    `
	itr := Session.Query(q, from, uuid).Consistency(gocql.One).Iter()
	for itr.MapScan(m) {
		return true
	}

	return false
}

func createUser(user *User) error {
	q := `
        INSERT INTO users (
			zaddress,
            bank,
			account,
			public_key,
            salt,
            verified,
			zode,
			active,
			timestamp
        )
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
	err := Session.Query(q,
		user.Zaddress,
		user.Bank,
		user.Account,
		user.PublicKey,
		user.Salt,
		user.Verified,
		user.Zode,
		user.Active,
		user.Timestamp).Exec()

	if err != nil {
		println(err.Error())
	}

	return err
}

func getUser(zaddress string) (*User, error) {
	m := map[string]interface{}{}
	q := `
        SELECT account, salt, public_key, verified, zode, active
        FROM users
        WHERE zaddress = ?
        LIMIT 1
    `
	itr := Session.Query(q, zaddress).Consistency(gocql.One).Iter()
	for itr.MapScan(m) {
		user := &User{}
		user.Zaddress = zaddress
		user.Account = m["account"].(string)
		user.Salt = m["salt"].(string)
		user.PublicKey = m["public_key"].(string)
		user.Verified = m["verified"].(bool)
		user.Zode = m["zode"].(string)
		user.Active = m["active"].(bool)

		return user, nil
	}

	return nil, errors.New("Not found promize")
}

func setUserVerified(verified bool, zaddress string) error {
	q := `
		UPDATE users 
			SET verified = ? 
        WHERE zaddress = ?
        `
	err := Session.Query(q, verified, zaddress).Exec()
	if err != nil {
		println(err.Error())
		return err
	}

	return nil
}

func setUserActive(active bool, zaddress string) error {
	q := `
		UPDATE users 
			SET active = ? 
        WHERE zaddress = ?
        `
	err := Session.Query(q, active, zaddress).Exec()
	if err != nil {
		println(err.Error())
		return err
	}

	return nil
}

func uuid() gocql.UUID {
	return gocql.TimeUUID()
}

func uuidstr() string {
	return gocql.TimeUUID().String()
}

func cuuid(cid string) (gocql.UUID, error) {
	return gocql.ParseUUID(cid)
}

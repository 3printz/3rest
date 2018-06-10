package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func initHttpz() {
	// router
	r := mux.NewRouter()
	r.HandleFunc("/transactions", transactions).Methods("POST")

	// start server
	err := http.ListenAndServe(":7070", r)
	if err != nil {
		println(err.Error)
		os.Exit(1)
	}
}

func transactions(w http.ResponseWriter, r *http.Request) {
	// read body
	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	println(string(b))

	// unmarshel json and parse senz
	var senzMsg SenzMsg
	json.Unmarshal(b, &senzMsg)
	senz := parse(senzMsg.Msg)

	// get senzie key
	user, err := getUser(senz.Sender)
	if err != nil {
		senzResponse(w, "ERROR", senz.Attr["uid"], senz.Sender)
		return
	}

	// user needs to be active
	if !user.Active {
		senzResponse(w, "ERROR", senz.Attr["uid"], senz.Sender)
		return
	}

	// verify signature
	payload := strings.Replace(senz.Msg, senz.Digsig, "", -1)
	err = verify(payload, senz.Digsig, getSenzieRsaPub(user.PublicKey))
	if err != nil {
		senzResponse(w, "ERROR", senz.Attr["uid"], senz.Sender)
		return
	}

	// check for double spend first
	if isDoubleSpend(senz.Sender, senz.Attr["id"]) {
		// double spend response
		senzResponse(w, "DOUBLE_SPEND", senz.Attr["uid"], senz.Sender)
		return
	}

	// and new trans
	trans := senzToTrans(&senz)
	createTrans(trans)

	// TODO handle create failures

	// success response
	senzResponse(w, "SUCCESS", senz.Attr["uid"], senz.Sender)
	return
}

func senzResponse(w http.ResponseWriter, status string, uid string, to string) {
	// marshel and return error
	zmsg := SenzMsg{
		Uid: uid,
		Msg: statusSenz("ERROR", uid, to),
	}
	var zmsgs []SenzMsg
	zmsgs = append(zmsgs, zmsg)
	j, _ := json.Marshal(zmsgs)
	http.Error(w, string(j), 400)
}

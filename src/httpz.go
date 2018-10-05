package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type Zpreq struct {
	Uid              string
	Cid              string
	PrId             string
	CustomerId       string
	CustomerContact  string
	DeliveryDate     string
	DeliveryDue      int
	DeliveryLocation string
	ItemNo           string
	ItemQun          int
	ItemPrice        string
}

type Zporder struct {
	Uid    string
	Cid    string
	PoId   string
	OemId  string
	AmcId  string
	AmcApi string
	OemApi string
}

type Zdprep struct {
	Uid         string
	Cid         string
	PoId        string
	DataFile    string
	Instruction string
	AmcId       string
	AmcApi      string
	OemId       string
}

type Zprnt struct {
	Uid          string
	Cid          string
	SerialNumber string
	PrId         string
}

type Zdelnote struct {
	Uid       string
	CreatedAt string
	PoId      string
	Cid       string
	Status    string
	UpdatedAt string
}

type Zinvoice struct {
	Uid        string
	Cid        string
	DnId       string
	InId       string
	PoId       string
	Status     string
	TotalPrice string
	TotalQun   string
	CustomerId string
	Callback   string
}

type Zack struct {
	Uid        string
	Cid        string
	InId       string
	Status     string
	CustomerId string
}

type Zresp struct {
	Uid    string
	Status string
}

var rchans = make(map[string](chan string))

func initHttpz() {
	// router
	r := mux.NewRouter()
	r.HandleFunc("/prcontracts", prcontracts).Methods("POST")
	r.HandleFunc("/pocontracts", pocontracts).Methods("POST")
	r.HandleFunc("/dpcontracts", dpcontracts).Methods("POST")
	r.HandleFunc("/pcontracts", pcontracts).Methods("POST")
	r.HandleFunc("/delnotecontracts", delnotecontracts).Methods("POST")
	r.HandleFunc("/invoicecontracts", invoicecontracts).Methods("POST")
	r.HandleFunc("/ackcontracts", ackcontracts).Methods("POST")

	// start server
	err := http.ListenAndServe(":7070", r)
	if err != nil {
		println(err.Error)
		os.Exit(1)
	}
}

func prcontracts(w http.ResponseWriter, r *http.Request) {
	// read body
	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	println(string(b))

	// unmarshel json
	var req Zpreq
	json.Unmarshal(b, &req)
	senz := preqSenz(req)

	// create channel and add to rchans with uuid
	rchan := make(chan string)
	uid := req.Uid
	rchans[uid] = rchan

	// send to orderz(publish message to orderz topic)
	kmsg := Kmsg{
		Topic: "opsreq",
		Msg:   senz,
	}
	kchan <- kmsg

	senzResponse(w, "DONE")
}

func pocontracts(w http.ResponseWriter, r *http.Request) {
	// read body
	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	println(string(b))

	// unmarshel json
	var req Zporder
	json.Unmarshal(b, &req)
	senz := pordSenz(req)

	// send to orderz(publish message to orderz topic)
	kmsg := Kmsg{
		Topic: "opsreq",
		Msg:   senz,
	}
	kchan <- kmsg

	senzResponse(w, "DONE")
}

func dpcontracts(w http.ResponseWriter, r *http.Request) {
	// read body
	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	println(string(b))

	// unmarshel json
	var req Zdprep
	json.Unmarshal(b, &req)
	senz := dprepSenz(req)

	// send to orderz(publish message to orderz topic)
	kmsg := Kmsg{
		Topic: "opsreq",
		Msg:   senz,
	}
	kchan <- kmsg

	senzResponse(w, "DONE")
}

func pcontracts(w http.ResponseWriter, r *http.Request) {
	// read body
	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	println(string(b))

	// unmarshel json
	var req Zprnt
	json.Unmarshal(b, &req)
	senz := prntSenz(req)

	// send to orderz(publish message to orderz topic)
	kmsg := Kmsg{
		Topic: "opsreq",
		Msg:   senz,
	}
	kchan <- kmsg

	senzResponse(w, "DONE")
}

func delnotecontracts(w http.ResponseWriter, r *http.Request) {
	// read body
	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	println(string(b))

	// unmarshel json
	var req Zdelnote
	json.Unmarshal(b, &req)
	senz := delnoteSenz(req)

	// send to orderz(publish message to orderz topic)
	kmsg := Kmsg{
		Topic: "opsreq",
		Msg:   senz,
	}
	kchan <- kmsg

	senzResponse(w, "DONE")
}

func invoicecontracts(w http.ResponseWriter, r *http.Request) {
	// read body
	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	println(string(b))

	// unmarshel json
	var req Zinvoice
	json.Unmarshal(b, &req)
	senz := invoiceSenz(req)

	// send to orderz(publish message to orderz topic)
	kmsg := Kmsg{
		Topic: "opsreq",
		Msg:   senz,
	}
	kchan <- kmsg

	senzResponse(w, "DONE")
}

func ackcontracts(w http.ResponseWriter, r *http.Request) {
	// read body
	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	println(string(b))

	// unmarshel json
	var req Zack
	json.Unmarshal(b, &req)
	senz := ackSenz(req)

	// send to orderz(publish message to orderz topic)
	kmsg := Kmsg{
		Topic: "opsreq",
		Msg:   senz,
	}
	kchan <- kmsg

	senzResponse(w, "DONE")
}

func senzResponse(w http.ResponseWriter, status string) {
	zresp := Zresp{
		Uid:    "3223323",
		Status: status,
	}
	j, _ := json.Marshal(zresp)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	io.WriteString(w, string(j))
}

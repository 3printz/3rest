package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type Zreq struct {
	Uid          string
	ItemNo       string
	Quantity     int
	DeliveryDate string
}

type Zresp struct {
	Uid    string
	Status string
}

var rchans = make(map[string](chan string))

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

	// unmarshel json
	var zreq Zreq
	json.Unmarshal(b, &zreq)
	senz := kafkaSenz(zreq)

	// create channel and add to rchans with uuid
	rchan := make(chan string)
	uid := zreq.Uid
	rchans[uid] = rchan

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

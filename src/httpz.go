package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

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

	// unmarshel json and parse senz
	var zmsg Zmsg
	json.Unmarshal(b, &zmsg)
	senz := parse(zmsg.Msg)

	// create channel and add to rchans with senz uuid
	rchan := make(chan string)
	uid := senz.Attr["uid"]
	rchans[uid] = rchan

	// send to orderz(publish message to orderz topic)
	kmsg := Kmsg{
		Topic: "opsreq",
		Msg:   zmsg.Msg,
	}
	kchan <- kmsg

	senzResponse(w, "DONE")
}

func waitForResponse(w http.ResponseWriter, r *http.Request, rchan chan string, uid string) {
	for {
		select {
		case resp := <-rchan:
			// send senzResponse back
			println(resp)

			// senz respone
			senzResponse(w, resp)

			// clear map
			delete(rchans, uid)

			return
		}
	}
}

func senzResponse(w http.ResponseWriter, resp string) {
	zmsg := Zmsg{
		Uid: "3223323",
		Msg: resp,
	}
	j, _ := json.Marshal(zmsg)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	io.WriteString(w, string(j))
}

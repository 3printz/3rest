package main

type Senz struct {
	Msg      string
	Uid      string
	Ztype    string
	Sender   string
	Receiver string
	Attr     map[string]string
	Digsig   string
}

func main() {
	// first init key pair
	setUpKeys()

	// init cassandra session
	//initCStarSession()

	// init kafka producer, consumer
	initKafkaz()

	// start http server
	initHttpz()
}

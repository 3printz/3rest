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

type Zmsg struct {
	Uid string
	Msg string
}

type Kmsg struct {
	Topic string
	Msg   string
}

func main() {
	// first init key pair
	setUpKeys()

	// init cassandra session
	initCStarSession()

	// init kafka producer, consumer
	initKafkaz()

	// start http server
	initHttpz()
}

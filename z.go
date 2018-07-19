package main

import (
	"fmt"
	"time"
)

var rchans = make(map[string](chan string))

func m() {
	c := make(chan string, 10)
	rchans["001"] = c
	go ops("001", c, 5)

	c = make(chan string)
	rchans["002"] = c
	go ops("002", c, 5)

	c = make(chan string)
	rchans["003"] = c
	go ops("003", c, 4)

	fmt.Printf("size: %d\n", len(rchans))

	// ticking
	tick := time.Tick(1000 * time.Millisecond)
	for {
		select {
		case <-tick:
			println("tick")
			push("001")
			push("002")
			push("003")
		}
	}
}

func push(uid string) {
	if c, ok := rchans[uid]; ok {
		c <- uid
	}
}

func ops(uid string, rchan chan string, peers int) {
	var i int = 0
	for {
		select {
		case r := <-rchan:
			fmt.Printf("ops %s %d\n", r, i)
			i = i + 1
			if i == peers {
				fmt.Print("done %s\n", uid)
				delete(rchans, uid)
				return
			}
		}
	}
}

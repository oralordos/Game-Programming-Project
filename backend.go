package main

import (
	"fmt"
	"time"

	"github.com/Oralordos/Game-Programming-Project/events"
)

const frameDelta = 33 * time.Millisecond

func init() {
	go backend()
}

func backend() {
	inChn := make(chan events.Event)
	eta := time.Now()
	for {
		select {
		case todo := <-inChn:
			fmt.Println("Something")
			fmt.Println(todo)
		case <-time.After(eta.Sub(time.Now())):
			fmt.Println("timeout")
			eta = time.Now().Add(frameDelta)
		}
	}
}

package main

import (
	"fmt"
	"time"

	"github.com/Oralordos/Game-Programming-Project/events"
)

func init() {
	go backend()
}

func backend() {
	inChn := make(chan events.Event)
	select {
	case todo := <-inChn:
		fmt.Println("Something")
		fmt.Println(todo)
	case <-time.After(time.Millisecond * 500):
		fmt.Println("timeout")
	}
}

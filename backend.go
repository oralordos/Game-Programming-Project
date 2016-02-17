package main

import (
	"fmt"
	"time"

	"github.com/Oralordos/Game-Programming-Project/events"
)

func main() {
	inChn := make(chan events.Event)
	go func() {
		select {
		case todo := <-inChn:
			fmt.Println("Something")
			fmt.Println(todo)
		case <-time.After(time.Millisecond * 500):
			fmt.Println("timeout")
		}
	}()
}

package main

import (
	"fmt"
	"time"

	"github.com/Oralordos/Game-Programming-Project/events"
	"github.com/Oralordos/Game-Programming-Project/graphics"
)

type testEvent string

func (t *testEvent) GetDirection() int {
	return events.DirSystem
}

func (t *testEvent) GetSubValue() int {
	return 1
}

func (t testEvent) String() string {
	return string(t)
}

func main() {
	err := graphics.Init()
	if err != nil {
		panic(err)
	}
	defer graphics.Quit()
	win, err := graphics.CreateWindow(800, 600, "Test")
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

	ch := make(chan events.Event)
	events.AddListener(ch, events.DirSystem, 0)
	e := testEvent("Testing Events!")
	events.SendEvent(&e)

	val := <-ch

	fmt.Println(val)

	time.Sleep(5 * time.Second)
}

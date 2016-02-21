package main

import (
	"fmt"
	"time"

	"github.com/Oralordos/Game-Programming-Project/events"
	"github.com/Oralordos/Game-Programming-Project/graphics"
	"github.com/deckarep/golang-set"
)

type testEvent string

func (_ *testEvent) GetDirection() int {
	return events.DirSystem
}

func (_ *testEvent) GetSubValue() mapset.Set {
	s := mapset.NewThreadUnsafeSet()
	s.Add(1)
	return s
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
	events.AddListener(ch, events.DirSystem, []int{1, 2, 3, 4})
	e := testEvent("Testing Events!")
	events.SendEvent(&e)

	val := <-ch

	fmt.Println(val)

	time.Sleep(5 * time.Second)
}

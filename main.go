package main

import (
	"time"

	"github.com/Oralordos/Game-Programming-Project/events"
	"github.com/Oralordos/Game-Programming-Project/graphics"
)

func main() {
	err := graphics.Init()
	if err != nil {
		panic(err)
	}
	defer graphics.Quit()

	u := graphics.NewUnit(30, 30, 32, 32, 1)
	defer u.Destroy()

	win, err := graphics.CreateWindow(800, 600, "Test", u)
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

	e := events.UnitMoved{
		ID:   1,
		NewX: 50,
		NewY: 50,
	}

	err = win.Update()
	if err != nil {
		panic(err)
	}

	time.Sleep(2000 * time.Millisecond)
	events.SendEvent(&e)
	time.Sleep(100 * time.Millisecond)

	err = win.Update()
	if err != nil {
		panic(err)
	}
	time.Sleep(2500 * time.Millisecond)
}

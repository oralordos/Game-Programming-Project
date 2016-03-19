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
	win, err := graphics.CreateWindow(800, 600, "Test")
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

	u := graphics.NewUnit(30, 30, 32, 32, 1)
	defer u.Close()

	e := events.UnitMoved{
		ID:   1,
		NewX: 50,
		NewY: 50,
	}

	err = win.Clear()
	if err != nil {
		panic(err)
	}
	err = u.Draw(win)
	if err != nil {
		panic(err)
	}
	win.Update()

	time.Sleep(2500 * time.Millisecond)
	events.SendEvent(&e)
	time.Sleep(100 * time.Millisecond)
	err = win.Clear()
	if err != nil {
		panic(err)
	}
	err = u.Draw(win)
	if err != nil {
		panic(err)
	}
	win.Update()
	time.Sleep(2500 * time.Millisecond)
}

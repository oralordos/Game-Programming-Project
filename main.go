package main

import (
	"log"
	"time"

	"github.com/Oralordos/Game-Programming-Project/events"
	"github.com/Oralordos/Game-Programming-Project/graphics"
)

func main() {
	err := graphics.Init()
	if err != nil {
		log.Fatalln(err)
	}
	defer graphics.Quit()

	u := graphics.NewUnit(30, 30, 32, 32, 1)
	defer u.Destroy()

	win, err := graphics.CreateWindow(800, 600, "Test", u)
	if err != nil {
		log.Fatalln(err)
	}
	defer win.Destroy()

	e := events.InputUpdate{
		ID: 1,
		X:  0.707,
		Y:  0.707,
	}

	e2 := events.InputUpdate{
		ID: 1,
		X:  0,
		Y:  0,
	}

	err = win.Update()
	if err != nil {
		log.Fatalln(err)
	}

	time.Sleep(900 * time.Millisecond)
	events.SendEvent(&e)
	time.Sleep(100 * time.Millisecond)

	err = win.Update()
	if err != nil {
		log.Fatalln(err)
	}

	time.Sleep(1900 * time.Millisecond)
	events.SendEvent(&e2)
	time.Sleep(100 * time.Millisecond)

	err = win.Update()
	if err != nil {
		log.Fatalln(err)
	}

	time.Sleep(2 * time.Second)
}

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

	create := &events.CreateUnit{
		ID: 1,
		X:  30,
		Y:  30,
		W:  32,
		H:  32,
	}

	frontend := NewPlayerFrontend()
	defer frontend.Destroy()

	events.SendEvent(create)
	time.Sleep(20 * time.Millisecond)
	frontend.AttachUnit(1)

	win, err := graphics.CreateWindow(800, 600, "Test", frontend)
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

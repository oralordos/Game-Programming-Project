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

	win, err := graphics.CreateWindow(800, 600, "Test")
	if err != nil {
		log.Fatalln(err)
	}
	defer win.Destroy()

	frontend := NewPlayerFrontend(win)

	events.SendEvent(create)
	time.Sleep(20 * time.Millisecond)
	frontend.AttachUnit(1)

	frontend.Mainloop()
}

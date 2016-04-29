package main

import (
	"log"
	"os"

	"github.com/Oralordos/Game-Programming-Project/events"
	"github.com/Oralordos/Game-Programming-Project/graphics"
)

func main() {
	err := graphics.Init()
	if err != nil {
		log.Fatalln(err)
	}
	defer graphics.Quit()

	win, err := graphics.CreateWindow(800, 600, "Test")
	if err != nil {
		log.Fatalln(err)
	}
	defer win.Destroy()

	if len(os.Args) == 1 {
		newBackEnd()
	}

	frontend := NewPlayerFrontend(win)

	if len(os.Args) == 1 {
		joinEvent := &events.PlayerJoin{
			UUID: frontend.id,
		}
		events.SendEvent(joinEvent)

		go StartNetworkListener()
		loadlevel := events.LoadLevel{
			FileName: "assets/testTiles.json",
		}
		events.SendEvent(&loadlevel)
	} else {
		NewNetworkBackend(os.Args[1])
	}
	frontend.Mainloop()
}

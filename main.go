package main

import (
	"log"
	"os"
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

	win, err := graphics.CreateWindow(800, 600, "Test")
	if err != nil {
		log.Fatalln(err)
	}
	defer win.Destroy()

	frontend := NewPlayerFrontend(win)

	if len(os.Args) == 1 {
		go backendLoop()
		create := &events.CreateUnit{
			ID: 1,
			X:  48,
			Y:  48,
			W:  32,
			H:  32,
		}
		create2 := &events.CreateUnit{
			ID: 2,
			X:  105,
			Y:  105,
			W:  32,
			H:  32,
		}

		tiles := [][]int{}
		collide := [][]bool{}
		for i := 0; i < 15; i++ {
			tiles = append(tiles, []int{})
			collide = append(collide, []bool{})
			for j := 0; j < 25; j++ {
				if i == 0 || j == 0 || i == 14 || j == 24 {
					tiles[i] = append(tiles[i], 1)
					collide[i] = append(collide[i], true)
				} else {
					tiles[i] = append(tiles[i], 0)
					collide[i] = append(collide[i], false)
				}
			}
		}
		change := &events.ChangeLevel{
			Images:     tiles,
			TileWidth:  32,
			TileHeight: 32,
			CollideMap: collide,
			Units:      []events.Event{create, create2},
		}

		go StartNetworkListener()
		time.Sleep(10 * time.Millisecond)
		events.SendEvent(change)
		frontend.AttachUnit(1)
	} else {
		NewNetworkBackend(os.Args[1])
		frontend.AttachUnit(2)
	}
	frontend.Mainloop()
}

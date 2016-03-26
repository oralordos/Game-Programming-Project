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

	tiles := [][]graphics.Tile{}
	for i := 0; i < 10; i++ {
		tiles = append(tiles, []graphics.Tile{})
		for j := 0; j < 10; j++ {
			if i == 0 || j == 0 || i == 9 || j == 9 {
				tiles[i] = append(tiles[i], graphics.NewTile(255, 255, 255, 255))
			} else {
				tiles[i] = append(tiles[i], graphics.NewTile(127, 127, 127, 255))
			}
		}
	}
	level := graphics.NewTilemap(tiles, 32, 32)

	win, err := graphics.CreateWindow(800, 600, "Test")
	if err != nil {
		log.Fatalln(err)
	}
	defer win.Destroy()

	frontend := NewPlayerFrontend(win)
	frontend.SetLevel(level)

	events.SendEvent(create)
	time.Sleep(20 * time.Millisecond)
	frontend.AttachUnit(1)

	frontend.Mainloop()
}

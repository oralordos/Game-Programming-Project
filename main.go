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
		X:  48,
		Y:  48,
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
	}
	// level := graphics.NewTilemap(tiles, 32, 32)

	win, err := graphics.CreateWindow(800, 600, "Test")
	if err != nil {
		log.Fatalln(err)
	}
	defer win.Destroy()

	frontend := NewPlayerFrontend(win)
	// frontend.SetLevel(level)

	events.SendEvent(change)
	events.SendEvent(create)
	time.Sleep(20 * time.Millisecond)
	frontend.AttachUnit(1)

	frontend.Mainloop()
}

package main

import (
	"time"

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
	time.Sleep(5 * time.Second)
}

package main

import (
	"log"
	"time"

	"github.com/Oralordos/Game-Programming-Project/events"
	"github.com/Oralordos/Game-Programming-Project/graphics"
	"github.com/veandco/go-sdl2/sdl"
)

type PlayerFrontend struct {
	player  int
	window  *graphics.Window
	units   []*graphics.Unit
	level   *graphics.Tilemap
	inputs  []InputSystem
	eventCh chan events.Event
	close   chan struct{}
}

func NewPlayerFrontend(win *graphics.Window) *PlayerFrontend {
	p := PlayerFrontend{
		window:  win,
		units:   []*graphics.Unit{},
		inputs:  []InputSystem{ExitInput{}, &KeyboardInput{}},
		eventCh: make(chan events.Event),
		close:   make(chan struct{}),
	}
	events.AddListener(p.eventCh, events.DirFront, 0)
	return &p
}

func (p *PlayerFrontend) SetLevel(level *graphics.Tilemap) {
	p.level = level
}

func (p *PlayerFrontend) Mainloop() {
	nextUpdate := time.After(10 * time.Microsecond)
loop:
	for {
		select {
		case ev := <-p.eventCh:
			p.processEvent(ev)
		case <-nextUpdate:
			nextUpdate = time.After(time.Second / 60)
			p.processInput()
			if err := p.window.Update(p.getDraw()); err != nil {
				log.Fatalln(err)
			}
		case _, ok := <-p.close:
			if !ok {
				break loop
			}
		}
	}
	events.RemoveListener(p.eventCh, events.DirFront, 0)
	for _, v := range p.units {
		events.SendEvent(&events.DestroyUnit{ID: v.GetID()})
	}
}

func (p *PlayerFrontend) Destroy() {
	close(p.close)
}

func (p *PlayerFrontend) processInput() {
	for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {
		input := Input{}
		for _, v := range p.inputs {
			earlyExit, in := v.ProcessEvent(ev, p)
			if in != nil {
				input.Combine(in)
			}
			if earlyExit {
				break
			}
		}
		input.Normalize()
		p.sendInput(input)
	}
}

func (p *PlayerFrontend) sendInput(in Input) {
	player := p.GetUnit(p.player)
	if player == nil {
		return
	}
	e := events.InputUpdate{
		ID: p.player,
		X:  in.X,
		Y:  in.Y,
	}
	events.SendEvent(&e)
}

func (p *PlayerFrontend) processEvent(ev events.Event) {
	switch e := ev.(type) {
	case *events.CreateUnit:
		p.units = append(p.units, graphics.NewUnit(e.X, e.Y, e.W, e.H, e.ID))
	case *events.DestroyUnit:
		for i, v := range p.units {
			if v.GetID() == e.ID {
				p.units = append(p.units[:i], p.units[i+1:]...)
				break
			}
		}
	case *events.ChangeLevel:
		tiles := [][]graphics.Tile{}
		for y, row := range e.Images {
			tiles = append(tiles, []graphics.Tile{})
			for _, img := range row {
				var newTile graphics.Tile
				if img == 0 {
					newTile = graphics.NewTile(e.TileWidth, e.TileHeight, 127, 127, 127, 255)
				} else {
					newTile = graphics.NewTile(e.TileWidth, e.TileHeight, 191, 191, 191, 255)
				}
				tiles[y] = append(tiles[y], newTile)
			}
		}
		p.level = graphics.NewTilemap(tiles, e.TileWidth, e.TileHeight)
	}
}

func (p *PlayerFrontend) AttachUnit(id int) {
	p.player = id
}

func (p *PlayerFrontend) GetUnit(id int) *graphics.Unit {
	for _, v := range p.units {
		if v.GetID() == id {
			return v
		}
	}
	return nil
}

func (p *PlayerFrontend) getDraw() graphics.Drawable {
	draw := make(graphics.CombinedDrawer, 0, len(p.units)+1)
	if p.level != nil {
		draw = append(draw, p.level.GetDrawable())
	}
	for _, v := range p.units {
		draw = append(draw, v.GetDrawable())
	}
	return &draw
}

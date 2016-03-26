package main

import (
	"log"
	"time"

	"github.com/Oralordos/Game-Programming-Project/events"
	"github.com/Oralordos/Game-Programming-Project/graphics"
	"github.com/veandco/go-sdl2/sdl"
)

type PlayerFrontend struct {
	player  *graphics.Unit
	window  *graphics.Window
	units   []*graphics.Unit
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

func (p *PlayerFrontend) Mainloop() {
	nextUpdate := time.After(10 * time.Microsecond)
loop:
	for {
		select {
		case ev := <-p.eventCh:
			p.processEvent(ev)
		case <-nextUpdate:
			nextUpdate = time.After(16666 * time.Microsecond)
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
			earlyExit, in := v.ProcessEvent(ev)
			if earlyExit {
				break
			}
			if in != nil {
				input.Combine(in)
			}
		}
		input.Normalize()
		p.sendInput(input)
	}
}

func (p *PlayerFrontend) sendInput(in Input) {
	if p.player == nil {
		return
	}
	e := events.InputUpdate{
		ID: p.player.GetID(),
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
	}
}

func (p *PlayerFrontend) AttachUnit(id int) {
	for _, v := range p.units {
		if v.GetID() == id {
			p.player = v
			break
		}
	}
}

func (p *PlayerFrontend) getDraw() graphics.Drawable {
	draw := make(graphics.CombinedDrawer, 0, len(p.units))
	for _, v := range p.units {
		draw = append(draw, v.GetDrawable())
	}
	return &draw
}

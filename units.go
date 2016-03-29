package main

import (
	"math"
	"time"

	"github.com/Oralordos/Game-Programming-Project/events"
)

type unit struct {
	unitID     int
	unitHealth int
	x, y       float64
	xV, yV     float64
	xAcl, yAcl float64
	typ        *unitType
	eventCh    chan events.Event
	close      chan struct{}
}

type unitType struct {
	//basic stats go here
	maxHealth int
	movement  float64
	W, H      int32
}

var PlayerT = &unitType{
	maxHealth: 10,
	movement:  5,
	W:         32,
	H:         32,
}

func NewUnit(x, y float64, typ *unitType, id int) *unit {
	u := &unit{
		id,
		typ.maxHealth,
		x, y,
		0, 0, 0, 0,
		typ,
		make(chan events.Event),
		make(chan struct{}),
	}
	events.AddListener(u.eventCh, events.DirSystem, u.unitID)
	go u.unitloop()
	return u
}

func (u *unit) unitloop() {
	eta := time.After(frameDelta)
loop:
	for {
		select {
		case ev := <-u.eventCh:
			u.processEvent(ev)
		case <-eta:
			eta = time.After(frameDelta)
			u.updateUnit()
		case _, ok := <-u.close:
			if !ok {
				break loop
			}
		}
	}
	events.RemoveListener(u.eventCh, events.DirSystem, u.unitID)
}

//has to do with everything about the unit
func (u *unit) updateUnit() {
	u.x += u.xV
	u.y += u.yV
	u.xV += u.xAcl + (-u.xV * 0.8)
	u.yV += u.yAcl + (-u.yV * 0.8)
	u.x = math.Max(0, math.Min(800, u.x))
	u.y = math.Max(0, math.Min(600, u.y))
	e := events.UnitMoved{
		ID:   u.unitID,
		NewX: u.x,
		NewY: u.y,
	}
	events.SendEvent(&e)
}

func (u *unit) processEvent(ev events.Event) {
	switch e := ev.(type) {
	case *events.InputUpdate:
		u.xAcl = e.X * u.typ.movement
		u.yAcl = e.Y * u.typ.movement
	}
}

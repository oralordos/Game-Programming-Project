package main

import (
	"fmt"
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
	backAccess *BackEnd
	eventCh    chan events.Event
	close      chan struct{}
}

type unitType struct {
	//basic stats go here
	maxHealth int
	movement  float64
	hitDetect mainRect
}

type mainRect struct {
	topBox    hitRect
	bottomBox hitRect
	leftBox   hitRect
	rightBox  hitRect
}

type hitRect struct {
	top    int32
	bottom int32
	left   int32
	right  int32
}

var PlayerT = createUnitType(10, 5, 32, 32)

func createUnitType(maxHealth int, movement float64, w, h int32) *unitType {
	return &unitType{
		maxHealth: maxHealth,
		movement:  movement,
		hitDetect: mainRect{
			topBox: hitRect{
				top:    h / -2,
				bottom: 0,
				left:   w / -2,
				right:  w / 2,
			},
			bottomBox: hitRect{
				top:    0,
				bottom: h / 2,
				left:   w / -2,
				right:  w / 2,
			},
			leftBox: hitRect{
				top:    h / -2,
				bottom: h / 2,
				left:   w / -2,
				right:  0,
			},
			rightBox: hitRect{
				top:    h / -2,
				bottom: h / 2,
				left:   0,
				right:  w / 2,
			},
		},
	}
}

func (u hitRect) boxPosition(x, y int32) hitRect {
	return hitRect{
		top:    u.top + y,
		bottom: u.bottom - y,
		left:   u.left - x,
		right:  u.right + x,
	}
}

func (u mainRect) collusionBox(x, y int32) mainRect {
	return mainRect{
		topBox:    u.topBox.boxPosition(x, y),
		bottomBox: u.bottomBox.boxPosition(x, y),
		leftBox:   u.leftBox.boxPosition(x, y),
		rightBox:  u.rightBox.boxPosition(x, y),
	}
}

func NewUnit(x, y float64, typ *unitType, id int, b *BackEnd) *unit {
	u := &unit{
		id,
		typ.maxHealth,
		x, y,
		0, 0, 0, 0,
		typ,
		b,
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
	fmt.Println(u.backAccess.lastLevel.CollideMap[int32(u.y)/u.backAccess.lastLevel.TileHeight][int32(u.x)/u.backAccess.lastLevel.TileWidth])
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

func (u *unit) Destroy() {
	close(u.close)
}

func (u *unit) processEvent(ev events.Event) {
	switch e := ev.(type) {
	case *events.InputUpdate:
		u.xAcl = e.X * u.typ.movement
		u.yAcl = e.Y * u.typ.movement
	}
}

package main

import (
	"container/heap"
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

type coord struct {
	x, y int
}
type sortCoord struct {
	coords     []coord
	dirX, dirY int
}

func (s sortCoord) Len() int {
	return len(s.coords)
}

func (s sortCoord) Less(i, j int) bool {
	first := s.coords[i]
	second := s.coords[j]
	if s.dirX == 1 {
		return first.x < second.x
	} else if s.dirX == -1 {
		return first.x > second.x
	} else if s.dirY == 1 {
		return first.y < second.y
	}
	return first.y > second.y
}

func (s sortCoord) Swap(i, j int) {
	s.coords[i], s.coords[j] = s.coords[j], s.coords[i]
}

func (s *sortCoord) Push(x interface{}) {
	s.coords = append(s.coords, x.(coord))
}

func (s *sortCoord) Pop() interface{} {
	old := s.coords
	n := len(old)
	x := old[n-1]
	s.coords = old[:n-1]
	return x
}

func (u hitRect) checkRect(collide [][]bool, dirX, dirY int, tH, tW int32) (int, int) {
	collision := sortCoord{
		coords: []coord{},
		dirX:   dirX,
		dirY:   dirY,
	}

	leftT := u.left / tW
	rightT := u.right / tW
	topT := u.top / tH
	bottomT := u.bottom / tH
	for y, row := range collide[topT : bottomT+1] {
		for x, tile := range row[leftT : rightT+1] {
			if tile {
				collision.coords = append(collision.coords, coord{x + int(leftT), y + int(topT)})
			}
		}
	}

	if len(collision.coords) == 0 {
		return -1, -1
	}
	heap.Init(&collision)

	return collision.coords[0].x, collision.coords[0].y
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
				left:   w/-2 + 2,
				right:  w/2 - 1 - 2,
			},
			bottomBox: hitRect{
				top:    0,
				bottom: h/2 - 1,
				left:   w/-2 + 2,
				right:  w/2 - 1 - 2,
			},
			leftBox: hitRect{
				top:    h/-2 + 2,
				bottom: h/2 - 1 - 2,
				left:   w / -2,
				right:  0,
			},
			rightBox: hitRect{
				top:    h/-2 + 2,
				bottom: h/2 - 1 - 2,
				left:   0,
				right:  w/2 - 1,
			},
		},
	}
}

func (u hitRect) boxPosition(x, y int32) hitRect {
	return hitRect{
		top:    u.top + y,
		bottom: u.bottom + y,
		left:   u.left + x,
		right:  u.right + x,
	}
}

func (u mainRect) collisionBox(x, y int32) mainRect {
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
	curRect := u.typ.hitDetect.collisionBox(int32(u.x), int32(u.y))
	newX := u.x + u.xV
	if u.xV > 0 {
		//right
		currBox := curRect.rightBox
		currBox.right += int32(u.xV)
		tileX, _ := currBox.checkRect(u.backAccess.lastLevel.CollideMap, 1, 0, u.backAccess.lastLevel.TileWidth, u.backAccess.lastLevel.TileHeight)
		if tileX != -1 {
			newX = float64(int32(tileX)*u.backAccess.lastLevel.TileWidth - u.typ.hitDetect.rightBox.right - 1)
			u.xV = 0
		}
	} else {
		//left
		currbox := curRect.leftBox
		currbox.left += int32(u.xV)
		tileX, _ := currbox.checkRect(u.backAccess.lastLevel.CollideMap, -1, 0, u.backAccess.lastLevel.TileWidth, u.backAccess.lastLevel.TileHeight)
		if tileX != -1 {
			newX = float64(int32(tileX+1)*u.backAccess.lastLevel.TileWidth - u.typ.hitDetect.leftBox.left)
			u.xV = 0
		}
	}

	curRect = u.typ.hitDetect.collisionBox(int32(newX), int32(u.y))
	newY := u.y + u.yV
	if u.yV > 0 {
		//down
		currbox := curRect.bottomBox
		currbox.bottom += int32(u.yV)
		_, tileY := currbox.checkRect(u.backAccess.lastLevel.CollideMap, 0, 1, u.backAccess.lastLevel.TileWidth, u.backAccess.lastLevel.TileHeight)
		if tileY != -1 {
			newY = float64(int32(tileY)*u.backAccess.lastLevel.TileHeight - u.typ.hitDetect.bottomBox.bottom - 1)
			u.yV = 0
		}
	} else {
		//up
		currBox := curRect.topBox
		currBox.top += int32(u.yV)
		_, tileY := currBox.checkRect(u.backAccess.lastLevel.CollideMap, 0, -1, u.backAccess.lastLevel.TileWidth, u.backAccess.lastLevel.TileHeight)
		if tileY != -1 {
			newY = float64(int32(tileY+1)*u.backAccess.lastLevel.TileHeight - u.typ.hitDetect.topBox.top)
			u.yV = 0
		}
	}

	u.x = newX
	u.y = newY
	u.xV += u.xAcl + (-u.xV * 0.8)
	u.yV += u.yAcl + (-u.yV * 0.8)
	if u.backAccess.lastLevel.Pits[int(u.y)][int(u.x)] {
		u.x = u.backAccess.lastLevel.StartX
		u.y = u.backAccess.lastLevel.StartY
	}
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

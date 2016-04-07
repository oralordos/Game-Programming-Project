package main

import (
	"time"

	"github.com/Oralordos/Game-Programming-Project/events"
)

//BackEnd is struct for backend mechanics and info
type BackEnd struct {
	unitInfo  []*unit
	lastLevel *events.ChangeLevel
}

const frameDelta = 33 * time.Millisecond

func backendLoop() {
	b := &BackEnd{}
	b.unitInfo = []*unit{}
	inChn := make(chan events.Event)
	events.AddListener(inChn, events.DirSystem, 0)
	for {
		select {
		case todo := <-inChn:
			b.processEvent(todo)
		}
	}
}

func (b *BackEnd) processEvent(ev events.Event) {
	switch e := ev.(type) {
	case *events.CreateUnit:
		b.unitInfo = append(b.unitInfo, NewUnit(e.X, e.Y, PlayerT, e.ID, b))
	case *events.ChangeLevel:
		b.lastLevel = e
		for _, unit := range b.unitInfo {
			unit.Destroy()
		}
		b.unitInfo = make([]*unit, 0, len(e.Units))
		for _, unit := range e.Units {
			b.processEvent(&unit)
		}
	case events.ReloadLevel:
		units := make([]events.CreateUnit, len(b.unitInfo))
		for i, unit := range b.unitInfo {
			units[i] = events.CreateUnit{
				ID: unit.unitID,
				X:  unit.x,
				Y:  unit.y,
				W:  32,
				H:  32,
			}
		}
		newLevel := &events.ChangeLevel{
			Tilemap:    b.lastLevel.Tilemap,
			Images:     b.lastLevel.Images,
			TileWidth:  b.lastLevel.TileWidth,
			TileHeight: b.lastLevel.TileHeight,
			CollideMap: b.lastLevel.CollideMap,
			Units:      units,
		}
		events.SendEvent(newLevel)
	}
}

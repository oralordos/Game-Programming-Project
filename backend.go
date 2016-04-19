package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Oralordos/Game-Programming-Project/events"
)

//BackEnd is struct for backend mechanics and info
type BackEnd struct {
	unitInfo  []*unit
	lastLevel *events.ChangeLevel
	players   map[string]int
	nextID    int
	inChn     chan events.Event
}

const frameDelta = time.Second / 30

func newBackEnd() *BackEnd {
	b := &BackEnd{
		nextID:   1,
		unitInfo: []*unit{},
		players:  map[string]int{},
		inChn:    make(chan events.Event),
	}
	events.AddListener(b.inChn, events.DirSystem, 0)
	go b.backendLoop()
	return b
}

func (b *BackEnd) backendLoop() {
	for {
		select {
		case todo := <-b.inChn:
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
		b.nextID = 1
		for _, unit := range e.Units {
			b.processEvent(&unit)
			if b.nextID <= unit.ID {
				b.nextID = unit.ID + 1
			}
		}
		for k, v := range b.players {
			if _, ok := e.Players[k]; !ok {
				if b.nextID <= v {
					b.nextID = v + 1
				}
				b.createPlayerUnit(v, k)
			}
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
		players := make(map[string]int)
		for k, v := range b.players {
			players[k] = v
		}
		newLevel := &events.ChangeLevel{
			Tilemap:    b.lastLevel.Tilemap,
			Images:     b.lastLevel.Images,
			TileWidth:  b.lastLevel.TileWidth,
			TileHeight: b.lastLevel.TileHeight,
			StartX:     b.lastLevel.StartX,
			StartY:     b.lastLevel.StartY,
			CollideMap: b.lastLevel.CollideMap,
			Units:      units,
			Players:    players,
		}
		events.SendEvent(newLevel)
	case *events.PlayerJoin:
		id := b.nextID
		b.players[e.UUID] = id
		b.nextID++
		if b.lastLevel == nil {
			return
		}
		b.createPlayerUnit(id, e.UUID)
	case *events.LoadLevel:
		b.loadLevel(e)
	}
}

func (b *BackEnd) createPlayerUnit(id int, uuid string) {
	createPlayer := events.CreateUnit{
		ID:       id,
		X:        b.lastLevel.StartX,
		Y:        b.lastLevel.StartY,
		W:        32,
		H:        32,
		AttachTo: uuid,
	}
	events.SendEvent(&createPlayer)
}

func (b *BackEnd) loadLevel(e *events.LoadLevel) {
	type level struct {
		Height int32
		Width  int32
		Layers []struct {
			Data       []int32
			Properties map[string]string
		}
		Tilesets []struct {
			Image      string
			TileWidth  int32
			TileHeight int32
		}
	}

	file, err := os.Open(e.FileName)
	if err != nil {
		log.Printf("failed loading file: %s\n", e.FileName)
		return
	}
	defer file.Close()

	var x level
	err = json.NewDecoder(file).Decode(&x)
	if err != nil {
		log.Printf("failed loading file: %s\n", e.FileName)
		return
	}
	fmt.Println(x)
}

package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
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
	case *events.DestroyUnit:
		for i, unit := range b.unitInfo {
			if unit.unitID == e.ID {
				b.unitInfo = append(b.unitInfo[:i], b.unitInfo[i+1:]...)
				break
			}
		}
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
		newLevel := &events.ChangeLevel{
			// TODO Change this to copy the contents of Tilemaps
			Tilemaps:   b.lastLevel.Tilemaps,
			Images:     b.lastLevel.Images,
			TileWidth:  b.lastLevel.TileWidth,
			TileHeight: b.lastLevel.TileHeight,
			StartX:     b.lastLevel.StartX,
			StartY:     b.lastLevel.StartY,
			CollideMap: b.lastLevel.CollideMap,
			Pits:       b.lastLevel.Pits,
			Units:      units,
			Players:    b.players,
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
	case *events.PlayerLeave:
		destroy := events.DestroyUnit{
			ID: b.players[e.UUID],
		}
		events.SendEvent(&destroy)
		delete(b.players, e.UUID)
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
		Height     int32
		Width      int32
		Tileheight int32
		Tilewidth  int32
		Layers     []struct {
			Data       []int
			Properties map[string]string
			Height     int32
			Width      int32
		}
		Tilesets []struct {
			Image          string
			Tilewidth      int32
			Tileheight     int32
			Tileproperties map[string]map[string]string
		}
		Properties map[string]string
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

	//changeLevel struct
	//load the change

	startX, err := strconv.Atoi(x.Properties["StartX"])
	startY, err := strconv.Atoi(x.Properties["StartY"])
	cLevel := events.ChangeLevel{
		Tilemaps:   make([]events.Tilemap, 0, len(x.Tilesets)),
		Images:     [][][]int{},
		TileWidth:  x.Tilewidth,
		TileHeight: x.Tileheight,
		StartX:     float64(int32(startX) * x.Tilewidth),
		StartY:     float64(int32(startY) * x.Tileheight),
		CollideMap: make([][]bool, x.Height),
		Units:      make([]events.CreateUnit, 0, len(b.players)),
		Players:    map[string]int{},
		Pits:       make([][]bool, x.Height),
	}

	b.nextID = 1
	for playerID := range b.players {
		create := events.CreateUnit{
			ID:       b.nextID,
			X:        cLevel.StartX,
			Y:        cLevel.StartY,
			W:        32,
			H:        32,
			AttachTo: playerID,
		}
		b.players[playerID] = b.nextID
		cLevel.Players[playerID] = b.nextID
		cLevel.Units = append(cLevel.Units, create)
		b.nextID++
	}

	for i := range cLevel.CollideMap {
		cLevel.CollideMap[i] = make([]bool, x.Width)
		cLevel.Pits[i] = make([]bool, x.Width)
	}

	for _, t := range x.Tilesets {
		tm := events.Tilemap{
			Filename:   t.Image,
			TileWidth:  t.Tilewidth,
			TileHeight: t.Tileheight,
		}
		cLevel.Tilemaps = append(cLevel.Tilemaps, tm)
	}

	for z, layer := range x.Layers {
		cLevel.Images = append(cLevel.Images, [][]int{})
		var i int32
		for i = 0; i < layer.Height; i++ {
			cLevel.Images[z] = append(cLevel.Images[z], layer.Data[i*layer.Width:(i+1)*layer.Width])
		}
		for yC, row := range cLevel.Images[z] {
			for xC, tile := range row {
				for _, tileset := range x.Tilesets {
					for k, v := range tileset.Tileproperties {
						if v["Pit"] == "True" {
							if strconv.Itoa(tile) == k {
								cLevel.Pits[yC][xC] = true
							}
						}
					}
				}
				if layer.Properties["collide"] == "true" {
					if tile != 0 {
						cLevel.CollideMap[yC][xC] = true
					}
				}
			}
		}
	}
	events.SendEvent(&cLevel)
}

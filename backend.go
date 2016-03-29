package main

import (
	"time"

	"github.com/Oralordos/Game-Programming-Project/events"
)

//BackEnd is struct for backend mechanics and info
type BackEnd struct {
	unitInfo []*unit
	eventCh  chan events.Event
}

const frameDelta = 33 * time.Millisecond

func init() {
	go backendLoop()
}

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
		b.unitInfo = append(b.unitInfo, NewUnit(e.X, e.Y, PlayerT, e.ID))
	}
}

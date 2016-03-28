package main

import (
	"fmt"
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
	go backendLoop() //potentual problem here
}

func backendLoop() {
	b := &BackEnd{}
	b.unitInfo = []*unit{}
	inChn := make(chan events.Event)
	events.AddListener(inChn, events.DirSystem, 0)
	for {
		ev := <-inChn
		fmt.Printf("%T\n", ev)
		select {
		case todo := <-inChn:
			b.processEvent(todo) //potentual problem here
		}
	}
}

func (b *BackEnd) processEvent(ev events.Event) {
	switch e := ev.(type) {
	case *events.CreateUnit:
		b.unitInfo = append(b.unitInfo, NewUnit(e.X, e.Y, PlayerT, e.ID))
	}
}

// switch e := todo.(type) {
// case *events.UnitMoved:
// 	ev := events.UnitMoved{
// 		ID: 1,
// 		X:  0,
// 		Y:  -1,
// 	}
// 	events.SendEvent(ev)
// }

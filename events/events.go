package events

const (
	DirSystem = 1
	DirFront  = 2
)

type Event interface {
	getDirection() int
	getSubValue() int
}

type EventManager struct {
	input      chan Event
	systemOuts []chan<- Event
	frontOuts  []chan<- Event
}

func NewEventManager() *EventManager {
	return &EventManager{
		input:      make(chan Event),
		systemOuts: []chan<- Event{},
		frontOuts:  []chan<- Event{},
	}
}

func (e *EventManager) AddListener(listener chan<- Event, direction int) {
	if direction&DirSystem == DirSystem {
		e.systemOuts = append(e.systemOuts, listener)
	}
	if direction&DirFront == DirFront {
		e.frontOuts = append(e.systemOuts, listener)
	}
}

func (e *EventManager) SendEvent(event Event) {
	go func() {
		e.input <- event
	}()
}

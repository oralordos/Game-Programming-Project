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

var DefaultEventManager = NewEventManager()

func NewEventManager() *EventManager {
	return &EventManager{
		input:      make(chan Event),
		systemOuts: []chan<- Event{},
		frontOuts:  []chan<- Event{},
	}
}

func (e *EventManager) AddListener(listener chan<- Event, direction, subVal int) {
	if direction&DirSystem == DirSystem {
		e.systemOuts = append(e.systemOuts, listener)
	}
	if direction&DirFront == DirFront {
		e.frontOuts = append(e.systemOuts, listener)
	}
}

func (e *EventManager) RemoveListener(listener chan<- Event, direction, subVal int) {
	if direction&DirSystem == DirSystem {
		for i, v := range e.systemOuts {
			if v == listener {
				e.systemOuts = append(e.systemOuts[:i], e.systemOuts[i+1:]...)
				break
			}
		}
	}
	if direction&DirFront == DirFront {
		for i, v := range e.frontOuts {
			if v == listener {
				e.frontOuts = append(e.frontOuts[:i], e.frontOuts[i+1:]...)
				break
			}
		}
	}
}

func (e *EventManager) SendEvent(event Event) {
	go func() {
		e.input <- event
	}()
}

func AddListener(listener chan<- Event, direction, subVal int) {
	DefaultEventManager.AddListener(listener, direction, subVal)
}

func RemoveListener(listener chan<- Event, direction, subVal int) {
	DefaultEventManager.RemoveListener(listener, direction, subVal)
}

func SendEvent(event Event) {
	DefaultEventManager.SendEvent(event)
}

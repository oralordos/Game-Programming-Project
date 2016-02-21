package events

import "github.com/deckarep/golang-set"

const (
	DirSystem = 1
	DirFront  = 2
)

type Event interface {
	GetDirection() int
	GetSubValue() mapset.Set
}

type listener struct {
	ch        chan<- Event
	direction int
	subVal    mapset.Set
}

type EventManager struct {
	addList    chan *listener
	removeList chan *listener
	eInput     chan Event
	systemOuts []*listener
	frontOuts  []*listener
	close      chan struct{}
}

var DefaultEventManager = NewEventManager()

func NewEventManager() *EventManager {
	man := &EventManager{
		addList:    make(chan *listener),
		removeList: make(chan *listener),
		eInput:     make(chan Event),
		systemOuts: []*listener{},
		frontOuts:  []*listener{},
		close:      make(chan struct{}),
	}
	go man.mainloop()
	return man
}

func (e *EventManager) mainloop() {
eventLoop:
	for {
		select {
		case newList := <-e.addList:
			e.add(newList)
		case oldList := <-e.removeList:
			e.remove(oldList)
		case event := <-e.eInput:
			e.event(event)
		case _, ok := <-e.close:
			if !ok {
				break eventLoop
			}
		}
	}
}

func (e *EventManager) add(list *listener) {
	if list.direction&DirSystem == DirSystem {
		e.systemOuts = append(e.systemOuts, list)
	}
	if list.direction&DirFront == DirFront {
		e.frontOuts = append(e.systemOuts, list)
	}
}

func (e *EventManager) remove(list *listener) {
	if list.direction&DirSystem == DirSystem {
		for i, v := range e.systemOuts {
			if v.ch == list.ch && v.direction == list.direction && v.subVal.Equal(list.subVal) {
				e.systemOuts = append(e.systemOuts[:i], e.systemOuts[i+1:]...)
				break
			}
		}
	}
	if list.direction&DirFront == DirFront {
		for i, v := range e.frontOuts {
			if v.ch == list.ch && v.direction == list.direction && v.subVal.Equal(list.subVal) {
				e.frontOuts = append(e.frontOuts[:i], e.frontOuts[i+1:]...)
				break
			}
		}
	}
}

func (e *EventManager) event(ev Event) {
	dir := ev.GetDirection()
	var list []*listener
	if dir&DirFront == DirFront {
		list = e.frontOuts
	} else if dir&DirSystem == DirSystem {
		list = e.systemOuts
	}
	for _, v := range list {
		if v.subVal.Union(ev.GetSubValue()).Cardinality() != 0 {
			go func(ch chan<- Event) {
				ch <- ev
			}(v.ch)
		}
	}
}

func (e *EventManager) AddListener(list chan<- Event, direction int, subVal []int) {
	s := mapset.NewThreadUnsafeSet()
	for _, v := range subVal {
		s.Add(v)
	}
	l := &listener{
		ch:        list,
		direction: direction,
		subVal:    s,
	}
	e.addList <- l
}

func (e *EventManager) RemoveListener(list chan<- Event, direction int, subVal []int) {
	s := mapset.NewThreadUnsafeSet()
	for _, v := range subVal {
		s.Add(v)
	}
	l := &listener{
		ch:        list,
		direction: direction,
		subVal:    s,
	}
	e.removeList <- l
}

func (e *EventManager) SendEvent(event Event) {
	go func() {
		e.eInput <- event
	}()
}

func (e *EventManager) Close() {
	close(e.close)
}

func AddListener(listener chan<- Event, direction int, subVal []int) {
	DefaultEventManager.AddListener(listener, direction, subVal)
}

func RemoveListener(listener chan<- Event, direction int, subVal []int) {
	DefaultEventManager.RemoveListener(listener, direction, subVal)
}

func SendEvent(event Event) {
	DefaultEventManager.SendEvent(event)
}

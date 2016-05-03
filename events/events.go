package events

const (
	DirSystem = 1
	DirFront  = 2
)

type Event interface {
	GetDirection() int
	GetSubValue() int
	SetDuplicate(d bool)
	HasDuplicate() bool
	GetTypeID() int
}

type duplicateOnce bool

func (o *duplicateOnce) SetDuplicate(d bool) {
	*o = duplicateOnce(d)
}

func (o *duplicateOnce) HasDuplicate() bool {
	return bool(*o)
}

type noDuplicate struct{}

func (n noDuplicate) SetDuplicate(d bool) {}

func (n noDuplicate) HasDuplicate() bool {
	return true
}

type listener struct {
	ch        chan<- Event
	direction int
	subVal    int
}

type EventManager struct {
	eventQueue   []Event
	sendChannels []*listener
	currEv       Event
	addList      chan *listener
	removeList   chan *listener
	eInput       chan Event
	systemOuts   []*listener
	frontOuts    []*listener
	close        chan struct{}
}

var DefaultEventManager = NewEventManager()

func NewEventManager() *EventManager {
	man := &EventManager{
		eventQueue:   []Event{},
		sendChannels: []*listener{},
		addList:      make(chan *listener),
		removeList:   make(chan *listener),
		eInput:       make(chan Event),
		systemOuts:   []*listener{},
		frontOuts:    []*listener{},
		close:        make(chan struct{}),
	}
	go man.mainloop()
	return man
}

func (e *EventManager) mainloop() {
eventLoop:
	for {
		e.updateCurrEvent()
		if e.currEv == nil {
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
		} else {
			select {
			case e.sendChannels[0].ch <- e.currEv:
				e.sendChannels = e.sendChannels[1:]
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
}

func (e *EventManager) updateCurrEvent() {
	if len(e.sendChannels) == 0 {
		e.currEv = nil
	}
	for e.currEv == nil && len(e.eventQueue) > 0 {
		e.currEv = e.eventQueue[0]
		e.eventQueue = e.eventQueue[1:]
		list := []*listener{}
		if e.currEv.GetDirection()&DirSystem == DirSystem {
			list = append(list, e.systemOuts...)
		}
		if e.currEv.GetDirection()&DirFront == DirFront {
			list = append(list, e.frontOuts...)
		}
		for _, v := range list {
			if v.subVal == 0 || v.subVal == e.currEv.GetSubValue() {
				e.sendChannels = append(e.sendChannels, v)
			}
		}
		if len(e.sendChannels) == 0 {
			e.currEv = nil
		}
	}
}

func (e *EventManager) add(list *listener) {
	if list.direction&DirSystem == DirSystem {
		e.systemOuts = append(e.systemOuts, list)
	} else if list.direction&DirFront == DirFront {
		e.frontOuts = append(e.frontOuts, list)
	}
}

func (e *EventManager) remove(list *listener) {
	var l *[]*listener
	if list.direction&DirSystem == DirSystem {
		l = &e.systemOuts
	} else if list.direction&DirFront == DirFront {
		l = &e.frontOuts
	}
	for i, v := range *l {
		if list.equals(v) {
			*l = append((*l)[:i], (*l)[i+1:]...)
			break
		}
	}
	for i, v := range e.sendChannels {
		if list.equals(v) {
			e.sendChannels = append(e.sendChannels[:i], e.sendChannels[i+1:]...)
			break
		}
	}
}

func (e *EventManager) event(ev Event) {
	e.eventQueue = append(e.eventQueue, ev)
}

func (e *EventManager) AddListener(list chan<- Event, direction int, subVal int) {
	l := &listener{
		ch:        list,
		direction: direction,
		subVal:    subVal,
	}
	e.addList <- l
}

func (e *EventManager) RemoveListener(list chan<- Event, direction int, subVal int) {
	l := &listener{
		ch:        list,
		direction: direction,
		subVal:    subVal,
	}
	e.removeList <- l
}

func (e *EventManager) SendEvent(event Event) {
	e.eInput <- event
}

func (e *EventManager) Close() {
	close(e.close)
}

func (l *listener) equals(other *listener) bool {
	return other.ch == l.ch && other.direction == l.direction && other.subVal == l.subVal
}

func AddListener(listener chan<- Event, direction int, subVal int) {
	DefaultEventManager.AddListener(listener, direction, subVal)
}

func RemoveListener(listener chan<- Event, direction int, subVal int) {
	DefaultEventManager.RemoveListener(listener, direction, subVal)
}

func SendEvent(event Event) {
	DefaultEventManager.SendEvent(event)
}

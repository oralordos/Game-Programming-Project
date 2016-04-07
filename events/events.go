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

type listener struct {
	ch        chan<- Event
	direction int
	subVal    int
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
		e.frontOuts = append(e.frontOuts, list)
	}
}

func (e *EventManager) remove(list *listener) {
	if list.direction&DirSystem == DirSystem {
		for i, v := range e.systemOuts {
			if list.equals(v) {
				e.systemOuts = append(e.systemOuts[:i], e.systemOuts[i+1:]...)
				break
			}
		}
	}
	if list.direction&DirFront == DirFront {
		for i, v := range e.frontOuts {
			if list.equals(v) {
				e.frontOuts = append(e.frontOuts[:i], e.frontOuts[i+1:]...)
				break
			}
		}
	}
}

func (e *EventManager) event(ev Event) {
	dir := ev.GetDirection()
	list := []*listener{}
	if dir&DirFront == DirFront {
		list = append(list, e.frontOuts...)
	}
	if dir&DirSystem == DirSystem {
		list = append(list, e.systemOuts...)
	}
	for _, v := range list {
		if v.subVal == 0 || v.subVal == ev.GetSubValue() {
			go func(ch chan<- Event) {
				ch <- ev
			}(v.ch)
		}
	}
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
	go func() {
		e.eInput <- event
	}()
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

package graphics

import "github.com/Oralordos/Game-Programming-Project/events"

type Unit struct {
	id      int
	x, y    float64
	w, h    int32
	eventCh chan events.Event
	drawCh  chan Drawable
	close   chan struct{}
}

func NewUnit(x, y float64, w, h int32, id int) *Unit {
	u := &Unit{
		id,
		x, y, w, h,
		make(chan events.Event),
		make(chan Drawable),
		make(chan struct{}),
	}
	events.AddListener(u.eventCh, events.DirFront, u.id)
	go u.mainloop()
	return u
}

func (u *Unit) GetDrawable() Drawable {
	return <-u.drawCh
}

func (u *Unit) mainloop() {
loop:
	for {
		select {
		case ev := <-u.eventCh:
			u.processEvent(ev)
		case u.drawCh <- u.getDraw():
		case _, ok := <-u.close:
			if !ok {
				break loop
			}
		}
	}
	events.RemoveListener(u.eventCh, events.DirFront, u.id)
}

func (u *Unit) processEvent(ev events.Event) {
	switch e := ev.(type) {
	case *events.UnitMoved:
		u.x = e.NewX
		u.y = e.NewY
	case *events.DestroyUnit:
		u.Destroy()
	}
}

func (u *Unit) GetID() int {
	return u.id
}

func (u *Unit) getDraw() Drawable {
	return &RectDrawer{
		x: int32(u.x) - u.w/2,
		y: int32(u.y) - u.h/2,
		w: u.w,
		h: u.h,
		r: 255,
		g: 127,
		b: 0,
		a: 255,
	}
}

func (u *Unit) Destroy() {
	close(u.close)
}

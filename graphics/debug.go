package graphics

import (
	"github.com/Oralordos/Game-Programming-Project/events"
	"github.com/veandco/go-sdl2/sdl"
)

func drawRect(win *Window, r *sdl.Rect, c *sdl.Color) error {
	if err := win.rend.SetDrawColor(c.R, c.G, c.B, c.A); err != nil {
		return err
	}
	return win.rend.FillRect(r)
}

type Unit struct {
	id      int
	x, y    float64
	w, h    int32
	eventCh chan events.Event
	close   chan struct{}
}

func NewUnit(x, y float64, w, h int32, id int) *Unit {
	u := &Unit{
		id,
		x, y, w, h,
		make(chan events.Event),
		make(chan struct{}),
	}
	events.AddListener(u.eventCh, events.DirFront, u.id)
	go u.mainloop()
	return u
}

func (u *Unit) mainloop() {
loop:
	for {
		select {
		case ev := <-u.eventCh:
			switch e := ev.(type) {
			case *events.UnitMoved:
				u.x = e.NewX
				u.y = e.NewY
			}
		case _, ok := <-u.close:
			if !ok {
				break loop
			}
		}
	}
}

func (u *Unit) Draw(win *Window) error {
	rect := sdl.Rect{
		X: int32(u.x) - u.w/2,
		Y: int32(u.y) - u.h/2,
		W: u.w,
		H: u.h,
	}
	color := sdl.Color{
		R: 255,
		G: 127,
		B: 0,
		A: 255,
	}
	return drawRect(win, &rect, &color)
}

func (u *Unit) Close() {
	close(u.close)
}

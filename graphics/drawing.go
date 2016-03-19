package graphics

import "github.com/veandco/go-sdl2/sdl"

type UpdatableDrawable interface {
	GetDrawable() Drawable
}

type Drawable interface {
	Draw(r *sdl.Renderer) error
}

type RectDrawer struct {
	x, y, w, h int32
	r, g, b, a uint8
}

func (d *RectDrawer) Draw(r *sdl.Renderer) error {
	if err := r.SetDrawColor(d.r, d.g, d.b, d.a); err != nil {
		return err
	}
	rect := &sdl.Rect{
		X: d.x,
		Y: d.y,
		W: d.w,
		H: d.h,
	}
	return r.FillRect(rect)
}

type CombinedDrawer []Drawable

func (d CombinedDrawer) Draw(r *sdl.Renderer) error {
	for _, v := range d {
		if err := v.Draw(r); err != nil {
			return err
		}
	}
	return nil
}

package graphics

import "github.com/veandco/go-sdl2/sdl"

type Drawable interface {
	Draw(r *sdl.Renderer, offsetX, offsetY int32) error
}

type RectDrawer struct {
	x, y, w, h int32
	r, g, b, a uint8
}

func (d *RectDrawer) Draw(r *sdl.Renderer, offsetX, offsetY int32) error {
	if err := r.SetDrawColor(d.r, d.g, d.b, d.a); err != nil {
		return err
	}
	rect := &sdl.Rect{
		X: d.x + offsetX,
		Y: d.y + offsetY,
		W: d.w,
		H: d.h,
	}
	return r.FillRect(rect)
}

type CombinedDrawer []Drawable

func (d CombinedDrawer) Draw(r *sdl.Renderer, offsetX, offsetY int32) error {
	for _, v := range d {
		if err := v.Draw(r, offsetX, offsetY); err != nil {
			return err
		}
	}
	return nil
}

type OffsetDrawer struct {
	Drawing          Drawable
	OffsetX, OffsetY int32
}

func (d *OffsetDrawer) Draw(r *sdl.Renderer, offsetX, offsetY int32) error {
	d.Drawing.Draw(r, offsetX+d.OffsetX, offsetY+d.OffsetY)
	return nil
}

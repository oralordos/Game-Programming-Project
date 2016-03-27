package graphics

type Image struct {
	w, h       int32
	r, g, b, a uint8
}

func (i *Image) GetDrawable(x, y int32) Drawable {
	return &RectDrawer{
		x: x,
		y: y,
		w: i.w,
		h: i.h,
		r: i.r,
		g: i.g,
		b: i.b,
		a: i.a,
	}
}

package graphics

type Tile struct {
	r, g, b, a uint8
}

func NewTile(r, g, b, a uint8) Tile {
	return Tile{
		r: r,
		g: g,
		b: b,
		a: a,
	}
}

func (t *Tile) GetDrawable(x, y, w, h int32) Drawable {
	return &RectDrawer{
		x: x * w,
		y: y * h,
		w: w,
		h: h,
		r: t.r,
		g: t.g,
		b: t.b,
		a: t.a,
	}
}

type Tilemap struct {
	level                 [][]Tile
	tileWidth, tileHeight int32
}

func NewTilemap(level [][]Tile, tileWidth, tileHeight int32) *Tilemap {
	return &Tilemap{
		level:      level,
		tileWidth:  tileWidth,
		tileHeight: tileHeight,
	}
}

func (t *Tilemap) GetDrawable() Drawable {
	draw := CombinedDrawer{}
	for x, row := range t.level {
		for y, tile := range row {
			draw = append(draw, tile.GetDrawable(int32(x), int32(y), t.tileWidth, t.tileHeight))
		}
	}
	return draw
}

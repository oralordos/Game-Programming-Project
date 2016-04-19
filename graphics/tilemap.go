package graphics

type Tile struct {
	img Image
}

func NewTile(w, h int32, r, g, b, a uint8) Tile {
	return Tile{
		img: Image{
			w: w,
			h: h,
			r: r,
			g: g,
			b: b,
			a: a,
		},
	}
}

func (t *Tile) GetDrawable(x, y int32) Drawable {
	return t.img.GetDrawable(x, y)
}

type Tilemap struct {
	level                 [][][]Tile
	tileWidth, tileHeight int32
}

func NewTilemap(level [][][]Tile, tileWidth, tileHeight int32) *Tilemap {
	return &Tilemap{
		level:      level,
		tileWidth:  tileWidth,
		tileHeight: tileHeight,
	}
}

func (t *Tilemap) GetDrawable() Drawable {
	draw := CombinedDrawer{}
	for _, layer := range t.level {
		for y, row := range layer {
			for x, tile := range row {
				draw = append(draw, tile.GetDrawable(int32(x)*t.tileWidth, int32(y)*t.tileHeight))
			}
		}
	}
	return draw
}

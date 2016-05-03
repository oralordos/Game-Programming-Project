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

func (t *Tilemap) GetDrawable(offsetX, offsetY, screenW, screenH int32) Drawable {
	startX := -offsetX / t.tileWidth
	startY := -offsetY / t.tileHeight
	endX := startX + screenW/t.tileWidth + 2
	endY := startY + screenH/t.tileHeight + 2
	if startX < 0 {
		startX = 0
	}
	if startY < 0 {
		startY = 0
	}
	if endX >= int32(len(t.level[0][0])) {
		endX = int32(len(t.level[0][0]))
	}
	if endY >= int32(len(t.level[0])) {
		endY = int32(len(t.level[0]))
	}
	draw := CombinedDrawer{}
	for _, layer := range t.level {
		for y := startY; y < endY; y++ {
			row := layer[y]
			for x := startX; x < endX; x++ {
				tile := row[x]
				draw = append(draw, tile.GetDrawable(int32(x)*t.tileWidth, int32(y)*t.tileHeight))
			}
		}
		// for y, row := range layer {
		// 	for x, tile := range row {
		// 		draw = append(draw, tile.GetDrawable(int32(x)*t.tileWidth, int32(y)*t.tileHeight))
		// 	}
		// }
	}
	return draw
}

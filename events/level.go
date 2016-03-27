package events

type ChangeLevel struct {
	Tilemap               string
	Images                [][]int
	TileWidth, TileHeight int32
	CollideMap            [][]bool
}

func (c *ChangeLevel) GetDirection() int {
	return DirFront | DirSystem
}

func (c *ChangeLevel) GetSubValue() int {
	return 0
}

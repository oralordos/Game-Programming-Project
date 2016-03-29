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

func isChangeLevel(items []string) bool {
	return isMatch(items, []string{"Tilemap", "Images", "TileWidth", "TileHeight", "CollideMap"})
}

func getChangeLevel(data map[string]interface{}) Event {
	e := ChangeLevel{}

	tilemap, ok := data["Tilemap"].(string)
	if !ok {
		return nil
	}
	e.Tilemap = tilemap

	tilewidth, ok := data["TileWidth"].(float64)
	if !ok {
		return nil
	}
	e.TileWidth = int32(tilewidth + 0.5)

	tileheight, ok := data["TileWidth"].(float64)
	if !ok {
		return nil
	}
	e.TileHeight = int32(tileheight + 0.5)

	e.Images = get2Dint(data["Images"])
	if e.Images == nil {
		return nil
	}

	e.CollideMap = get2Dbool(data["CollideMap"])
	if e.CollideMap == nil {
		return nil
	}

	return &e
}

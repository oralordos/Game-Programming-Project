package events

type ReloadLevel struct{}

func (c ReloadLevel) GetDirection() int {
	return DirSystem
}

func (c ReloadLevel) GetSubValue() int {
	return 0
}

type ChangeLevel struct {
	Tilemap               string
	Images                [][]int
	TileWidth, TileHeight int32
	CollideMap            [][]bool
	Units                 []Event
}

func (c *ChangeLevel) GetDirection() int {
	return DirFront | DirSystem
}

func (c *ChangeLevel) GetSubValue() int {
	return 0
}

func isChangeLevel(items []string) bool {
	return isMatch(items, []string{"Tilemap", "Images", "TileWidth", "TileHeight", "CollideMap", "Units"})
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

	units, ok := data["Units"].([]interface{})
	if !ok {
		return nil
	}
	results := make([]Event, len(units))
	for i, v := range units {
		unitData, ok := v.(map[string]interface{})
		if !ok {
			return nil
		}
		unit := DecodeJSON(unitData)
		if unit == nil {
			return nil
		}
		results[i] = unit
	}

	return &e
}

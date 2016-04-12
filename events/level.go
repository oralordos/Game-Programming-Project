package events

type ReloadLevel struct {
	noDuplicate `json:"-"`
}

func (c ReloadLevel) GetDirection() int {
	return DirSystem
}

func (c ReloadLevel) GetSubValue() int {
	return 0
}

func (c ReloadLevel) GetTypeID() int {
	return TypeReloadLevel
}

type ChangeLevel struct {
	Tilemap               string
	Images                [][]int
	TileWidth, TileHeight int32
	StartX, StartY        float64
	CollideMap            [][]bool
	Units                 []CreateUnit
	Players               map[string]int
	duplicateOnce         `json:"-"`
}

func (c *ChangeLevel) GetDirection() int {
	return DirFront | DirSystem
}

func (c *ChangeLevel) GetSubValue() int {
	return 0
}

func (c *ChangeLevel) GetTypeID() int {
	return TypeChangeLevel
}

type LoadLevel struct {
	FileName    string
	noDuplicate `json:"-"`
}

func (c *LoadLevel) GetDirection() int {
	return DirSystem
}

func (c *LoadLevel) GetSubValue() int {
	return 0
}

func (c *LoadLevel) GetTypeID() int {
	return TypeLoadLevel
}

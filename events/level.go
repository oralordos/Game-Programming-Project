package events

type ReloadLevel struct{}

func (c ReloadLevel) GetDirection() int {
	return DirSystem
}

func (c ReloadLevel) GetSubValue() int {
	return 0
}

func (c ReloadLevel) SetDuplicate(d bool) {}

func (c ReloadLevel) HasDuplicate() bool {
	return true
}

func (c ReloadLevel) GetTypeID() int {
	return TypeReloadLevel
}

type ChangeLevel struct {
	Tilemap               string
	Images                [][]int
	TileWidth, TileHeight int32
	CollideMap            [][]bool
	Units                 []CreateUnit
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

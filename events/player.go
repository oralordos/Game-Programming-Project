package events

type PlayerJoin struct {
	UUID        string
	noDuplicate `json:"-"`
}

func (p *PlayerJoin) GetDirection() int {
	return DirSystem
}

func (p *PlayerJoin) GetSubValue() int {
	return 0
}

func (p *PlayerJoin) GetTypeID() int {
	return TypePlayerJoin
}

type PlayerLeave struct {
	UUID        string
	noDuplicate `json:"-"`
}

func (p *PlayerLeave) GetDirection() int {
	return DirSystem
}

func (p *PlayerLeave) GetSubValue() int {
	return 0
}

func (p *PlayerLeave) GetTypeID() int {
	return TypePlayerLeave
}

type SetUUID struct {
	UUID          string
	duplicateOnce `json:"-"`
}

func (p *SetUUID) GetDirection() int {
	return DirFront
}

func (p *SetUUID) GetSubValue() int {
	return 0
}

func (p *SetUUID) GetTypeID() int {
	return TypeSetUUID
}

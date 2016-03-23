package events

type UnitMoved struct {
	ID         int
	NewX, NewY float64
}

func (u *UnitMoved) GetDirection() int {
	return DirFront
}

func (u *UnitMoved) GetSubValue() int {
	return u.ID
}

type InputUpdate struct {
	ID   int
	X, Y float64
}

func (u *InputUpdate) GetDirection() int {
	return DirSystem
}

func (u *InputUpdate) GetSubValue() int {
	return u.ID
}

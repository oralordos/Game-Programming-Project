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

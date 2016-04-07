package events

type UnitMoved struct {
	ID            int
	NewX, NewY    float64
	duplicateOnce `json:"-"`
}

func (u *UnitMoved) GetDirection() int {
	return DirFront
}

func (u *UnitMoved) GetSubValue() int {
	return u.ID
}

func (u *UnitMoved) GetTypeID() int {
	return TypeUnitMoved
}

type InputUpdate struct {
	ID            int
	X, Y          float64
	duplicateOnce `json:"-"`
}

func (u *InputUpdate) GetDirection() int {
	return DirSystem
}

func (u *InputUpdate) GetSubValue() int {
	return u.ID
}

func (u *InputUpdate) GetTypeID() int {
	return TypeInputUpdate
}

type CreateUnit struct {
	ID            int
	X, Y          float64
	W, H          int32
	duplicateOnce `json:"-"`
}

func (u *CreateUnit) GetDirection() int {
	return DirFront | DirSystem
}

func (u *CreateUnit) GetSubValue() int {
	return 0
}

func (u *CreateUnit) GetTypeID() int {
	return TypeCreateUnit
}

type DestroyUnit struct {
	ID            int
	duplicateOnce `json:"-"`
}

func (u *DestroyUnit) GetDirection() int {
	return DirFront | DirSystem
}

func (u *DestroyUnit) GetSubValue() int {
	return u.ID
}

func (u *DestroyUnit) GetTypeID() int {
	return TypeDestroyUnit
}

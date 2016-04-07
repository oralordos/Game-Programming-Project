package events

type UnitMoved struct {
	ID         int
	NewX, NewY float64
	duplicateOnce
}

func (u *UnitMoved) GetDirection() int {
	return DirFront
}

func (u *UnitMoved) GetSubValue() int {
	return u.ID
}

func isUnitMoved(items []string) bool {
	return isMatch(items, []string{"ID", "NewX", "NewY", "duplicateOnce"})
}

func getUnitMoved(data map[string]interface{}) Event {
	e := UnitMoved{}

	id, ok := data["ID"].(float64)
	if !ok {
		return nil
	}
	e.ID = int(id + 0.5)

	newx, ok := data["NewX"].(float64)
	if !ok {
		return nil
	}
	e.NewX = newx

	newy, ok := data["NewY"].(float64)
	if !ok {
		return nil
	}
	e.NewY = newy

	dup, ok := data["duplicateOnce"].(bool)
	if !ok {
		return nil
	}
	e.duplicateOnce = duplicateOnce(dup)

	return &e
}

type InputUpdate struct {
	ID   int
	X, Y float64
	duplicateOnce
}

func (u *InputUpdate) GetDirection() int {
	return DirSystem
}

func (u *InputUpdate) GetSubValue() int {
	return u.ID
}

func isInputUpdate(items []string) bool {
	return isMatch(items, []string{"ID", "X", "Y", "duplicateOnce"})
}

func getInputUpdate(data map[string]interface{}) Event {
	e := InputUpdate{}

	id, ok := data["ID"].(float64)
	if !ok {
		return nil
	}
	e.ID = int(id + 0.5)

	x, ok := data["X"].(float64)
	if !ok {
		return nil
	}
	e.X = x

	y, ok := data["Y"].(float64)
	if !ok {
		return nil
	}
	e.Y = y

	dup, ok := data["duplicateOnce"].(bool)
	if !ok {
		return nil
	}
	e.duplicateOnce = duplicateOnce(dup)

	return &e
}

type CreateUnit struct {
	ID   int
	X, Y float64
	W, H int32
	duplicateOnce
}

func (u *CreateUnit) GetDirection() int {
	return DirFront | DirSystem
}

func (u *CreateUnit) GetSubValue() int {
	return 0
}

func isCreateUnit(items []string) bool {
	return isMatch(items, []string{"ID", "X", "Y", "W", "H", "duplicateOnce"})
}

func getCreateUnit(data map[string]interface{}) Event {
	e := CreateUnit{}

	id, ok := data["ID"].(float64)
	if !ok {
		return nil
	}
	e.ID = int(id + 0.5)

	x, ok := data["X"].(float64)
	if !ok {
		return nil
	}
	e.X = x

	y, ok := data["Y"].(float64)
	if !ok {
		return nil
	}
	e.Y = y

	w, ok := data["W"].(float64)
	if !ok {
		return nil
	}
	e.W = int32(w + 0.5)

	h, ok := data["H"].(float64)
	if !ok {
		return nil
	}
	e.H = int32(h + 0.5)

	dup, ok := data["duplicateOnce"].(bool)
	if !ok {
		return nil
	}
	e.duplicateOnce = duplicateOnce(dup)

	return &e
}

type DestroyUnit struct {
	ID int
	duplicateOnce
}

func (u *DestroyUnit) GetDirection() int {
	return DirFront | DirSystem
}

func (u *DestroyUnit) GetSubValue() int {
	return u.ID
}

func isDestroyUnit(items []string) bool {
	return isMatch(items, []string{"ID", "duplicateOnce"})
}

func getDestroyUnit(data map[string]interface{}) Event {
	e := DestroyUnit{}

	id, ok := data["ID"].(float64)
	if !ok {
		return nil
	}
	e.ID = int(id + 0.5)

	dup, ok := data["duplicateOnce"].(bool)
	if !ok {
		return nil
	}
	e.duplicateOnce = duplicateOnce(dup)

	return &e
}

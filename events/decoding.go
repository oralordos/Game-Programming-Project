package events

import (
	"encoding/json"
	"fmt"
)

const (
	TypeUnitMoved = iota
	TypeInputUpdate
	TypeCreateUnit
	TypeDestroyUnit
	TypeReloadLevel
	TypeChangeLevel
	TypePlayerJoin
	TypePlayerLeave
	TypeSetUUID
)

func DecodeJSON(typ int, data json.RawMessage) (Event, error) {
	var ev Event
	switch typ {
	case TypeUnitMoved:
		ev = new(UnitMoved)
	case TypeInputUpdate:
		ev = new(InputUpdate)
	case TypeCreateUnit:
		ev = new(CreateUnit)
	case TypeDestroyUnit:
		ev = new(DestroyUnit)
	case TypeReloadLevel:
		ev = ReloadLevel{}
	case TypeChangeLevel:
		ev = new(ChangeLevel)
	case TypePlayerJoin:
		ev = new(PlayerJoin)
	case TypePlayerLeave:
		ev = new(PlayerLeave)
	case TypeSetUUID:
		ev = new(SetUUID)
	default:
		return nil, fmt.Errorf("Unknown event type: %d\n", typ)
	}
	if err := json.Unmarshal(data, ev); err != nil {
		return nil, err
	}
	return ev, nil
}

package events

import (
	"encoding/json"
	"fmt"
	"log"
)

const (
	TypeUnitMoved = iota
	TypeInputUpdate
	TypeCreateUnit
	TypeDestroyUnit
	TypeReloadLevel
	TypeChangeLevel
)

func DecodeJSON(typ int, decod *json.Decoder) (Event, error) {
	switch typ {
	case TypeUnitMoved:
		var um UnitMoved
		if err := decod.Decode(&um); err != nil {
			log.Printf("Unable to decode: %d\n", typ)
			return nil, err
		}
		return &um, nil
	case TypeInputUpdate:
		var iu InputUpdate
		if err := decod.Decode(&iu); err != nil {
			return nil, err
		}
		return &iu, nil
	case TypeCreateUnit:
		var cu CreateUnit
		if err := decod.Decode(&cu); err != nil {
			return nil, err
		}
		return &cu, nil
	case TypeDestroyUnit:
		var du DestroyUnit
		if err := decod.Decode(&du); err != nil {
			return nil, err
		}
		return &du, nil
	case TypeReloadLevel:
		var rl ReloadLevel
		if err := decod.Decode(&rl); err != nil {
			return nil, err
		}
		return rl, nil
	case TypeChangeLevel:
		var cl ChangeLevel
		if err := decod.Decode(&cl); err != nil {
			log.Println("Got an error:", err)
			return nil, err
		}
		return &cl, nil
	}
	return nil, fmt.Errorf("Unknown event type: %d\n", typ)
}

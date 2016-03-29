package events

func DecodeJSON(data map[string]interface{}) Event {
	keys := []string{}
	for k := range data {
		keys = append(keys, k)
	}
	switch {
	case isChangeLevel(keys):
		return getChangeLevel(data)
	case isUnitMoved(keys):
		return getUnitMoved(data)
	case isInputUpdate(keys):
		return getInputUpdate(data)
	case isCreateUnit(keys):
		return getCreateUnit(data)
	case isDestroyUnit(keys):
		return getDestroyUnit(data)
	}
	return nil
}

func isMatch(items, allowed []string) bool {
	if len(items) != len(allowed) {
		return false
	}
	for _, v := range allowed {
		found := false
		for _, item := range items {
			if v == item {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func get2Dint(data interface{}) [][]int {
	dataTable, ok := data.([]interface{})
	if !ok {
		return nil
	}
	results := make([][]int, len(dataTable))
	for y, row := range dataTable {
		rowItems, ok := row.([]interface{})
		if !ok {
			return nil
		}
		results[y] = make([]int, len(rowItems))
		for x, item := range rowItems {
			itemVal, ok := item.(float64)
			if !ok {
				return nil
			}
			results[y][x] = int(itemVal + 0.5)
		}
	}
	return results
}

func get2Dbool(data interface{}) [][]bool {
	dataTable, ok := data.([]interface{})
	if !ok {
		return nil
	}
	results := make([][]bool, len(dataTable))
	for y, row := range dataTable {
		rowItems, ok := row.([]interface{})
		if !ok {
			return nil
		}
		results[y] = make([]bool, len(rowItems))
		for x, item := range rowItems {
			itemVal, ok := item.(bool)
			if !ok {
				return nil
			}
			results[y][x] = itemVal
		}
	}
	return results
}

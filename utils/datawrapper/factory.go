package datawrapper

import "github.com/mrizkifadil26/medix/utils"

func WrapData(val any) Data {
	switch v := val.(type) {
	case *utils.OrderedMap[string, any]:
		return NewOrderedMapData(v)
	case map[string]any:
		return NewMapData(v)
	case []any:
		return NewArrayData(v)
	default:
		return &ValueData{data: v}
	}
}

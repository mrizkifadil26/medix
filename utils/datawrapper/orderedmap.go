package datawrapper

import (
	"fmt"

	"github.com/mrizkifadil26/medix/utils"
)

type OrderedMapData struct {
	data *utils.OrderedMap[string, any]
}

func NewOrderedMapData(
	m *utils.OrderedMap[string, any],
) *OrderedMapData {
	return &OrderedMapData{data: m}
}

func (o *OrderedMapData) Get(key any) (Data, bool) {
	k, ok := key.(string)
	if !ok {
		return nil, false
	}

	v, exists := o.data.Get(k)
	if !exists {
		return nil, false
	}

	return WrapData(v), true
}

func (o *OrderedMapData) Set(key any, value any) error {
	k, ok := key.(string)
	if !ok {
		return fmt.Errorf("key must be string")
	}
	o.data.Set(k, value)
	return nil
}

func (o *OrderedMapData) Keys() []any {
	keys := make([]any, o.data.Len())
	for i, k := range o.data.Keys() {
		keys[i] = k
	}
	return keys
}

func (o *OrderedMapData) Append(value any) error {
	return fmt.Errorf("cannot append to ordered map")
}

func (o *OrderedMapData) Raw() any {
	return o.data
}

func (o *OrderedMapData) Type() string {
	return "map"
}

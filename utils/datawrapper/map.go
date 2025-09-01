package datawrapper

import "fmt"

type MapData struct {
	data map[string]any
}

func NewMapData(m map[string]any) *MapData {
	return &MapData{data: m}
}

func (m *MapData) Get(key any) (Data, bool) {
	k, ok := key.(string)
	if !ok {
		return nil, false
	}
	v, exists := m.data[k]
	if !exists {
		return nil, false
	}
	return WrapData(v), true
}

func (m *MapData) Set(key any, value any) error {
	k, ok := key.(string)
	if !ok {
		return fmt.Errorf("key must be string")
	}
	m.data[k] = value
	return nil
}

func (m *MapData) Keys() []any {
	keys := make([]any, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys
}

func (m *MapData) Append(value any) error {
	return fmt.Errorf("cannot append to map")
}

func (m *MapData) Raw() any {
	return m.data
}

func (m *MapData) Type() string {
	return "map"
}

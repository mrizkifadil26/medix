package datawrapper

import "fmt"

type ArrayData struct {
	data []any
}

func NewArrayData(a []any) *ArrayData {
	return &ArrayData{data: a}
}

func (a *ArrayData) Get(key any) (Data, bool) {
	idx, ok := key.(int)
	if !ok || idx < 0 || idx >= len(a.data) {
		return nil, false
	}
	return WrapData(a.data[idx]), true
}

func (a *ArrayData) Set(key any, value any) error {
	idx, ok := key.(int)
	if !ok || idx < 0 || idx >= len(a.data) {
		return fmt.Errorf("invalid index")
	}
	a.data[idx] = value
	return nil
}

func (a *ArrayData) Keys() []any {
	keys := make([]any, len(a.data))
	for i := range a.data {
		keys[i] = i
	}
	return keys
}

func (a *ArrayData) Append(value any) error {
	a.data = append(a.data, value)
	return nil
}

func (a *ArrayData) Raw() any {
	return a.data
}

func (a *ArrayData) Type() string {
	return "array"
}

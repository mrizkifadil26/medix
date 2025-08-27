package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type OrderedMap[K comparable, V any] struct {
	keys   []K
	values map[K]V
	index  map[K]int
}

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		values: make(map[K]V),
		index:  make(map[K]int),
	}
}

func (om *OrderedMap[K, V]) Set(key K, value V) {
	if _, exists := om.values[key]; !exists {
		om.keys = append(om.keys, key)
		om.index[key] = len(om.keys) - 1
	}

	om.values[key] = value
}

func (om *OrderedMap[K, V]) Get(key K) (V, bool) {
	v, ok := om.values[key]
	return v, ok
}

func (om *OrderedMap[K, V]) Delete(key K) {
	if _, exists := om.values[key]; !exists {
		return
	}
	delete(om.values, key)

	idx := om.index[key]
	delete(om.index, key)

	// shift keys slice
	om.keys = append(om.keys[:idx], om.keys[idx+1:]...)

	// rebuild index after idx
	for i := idx; i < len(om.keys); i++ {
		om.index[om.keys[i]] = i
	}
}

func (om *OrderedMap[K, V]) Keys() []K {
	return append([]K(nil), om.keys...) // copy
}

func (om *OrderedMap[K, V]) Values() []V {
	values := make([]V, 0, len(om.keys))
	for _, k := range om.keys {
		values = append(values, om.values[k])
	}
	return values
}

func (om *OrderedMap[K, V]) Len() int {
	return len(om.keys)
}

func (om *OrderedMap[K, V]) UnmarshalJSON(data []byte) error {
	om.keys = make([]K, 0)
	om.values = make(map[K]V)
	om.index = make(map[K]int)

	dec := json.NewDecoder(bytes.NewReader(data))

	tok, err := dec.Token()
	if err != nil {
		return err
	}

	if delim, ok := tok.(json.Delim); !ok || delim != '{' {
		return fmt.Errorf("OrderedMap: expected object")
	}

	for dec.More() {
		// read key
		tok, err := dec.Token()
		if err != nil {
			return err
		}
		key := tok.(string)

		// decode value (recursively)
		var raw json.RawMessage
		if err := dec.Decode(&raw); err != nil {
			return err
		}

		// decide type
		// safe assignment
		var val V
		if raw != nil && string(raw) != "null" {
			switch {
			case raw[0] == '{':
				child := NewOrderedMap[string, any]()
				if err := json.Unmarshal(raw, child); err != nil {
					return err
				}
				val = any(child).(V)
			case raw[0] == '[':
				var arr []json.RawMessage
				if err := json.Unmarshal(raw, &arr); err != nil {
					return err
				}
				newArr := make([]any, 0, len(arr))
				for _, item := range arr {
					if len(item) > 0 && item[0] == '{' {
						child := NewOrderedMap[string, any]()
						if err := json.Unmarshal(item, child); err != nil {
							return err
						}
						newArr = append(newArr, child)
					} else {
						var prim any
						if err := json.Unmarshal(item, &prim); err != nil {
							return err
						}
						newArr = append(newArr, prim)
					}
				}
				val = any(newArr).(V)
			default:
				var prim any
				if err := json.Unmarshal(raw, &prim); err != nil {
					return err
				}
				val = any(prim).(V)
			}
		} else {
			// assign zero value of V if JSON is null
			var zero V
			val = zero
		}

		om.Set(any(key).(K), val)
		// om.Set(any(key).(K), any(v).(V))
	}

	// read closing '}'
	_, err = dec.Token()
	return err
}

// MarshalJSON re-encodes OrderedMap preserving key order
func (om *OrderedMap[string, any]) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')

	for i, k := range om.keys {
		val := om.values[k]

		keyBytes, _ := json.Marshal(k)
		valBytes, err := json.Marshal(val)
		if err != nil {
			return nil, err
		}

		if i > 0 {
			buf.WriteByte(',')
		}

		buf.Write(keyBytes)
		buf.WriteByte(':')
		buf.Write(valBytes)
	}

	buf.WriteByte('}')
	return buf.Bytes(), nil
}

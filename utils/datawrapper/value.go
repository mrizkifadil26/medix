package datawrapper

import (
	"fmt"
	"strconv"
)

type ValueData struct {
	data any
}

func (v *ValueData) Get(key any) (Data, bool)     { return nil, false }
func (v *ValueData) Set(key any, value any) error { return fmt.Errorf("cannot set value") }
func (v *ValueData) Keys() []any                  { return nil }
func (v *ValueData) Append(value any) error       { return fmt.Errorf("cannot append") }
func (v *ValueData) Raw() any                     { return v.data }
func (v *ValueData) Type() string                 { return "value" }

// --- primitive helpers ---
func (v *ValueData) String() (string, bool) {
	if v == nil || v.data == nil {
		return "", false
	}

	switch val := v.data.(type) {
	case string:
		return val, true
	case fmt.Stringer:
		return val.String(), true
	case []byte:
		return string(val), true
	default:
		return fmt.Sprintf("%v", val), true
	}
}

func (v *ValueData) Int64() (int64, bool) {
	if v == nil || v.data == nil {
		return 0, false
	}
	switch val := v.data.(type) {
	case int:
		return int64(val), true
	case int64:
		return val, true
	case float64: // JSON numbers usually land here
		return int64(val), true
	case string:
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			return i, true
		}
	}
	return 0, false
}

func (v *ValueData) Float64() (float64, bool) {
	if v == nil || v.data == nil {
		return 0, false
	}
	switch val := v.data.(type) {
	case float32:
		return float64(val), true
	case float64:
		return val, true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case string:
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f, true
		}
	}
	return 0, false
}

func (v *ValueData) Bool() (bool, bool) {
	if v == nil || v.data == nil {
		return false, false
	}
	switch val := v.data.(type) {
	case bool:
		return val, true
	case string:
		if val == "true" || val == "1" {
			return true, true
		} else if val == "false" || val == "0" {
			return false, true
		}
	}
	return false, false
}

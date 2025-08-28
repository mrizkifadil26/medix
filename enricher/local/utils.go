package local

import (
	"fmt"

	"github.com/mrizkifadil26/medix/utils/datawrapper"
)

func safeString(d datawrapper.Data, key string, errs *[]error) string {
	if d == nil {
		return ""
	}
	child, ok := d.Get(key)
	if !ok || child == nil {
		return ""
	}
	if s, ok := child.(*datawrapper.ValueData); ok {
		val, ok := s.String()
		if ok {
			return val
		}
	}
	if s, ok := child.Raw().(string); ok {
		return s
	}
	*errs = append(*errs, fmt.Errorf("expected string at key %q but got %T", key, child.Raw()))
	return ""
}

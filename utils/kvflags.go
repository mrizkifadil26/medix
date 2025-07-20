package utils

import (
	"fmt"
	"strings"
)

type KVFlags struct {
	values map[string]string
	order  []string // optional: keep insertion order
}

func (k *KVFlags) String() string {
	if k.values == nil {
		return ""
	}
	var parts []string
	for _, key := range k.order {
		parts = append(parts, fmt.Sprintf("%s:%s", key, k.values[key]))
	}
	return strings.Join(parts, ",")
}

func (k *KVFlags) Set(value string) error {
	if k.values == nil {
		k.values = make(map[string]string)
	}
	// Try to split as label:path
	parts := strings.SplitN(value, ":", 2)

	if len(parts) == 2 {
		label := parts[0]
		path := parts[1]
		k.values[label] = path
		k.order = append(k.order, label)
	} else {
		// No label, fallback: use path as label
		path := parts[0]
		k.values[path] = path
		k.order = append(k.order, path)
	}
	return nil
}

func (k *KVFlags) ToMap() map[string]string {
	return k.values
}

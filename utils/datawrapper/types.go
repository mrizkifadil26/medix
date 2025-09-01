package datawrapper

type Data interface {
	// Get returns the value at a key or index
	Get(key any) (Data, bool)

	// Set sets a value at a key or index
	Set(key any, value any) error

	// Keys returns keys for a map or indices for an array
	Keys() []any

	// Append adds a value to a slice (if applicable)
	Append(value any) error

	// Raw returns the underlying data
	Raw() any

	// Type returns the kind: "map", "array", "value"
	Type() string
}

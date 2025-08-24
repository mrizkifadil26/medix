package utils

// --- Generic Registry ---
type Registry[T any] struct {
	items map[string]T
}

func NewRegistry[T any]() *Registry[T] {
	return &Registry[T]{items: make(map[string]T)}
}

func (r *Registry[T]) Register(name string, fn T) {
	r.items[name] = fn
}

func (r *Registry[T]) Get(name string) (T, bool) {
	fn, ok := r.items[name]
	return fn, ok
}

func (r *Registry[T]) All() map[string]T {
	return r.items
}

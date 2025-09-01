package cache

import (
	"encoding/json"
	"os"
	"sync"
)

type Manager[T any] struct {
	mu       sync.Mutex
	filepath string
	data     map[string]map[string]T // category → key → value
}

func NewManager[T any](filepath string) *Manager[T] {
	return &Manager[T]{
		filepath: filepath,
		data:     make(map[string]map[string]T),
	}
}

// Load file into memory (if exists)
func (m *Manager[T]) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := os.ReadFile(m.filepath)
	if err != nil {
		// no cache yet
		return nil
	}

	return json.Unmarshal(data, &m.data)
}

// Save memory → file
func (m *Manager[T]) Save() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	enc, err := json.MarshalIndent(m.data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.filepath, enc, 0644)
}

func (m *Manager[T]) Has(category, key string) bool {
	_, ok := m.Get(category, key)
	return ok
}

// Get returns cached value if exists
func (m *Manager[T]) Get(category, key string) (T, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if cat, ok := m.data[category]; ok {
		if val, ok := cat[key]; ok {
			return val, true
		}
	}

	var zero T
	return zero, false
}

// Put updates memory (does not auto-save)
func (m *Manager[T]) Put(category, key string, value T) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.data[category]; !ok {
		m.data[category] = make(map[string]T)
	}

	m.data[category][key] = value
}

func (m *Manager[T]) Delete(category, key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if cat, ok := m.data[category]; ok {
		delete(cat, key)
		if len(cat) == 0 {
			delete(m.data, category) // cleanup empty category
		}
	}
}

func (m *Manager[T]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = make(map[string]map[string]T)
}

func (m *Manager[T]) Categories() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	cats := make([]string, 0, len(m.data))
	for k := range m.data {
		cats = append(cats, k)
	}

	return cats
}

func (m *Manager[T]) Keys(category string) []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	if cat, ok := m.data[category]; ok {
		keys := make([]string, 0, len(cat))
		for k := range cat {
			keys = append(keys, k)
		}
		return keys
	}

	return nil
}

package cache

import (
	"encoding/json"
	"os"
	"sync"
)

type Manager struct {
	mu       sync.Mutex
	filepath string
	data     map[string]map[string]any // category → key → value
}

func NewManager(filepath string) *Manager {
	return &Manager{
		filepath: filepath,
		data:     make(map[string]map[string]any),
	}
}

// Load file into memory (if exists)
func (m *Manager) Load() error {
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
func (m *Manager) Save() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	enc, err := json.MarshalIndent(m.data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.filepath, enc, 0644)
}

// Get returns cached value if exists
func (m *Manager) Get(category, key string) (any, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if cat, ok := m.data[category]; ok {
		val, ok := cat[key]
		return val, ok
	}

	return nil, false
}

// Put updates memory (does not auto-save)
func (m *Manager) Put(category, key string, value any) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.data[category]; !ok {
		m.data[category] = make(map[string]any)
	}

	m.data[category][key] = value
}

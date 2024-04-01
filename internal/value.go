package internal

import (
	"reflect"
	"sync"
)

type Value[V any] struct {
	_    noCopy // go vet to alert when copying by value.
	mu   sync.RWMutex
	set  bool
	data V
}

// NewValue returns a new Value.
func NewValue[V any]() *Value[V] {
	return &Value[V]{}
}

// NewWithValue returns a new Value, set to the specified value.
func NewWithValue[V any](v V) *Value[V] {
	return &Value[V]{data: v, set: true}
}

// Load returns the value stored, ok indicates whether value was previously set.
func (m *Value[V]) Load() (v V, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data, m.set
}

// Store sets the value.
func (m *Value[V]) Store(value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = value
	m.set = true
}

// LoadOrStore returns the existing value if present. Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (m *Value[V]) LoadOrStore(value V) (actual V, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.set {
		return m.data, true
	}
	m.data = value
	m.set = true
	return value, false
}

// Swap swaps the value for a key and returns the previous value if any.
// The loaded result reports whether the key was present.
func (m *Value[V]) Swap(value V) (previous V, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	previous = m.data
	loaded = m.set
	m.data = value
	m.set = true
	return previous, loaded
}

// CompareAndSwap swaps the old and new values if the value stored in the map is equal to old.
func (m *Value[V]) CompareAndSwap(old, new V) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.set && reflect.DeepEqual(m.data, old) {
		m.data = new
		return true
	}
	return false
}

// IsZero returns true if the value is a zero value (not set). It relies on the set flag.
func (m *Value[V]) IsZero() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return !m.set
}

// Exclusive executes the function f exclusively, ensuring that no other goroutine is accessing the value.
func (m *Value[V]) Exclusive(update func(actual V, loaded bool) V) (updated V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = update(m.data, m.set)
	m.set = true
	return m.data
}

// Clear set the value with the zero value and set flag to false.
func (m *Value[V]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	var zero V
	m.data = zero
	m.set = false
}

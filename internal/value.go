package internal

import (
	"reflect"
	"sync"
)

type Value[V any] struct {
	mu   sync.RWMutex
	set  bool
	data V
}

func NewValue[V any]() *Value[V] {
	return &Value[V]{}
}

func NewWithValue[V any](v V) *Value[V] {
	return &Value[V]{data: v, set: true}
}

func (m *Value[V]) Load() (v V, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data, m.set
}

func (m *Value[V]) Store(value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = value
	m.set = true
}

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

func (m *Value[V]) Swap(value V) (previous V, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	previous = m.data
	loaded = m.set
	m.data = value
	m.set = true
	return previous, loaded
}

func (m *Value[V]) CompareAndSwap(old, new V) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.set && reflect.DeepEqual(m.data, old) {
		m.data = new
		return true
	}
	return false
}

func (m *Value[V]) IsZero() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return !m.set
}

func (m *Value[V]) Exclusive(update func(actual V, loaded bool) V) (updated V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = update(m.data, m.set)
	m.set = true
	return m.data
}

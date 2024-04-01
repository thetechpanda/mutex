package internal

import (
	"reflect"
	"sync"
)

// Map implements a simple thread-safe map that uses generics.
type Map[K comparable, V any] struct {
	_    noCopy // go vet to alert when copying by value.
	mu   sync.RWMutex
	data map[K]V
}

// New returns a new Map, initialized with the given map. if m is nil, an empty map is created.
// m key, values are copied, so that the caller can safely modify the map after creating a Map.
func NewMap[K comparable, V any](m map[K]V) *Map[K, V] {
	var v map[K]V = make(map[K]V, len(m))
	for key, value := range m {
		v[key] = value
	}
	return &Map[K, V]{data: v}
}

// Store sets the value for a key.
func (m *Map[K, V]) Store(key K, value V) {
	m.Swap(key, value)
}

// Load returns the value stored in the map for a key, or nil if no
// value is present.
// The ok result indicates whether value was found in the map.
func (m *Map[K, V]) Load(key K) (v V, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok = m.data[key]
	return v, ok
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	actual, loaded = m.data[key]
	if loaded {
		return actual, true
	}
	m.data[key] = value
	return value, false
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	value, loaded = m.data[key]
	if loaded {
		delete(m.data, key)
	}
	return value, loaded
}

// Delete removes the key from the map.
// This is a locking operation.
func (m *Map[K, V]) Delete(key K) {
	m.LoadAndDelete(key)
}

// Swap swaps the value for a key and returns the previous value if any.
// The loaded result reports whether the key was present.
func (m *Map[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	previous, loaded = m.data[key]
	m.data[key] = value
	return previous, loaded
}

// CompareAndSwap swaps the old and new values for key
// if the value stored in the map is equal to old.
//
// Returns true if the swap was performed.
//
// ! this function uses reflect.DeepEqual to compare the values.
func (m *Map[K, V]) CompareAndSwap(key K, old, new V) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.data[key]
	if !ok || !reflect.DeepEqual(v, old) {
		return false
	}
	m.data[key] = new
	return true
}

// CompareAndDelete deletes the entry for key if its value is equal to old.
//
// If there is no current value for key in the map, CompareAndDelete
// returns false (even if the old value is the nil interface value).
//
// ! this function uses reflect.DeepEqual to compare the values.
func (m *Map[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.data[key]
	if !ok || !reflect.DeepEqual(v, old) {
		return false
	}
	delete(m.data, key)
	return true
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, Range stops the iteration.
// Avoid invoking any map functions within 'f' to prevent a deadlock.
func (m *Map[K, V]) Range(f func(K, V) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for key, value := range m.data {
		if !f(key, value) {
			break
		}
	}
}

// Clear removes all items from the map.
// This is a locking operation.
func (m *Map[K, V]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = make(map[K]V)
}

// Has returns true if the map contains the key.
func (m *Map[K, V]) Has(key K) bool {
	_, ok := m.Load(key)
	return ok
}

// Update allows the caller to change the value associated with the key atomically guaranteeing that the value would not be changed by another goroutine during the operation.
//
// ! Do not invoke any Map functions within 'f' to prevent a deadlock.
func (m *Map[K, V]) Update(key K, f func(V, bool) V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.data[key]
	m.data[key] = f(v, ok)
}

// UpdateRange is a thread-safe version of Range that locks the map for the duration of the iteration and allows for the modification of the values.
// If f returns false, UpdateRange stops the iteration, without updating the corresponding value in the map.
//
// ! Do not invoke any Map functions within 'f' to prevent a deadlock.
func (m *Map[K, V]) UpdateRange(f func(K, V) (V, bool)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for key, value := range m.data {
		newValue, ok := f(key, value)
		if !ok {
			return
		}
		m.data[key] = newValue
	}
}

// Exclusive provides a way to perform  operations on the map ensuring that no other operation is performed on the map during the execution of the function.
//
// ! Do not invoke any Map functions within 'f' to prevent a deadlock.
func (m *Map[K, V]) Exclusive(f func(m map[K]V)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	f(m.data)
}

// Len returns the number of items in the map.
func (m *Map[K, V]) Len() (n int) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}

// Keys returns a slice of all the keys present in the map, an empty slice is returned if the map is empty.
func (m *Map[K, V]) Keys() (keys []K) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	max := len(m.data)
	if max == 0 {
		return make([]K, 0)
	}
	keys = make([]K, 0, max)
	for key := range m.data {
		keys = append(keys, key)
		if len(keys) >= max {
			break
		}
	}
	return keys
}

// Values returns a slice of all the values present in the map, an empty slice is returned if the map is empty.
func (m *Map[K, V]) Values() (values []V) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	max := len(m.data)
	if max == 0 {
		return make([]V, 0)
	}
	values = make([]V, 0, max)
	for _, value := range m.data {
		values = append(values, value)
		if len(values) >= max {
			break
		}
	}
	return values
}

// Entries returns two slices, one containing all the keys and the other containing all the values present in the map.
func (m *Map[K, V]) Entries() (keys []K, values []V) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	max := len(m.data)
	if max == 0 {
		return make([]K, 0), make([]V, 0)
	}
	keys = make([]K, 0, max)
	values = make([]V, 0, max)
	for key, value := range m.data {
		keys = append(keys, key)
		values = append(values, value)
		if len(keys) >= max {
			break
		}
	}
	return keys, values
}

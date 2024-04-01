package mutex

import (
	"github.com/thetechpanda/mutex/internal"
)

// Value is an interface that represents a thread-safe value store.
// It provides methods to load, store, and manipulate the stored value.
// The value can be of any type specified by the type parameter V.
// Pay attention when using pointer types as modifications to the value directly could lead to concurrency issues.
type Value[V any] interface {
	// Load returns the value stored, or zero value if no
	// value is present.
	// The ok result indicates whether value was set.
	Load() (v V, ok bool)
	// Store sets the value.
	Store(value V)
	// LoadOrStore returns the existing value if present.
	// Otherwise, it stores and returns the given value.
	// The loaded result is true if the value was loaded, false if stored.
	LoadOrStore(value V) (actual V, loaded bool)
	// Swap swaps the value for a key and returns the previous value if any.
	// The loaded result reports whether the key was present.
	Swap(value V) (previous V, loaded bool)
	// CompareAndSwap swaps the old and new values
	// if the value stored in the map is equal to old.
	//
	// Returns true if the swap was performed.
	//
	// ! this function uses reflect.DeepEqual to compare the values.
	CompareAndSwap(old, new V) bool
	// return true if the value is a zero value (not set)
	IsZero() bool
	// Exclusive executes the function f exclusively, ensuring that no other goroutine is accessing the value.
	//
	// The function f is passed the current value and a boolean indicating whether the value is set.
	// The function f should return the new value to be stored.
	//
	// ! Do not invoke any Value or Numeric functions within 'f' to prevent a deadlock.
	Exclusive(f func(v V, ok bool) V) V
	// Clear removes the value from the store.
	Clear()
}

// NewValue returns a new Value.
func NewValue[V any]() Value[V] {
	return internal.NewValue[V]()
}

// NewWithValue returns a new Value, set to the specified value.
func NewWithValue[V any](v V) Value[V] {
	return internal.NewWithValue(v)
}

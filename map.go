package mutex

import "github.com/thetechpanda/mutex/internal"

// Map is a generic interface that provides a way to interact with the map.
// its interface is identical to sync.Map and so are function definition and behaviour.
type Map[K comparable, V any] interface {
	// Load returns the value stored in the map for a key, or nil if no
	// value is present.
	// The ok result indicates whether value was found in the map.
	Load(key K) (v V, ok bool)
	// Store sets the value for a key.
	Store(key K, value V)
	// LoadOrStore returns the existing value for the key if present.
	// Otherwise, it stores and returns the given value.
	// The loaded result is true if the value was loaded, false if stored.
	LoadOrStore(key K, value V) (actual V, loaded bool)
	// LoadAndDelete deletes the value for a key, returning the previous value if any.
	// The loaded result reports whether the key was present.
	LoadAndDelete(key K) (value V, loaded bool)
	// Delete deletes the value for a key.
	Delete(key K)
	// Swap swaps the value for a key and returns the previous value if any.
	// The loaded result reports whether the key was present.
	Swap(key K, value V) (previous V, loaded bool)
	// CompareAndSwap swaps the old and new values for key
	// if the value stored in the map is equal to old.
	//
	// Returns true if the swap was performed.
	//
	// ! this function uses reflect.DeepEqual to compare the values.
	CompareAndSwap(key K, old, new V) bool
	// CompareAndDelete deletes the entry for key if its value is equal to old.
	//
	// If there is no current value for key in the map, CompareAndDelete
	// returns false (even if the old value is the nil interface value).
	//
	// ! this function uses reflect.DeepEqual to compare the values.
	CompareAndDelete(key K, old V) (deleted bool)
	// Range calls f sequentially for each key and value present in the map.
	// If f returns false, range stops the iteration.
	//
	// Range does not necessarily correspond to any consistent snapshot of the Map's
	// contents: no key will be visited more than once, but if the value for any key
	// is stored or deleted concurrently (including by f), Range may reflect any
	// mapping for that key from any point during the Range call. Range does not
	// block other methods on the receiver; even f itself may call any method on m.
	//
	// Range may be O(N) with the number of elements in the map even if f returns
	// false after a constant number of calls.
	Range(f func(K, V) bool)
	// Update allows the caller to change the value associated with the key atomically guaranteeing that the value would not be changed by another goroutine during the operation.
	//
	// ! Do not invoke any Map functions within 'f' to prevent a deadlock.
	Update(key K, f func(V, bool) V)
	// UpdateRange is a thread-safe version of Range that locks the map for the duration of the iteration and allows for the modification of the values.
	// If f returns false, UpdateRange stops the iteration, without updating the corresponding value in the map.
	//
	// ! Do not invoke any Map functions within 'f' to prevent a deadlock.
	UpdateRange(f func(K, V) (V, bool))
	// Exclusive provides a way to perform  operations on the map ensuring that no other operation is performed on the map during the execution of the function.
	//
	// ! Do not invoke any Map functions within 'f' to prevent a deadlock.
	Exclusive(f func(m map[K]V))
	// Clear removes all items from the map.
	Clear()
	// Has returns true if the map contains the key.
	Has(key K) bool
	// Keys returns a slice of all the keys present in the map, an empty slice is returned if the map is empty.
	Keys() (keys []K)
	// Values returns a slice of all the values present in the map, an empty slice is returned if the map is empty.
	Values() (values []V)
	// Entries returns two slices, one containing all the keys and the other containing all the values present in the map.
	Entries() (keys []K, values []V)
	// Len returns the number of unique keys in the map.
	Len() (n int)
}

// NewMap returns an empty Mutex Map.
func NewMap[K comparable, V any]() Map[K, V] {
	return internal.NewMap[K, V](nil)
}

// NewMapWithValue returns a Mutex Map with the provided map.
// m is copied into the Mutex Map.
func NewMapWithValue[K comparable, V any](m map[K]V) Map[K, V] {
	return internal.NewMap(m)
}

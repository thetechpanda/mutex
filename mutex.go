// # Mutex: Value, Numeric, and Map using Generics
//
// Package mutex provides a collection of thread-safe data structures using generics in Go.
// It offers a Value type for lock-protected values, a Numeric type for thread-safe numeric operations,
// and a Map type for a concurrent map with type safety. These structures are designed to be easy to use,
// providing a simple and familiar interface similar to well known atomic.Value and sync.Map, but with added type safety
// and the flexibility of generics. The package aims to simplify concurrent programming by ensuring safe
// access to shared data and reducing the boilerplate code associated with mutexes.
package mutex

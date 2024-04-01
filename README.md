# Mutex: Value, Numeric, and Map using Generics

![Test](https://github.com/thetechpanda/mutex/actions/workflows/go.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/thetechpanda/mutex)](https://goreportcard.com/report/github.com/thetechpanda/mutex)
[![Go Reference](https://pkg.go.dev/badge/github.com/thetechpanda/mutex.svg)](https://pkg.go.dev/github.com/thetechpanda/mutex)
[![Release](https://img.shields.io/github/release/thetechpanda/mutex.svg?style=flat-square)](https://github.com/thetechpanda/mutex/releases)
![Dependencies](https://img.shields.io/badge/Go_Dependencies-_None_-green.svg)

Package mutex provides a collection of thread-safe data structures using generics in Go.

It offers a Value type for lock-protected values, a Numeric type for thread-safe numeric operations, and a Map type for a concurrent map with type safety.  These structures are designed to be easy to use, providing a simple and familiar interface similar to well known atomic.Value and sync.Map, but with added type safety and the flexibility of generics. 

The package aims to simplify concurrent programming by ensuring safe access to shared data and reducing the boilerplate code associated with mutexes.

- `Value`, `Numeric` implements a simple thread-safe Value store that behaves similarly to `atomic.Value` but uses `sync.RWMutex` instead.
- `Numeric` extends `Value` with the `Add(delta V) V` function to simplify thread-safe counters.
- `Map` implements a simple thread-safe map that behaves similarly to `sync.Map` adding type safety and making it simple to know how many unique keys are in the map. 

## Data Types

`Value` and `Numeric` use generics, in general it is recommended to use them with built-in types, and they work best if you don't use of pointers (see below for more information about pointer values).

This means that when initialising a `Value` or `Numeric` a new variable will be created to hold the protected value.

For these reasons `NewMapWithValue` copies the map to the lock-protected map.

Please note that `Map`, `Value`, `Numeric` use `reflect.DeepEqual` in comparisons.

## About `Load()` and `Exclusive()`

When loading a value using `Load()` you have no guarantee that that value won't be changed when calling any Set method. This means that, similar to `atomic.Value`, you should not use `mutex.Value` or `mutex.Numeric` to develop a synchronisation algorithm unless you clearly understand the implications of doing so.
The same goes for `Map` values when using the `Load(K) V,bool` method.

When it is important to maintain consistency between read and write use `Exclusive()` that locks the value while the function is executed. For details on `Exclusive` argument review the relative interface.

## Pointers Values, Maps and Slices

Consider the following when using `Value` and `Numeric` with pointer values, maps and slices:

* **Concurrent Modification:** If multiple goroutines modify the data pointed to by the same pointer without proper synchronisation, it can lead to race conditions and unpredictable behaviour.
* **Data Race:** Even if `Value` itself is thread-safe, the data pointed to by the values is not automatically protected. Accessing or modifying the data through pointers in concurrent goroutines can cause data races.

## Documentation

You can find the generated go doc [here](godoc.md).

## Key Features

`Map`, `Value` and `Numeric`:
* **Type Safety:** Uses generics to provide a type-safe, lock-protected access to values.
* **Thread-Safety:** Ensures safe concurrent access to the value through the use of a `sync.RWMutex`.
* **Atomic Updates:** Includes functions that allows for atomic modifications to values in the map.
* **Shortcut:** Use a simple zero-dependency package to avoid rewriting the same code around mutex value protection over and over.

`Map`:
* **Iteration:** Supports iterating over the map with the Range function, and provides methods to obtain slices of keys (Keys), values (Values), or both (Entries).
* **Map Size:** Offers a Len function to easily retrieve the number of items in the map.

## Motivation

During the implementation of [TypedMap](https://github.com/thetechpanda/typedmap) I originally used (or rather misused) `atomic.Value`, a more in-depth review of the code and articles on the subject, made it clear that my use case did not justify the use of `atomic`. Additionally, citing Go's `atomic` package:

> These functions require great care to be used correctly. Except for special, low-level applications, synchronisation is better done with channels or the facilities of the sync package. 

I did however liked how *easy* to use `atomic.Value` is despite the in-depth know how they required to be properly used, `Value` and `Numeric` aim to provide the same simplicity and with a similar interface. The same goes for `sync.Map` interface.

## Usage

You can find more examples in [example_test.go](example_test.go)

```go
package main

import (
	"fmt"

	"github.com/thetechpanda/mutex"
)

func main() {
	value_example()
	numeric_example()
	map_example()
}

func value_example() {
	m := mutex.NewValue[string]()
	value := "42"
	m.Store(value)
	v, ok := m.Load()
	fmt.Println("value =", v, ", ok =", ok) // value = 42 , ok = true
}

func numeric_example() {
	m := mutex.NewNumeric[int]()
	value := 42
	m.Store(value)
	v, ok := m.Load()
	fmt.Println("value =", v, ", ok =", ok) // value = 42 , ok = true
}

func map_example() {
	m := mutex.NewMap[string, int]()
	m.Store("key", 42)
	v, ok := m.Load("key")
	fmt.Println("value =", v, ", ok =", ok) // value = 42 , ok = true
}
```

## Code coverage

```
$ go test -cover ./...
ok      github.com/thetechpanda/mutex   0.119s  coverage: 100.0% of statements
ok      github.com/thetechpanda/mutex/internal  0.443s  coverage: 100.0% of statements
```

## Installation

```bash
go get github.com/thetechpanda/mutex
```

## Bibliography

Below is a list of articles and videos that I reviewed to conclude that atomic operations were not the appropriate solution for the problem I was addressing. Some of these links provide documentation on how atomic operations work, while others discuss the nuances of using atomic operations effectively and sensibly.

* [Atomic Danger](https://abseil.io/docs/cpp/atomic_danger)
* [Atomics and Concurrency in C++](https://www.freecodecamp.org/news/atomics-and-concurrency-in-cpp/)
* [CON07-C. Ensure that compound operations on shared variables are atomic](https://wiki.sei.cmu.edu/confluence/display/c/CON07-C.+Ensure+that+compound+operations+on+shared+variables+are+atomic)
* [CON08-C. Do not assume that a group of calls to independently atomic methods is atomic](https://wiki.sei.cmu.edu/confluence/display/c/CON08-C.+Do+not+assume+that+a+group+of+calls+to+independently+atomic+methods+is+atomic)
* [atomic Weapons: The C++ Memory Model and Modern Hardware](https://herbsutter.com/2013/02/11/atomic-weapons-the-c-memory-model-and-modern-hardware/)
* [CppCon 2014: Herb Sutter "Lock-Free Programming (or, Juggling Razor Blades), Part I"](https://youtu.be/c1gO9aB9nbs)
* [CppCon 2014: Herb Sutter "Lock-Free Programming (or, Juggling Razor Blades), Part II"](https://youtu.be/CmxkPChOcvw)
* [CppCon 2015: Fedor Pikus PART 1 "Live Lock-Free or Deadlock (Practical Lock-free Programming)"](https://youtu.be/lVBvHbJsg5Y)
* [CppCon 2015: Fedor Pikus PART 2 "Live Lock-Free or Deadlock (Practical Lock-free Programming)"](https://youtu.be/1obZeHnAwz4)
* [CppCon 2017: Fedor Pikus "C++ atomics, from basic to advanced. What do they really do?"](https://youtu.be/ZQFzMfHIxng)

### Why most of the bibliography is about C++?

Reading and studying atomic operations in C++ can be applicable to Go because both languages deal with concurrency and share similar challenges related to memory consistency and synchronization. Understanding atomic operations in C++ can provide a deeper insight into low-level concurrency mechanisms, which can be beneficial when working with Go's concurrency primitives and atomic package, despite the differences in syntax and language features.

### Disclaimer
The content provided in the above links is not written or curated by the author of this package. All rights and attributions belong to the original authors. The inclusion of these links does not imply any endorsement or ownership by the author of this package.

## Roadmap

- [ ] Mutex Slices
- [ ] Benchmarks

## Contributing

Contributions are welcome and very much appreciated! 

Feel free to open an issue or submit a pull request.

## License

The `mutex` package is released under the MIT License. See the [LICENSE](LICENSE) file for details.

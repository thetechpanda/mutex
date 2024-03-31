# mutex.Value & mutex.Numeric

[![Go Report Card](https://goreportcard.com/badge/github.com/thetechpanda/mutex)](https://goreportcard.com/report/github.com/thetechpanda/mutex)
[![Go Reference](https://pkg.go.dev/badge/github.com/thetechpanda/mutex.svg)](https://pkg.go.dev/github.com/thetechpanda/mutex)
[![Release](https://img.shields.io/github/release/thetechpanda/mutex.svg?style=flat-square)](https://github.com/thetechpanda/mutex/releases)

`Value` implements a simple thread-safe Value store that behaves similarly to `atomic.Value` but uses `sync.RWMutex` instead.

In addition it implements `Numeric` that extends `Value` with the `Add(delta V) V` function to simplify thread-safe counters.

## About `Load()` and `Exclusive()`

When loading a value using `Load()` you have no guarantee that that value won't be changed when calling any Set method. This means that as `atomic.Value` you should not use `mutex.Value` or `mutex.Numeric` to develop a synchronisation algorithm.

When it is important to maintain consistency between read and write use `Exclusive()` that locks the value while `func(v V, ok bool) V` is executed.

## Pointers Values

Consider the following when using `Value` and `Numeric` with pointer values:

* **Concurrent Modification:** If multiple goroutines modify the data pointed to by the same pointer without proper synchronization, it can lead to race conditions and unpredictable behaviour.
* **Data Race:** Even if `Value` itself is thread-safe, the data pointed to by the values is not automatically protected. Accessing or modifying the data through pointers in concurrent goroutines can cause data races.

## Documentation

You can find the generated go doc [here](godoc.txt).

## Key Features

* **Type Safety:** Uses generics to provide a type-safe values.
* **Thread Safety:** Ensures safe concurrent access to the value through the use of a `sync.RWMutex`.
* **Shortcut:** Use a simple zero-dependency package to avoid rewriting the same code around mutex value protection over and over.

## Motivation

During the implementation of [TypedMap](https://github.com/thetechpanda/typedmap) I originally used (or rather misused) `atomic.Value`, a more in-depth review of the code and articles on the subject, eg: [Atomic Danger](https://abseil.io/docs/cpp/atomic_danger), made it clear that my use case did not justify the use of `atomic`. Additionally, citing Go's `atomic` package:

> These functions require great care to be used correctly. Except for special, low-level applications, synchronization is better done with channels or the facilities of the sync package. 

I did however liked how *easy* to use `atomic.Value` is despite the in-depth know how they required to be properly used, `Value` and `Numeric` aim to provide the same simplicity and with a similar interface.

## Usage

```go
package main

import (
	"fmt"

	"github.com/thetechpanda/mutex"
)

func main() {
	m := mutex.NewValue[int]()
	// or using NewNumeric
	// m := mutex.NewNumeric[int]()
	value := 42
	m.Store(value)
	v, ok := m.Load()
	fmt.Println("value =", v, ", ok =", ok) // value = 42 , ok = true
}
```

## Code coverage

```
$ go test -cover ./...
ok      github.com/thetechpanda/mutex      0.137s  coverage: 100.0% of statements
```

## Installation

```bash
go get github.com/thetechpanda/mutex
```

## Contributing

Contributions are welcome and very much appreciated! 

Feel free to open an issue or submit a pull request.

## License

`Value` is released under the MIT License. See the [LICENSE](LICENSE) file for details.

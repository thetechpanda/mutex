package mutex_test

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/thetechpanda/mutex"
)

func ExampleValue() {
	// create a new Value with initial value 0
	m := mutex.NewValue[string]()
	m.Store("")

	// create a wait group to synchronize goroutines
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			m.Store(strconv.Itoa(i))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			m.Store(strconv.Itoa(-i))
		}
	}()

	// wait for goroutines to finish
	wg.Wait()

	// value is either "-999" or "999"
	v, ok := m.Load()
	fmt.Println("value =", v, ", ok =", ok) // "value = -999 , ok = true" or "value = 999 , ok = true"
}

func ExampleNumeric() {
	// create a new Value with initial value 0
	m := mutex.NewNumeric[int]()
	m.Store(0)

	// create a wait group to synchronize goroutines
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			m.Add(1)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			m.Add(-1)
		}
	}()

	// wait for goroutines to finish
	wg.Wait()

	v, ok := m.Load()
	fmt.Println("value =", v, ", ok =", ok) // value = 0 , ok = true
}

package internal_test

import (
	"fmt"
	"reflect"
	"sync"
	"testing"

	"github.com/thetechpanda/mutex"
	"github.com/thetechpanda/mutex/internal"
)

func TestValue(t *testing.T) {
	t.Run("test without value", func(t *testing.T) {
		m := internal.NewValue[string]()
		if _, ok := m.Load(); ok {
			t.Errorf("Load(): Expected value to be absent")
		}
		if swapped := m.CompareAndSwap("A", "B"); swapped {
			t.Errorf("CompareAndSwap(): Expected value be swapped")
		}
		m.Store("B")
		if swapped := m.CompareAndSwap("B", "C"); !swapped {
			t.Errorf("CompareAndSwap(): Expected value be swapped")
		}

		mn := internal.NewValue[int]()
		if _, ok := mn.Load(); ok {
			t.Errorf("Load(): Expected value to be absent")
		}
		if swapped := mn.CompareAndSwap(-1, 42); swapped {
			t.Errorf("CompareAndSwap(): Expected value not to be swapped")
		}
		mn.Store(42)
		if swapped := mn.CompareAndSwap(42, -1); !swapped {
			t.Errorf("CompareAndSwap(): Expected value be swapped")
		}
	})
	t.Run("test with value", func(t *testing.T) {
		m := mutex.NewWithValue("42")
		v, ok := m.Load()
		if !ok || v != "42" {
			t.Errorf("Load(): Expected value 42 , got value %s", v)
		}
		if m.IsZero() {
			t.Errorf("IsZero(): Expected value to be present")
		}
	})
}

func TestValueNonComparable(t *testing.T) {
	m := internal.NewValue[[]any]()
	if _, ok := m.Load(); ok {
		t.Errorf("Load(): Expected value to be absent")
	}
	if swapped := m.CompareAndSwap([]any{1, 2}, []any{1, 2, 3}); swapped {
		t.Errorf("CompareAndSwap(): Expected value not to be swapped")
	}
	m.Store([]any{1, 2})
	if swapped := m.CompareAndSwap([]any{1, 2}, []any{1, 2, 3}); !swapped {
		t.Errorf("CompareAndSwap(): Expected value be swapped")
	}
}

func TestValueLoad(t *testing.T) {
	m := internal.NewValue[int]()
	value := 42
	m.Store(value)
	v, ok := m.Load()
	if !ok || v != value {
		t.Errorf("Load(): Expected value %d, got value %d", value, v)
	}
}

func TestValueStore(t *testing.T) {
	m := internal.NewValue[int]()
	value := 42
	m.Store(value)
	v, ok := m.Load()
	if !ok || v != value {
		t.Errorf("Load(): Expected value %d, got value %d", value, v)
	}
}

func TestValueLoadOrStore(t *testing.T) {
	m := internal.NewValue[int]()
	value := 42
	actual, loaded := m.LoadOrStore(value)
	if loaded {
		t.Errorf("LoadOrStore(): Expected value not present")
	}
	if actual != value {
		t.Errorf("Expected value %d, got value %d", value, actual)
	}

	actual, loaded = m.LoadOrStore(43)
	if !loaded {
		t.Errorf("LoadOrStore(): Expected value to be present")
	}
	if actual != value {
		t.Errorf("Expected value %d, got value %d", value, actual)
	}
}

func TestValueSwap(t *testing.T) {
	m := internal.NewValue[int]()
	value := 42
	previous, loaded := m.Swap(value)
	if loaded {
		t.Errorf("Swap(): value present")
	}
	if previous != 0 {
		t.Errorf("Expected previous value 0, got %d", previous)
	}

	previous, loaded = m.Swap(43)
	if !loaded {
		t.Errorf("Swap(): Expected value to be loaded")
	}
	if previous != value {
		t.Errorf("Expected previous value %d, got %d", value, previous)
	}
}

func TestValueCompareAndSwap(t *testing.T) {
	m := internal.NewValue[int]()
	current, swap := 42, 43
	m.Store(current)
	swapped := m.CompareAndSwap(current, swap)
	if !swapped {
		t.Errorf("CompareAndSwap(): Expected value to be swapped")
	}
	v, ok := m.Load()
	if !ok || v != swap {
		t.Errorf("Load(): Expected value %d, got value %d", swap, v)
	}
}

func testValues[V any](t *testing.T, value V) {
	m := internal.NewValue[V]()

	if !m.IsZero() {
		t.Errorf("IsZero(): Expected value to be zero")
	}

	if _, loaded := m.LoadOrStore(value); loaded {
		t.Errorf("LoadOrStore(): Expected value not to be present")
	}

	if m.IsZero() {
		t.Errorf("IsZero(): Expected value not to be zero")
	}

	v, ok := m.Load()
	if !ok {
		t.Errorf("Load(): Expected value to be present")
	}
	if !reflect.DeepEqual(v, value) {
		t.Errorf("Load(): Expected value %v, got value %v", value, v)
	}

	if _, loaded := m.LoadOrStore(value); !loaded {
		t.Errorf("LoadOrStore(): Expected value to be stored")
	}
	if !m.CompareAndSwap(value, value) {
		t.Errorf("CompareAndSwap(): Expected value to be swapped")
	}
	if _, ok := m.Swap(value); !ok {
		t.Errorf("Swap(): Expected value to be swapped")
	}

}
func TestValueTypes(t *testing.T) {

	// Test values of different types
	var b bool = true
	var s string = "string"
	var i int = 1
	var f float64 = 1.1
	var c complex128 = 1 + 1i
	var r rune = 'r'
	var by byte = byte('b')
	var st = struct {
		A int
	}{A: 1}
	var sl []int = []int{1, 2, 3}

	// Pointer to values
	var bP *bool = &b
	var sP *string = &s
	var iP *int = &i
	var fP *float64 = &f
	var cP *complex128 = &c
	var rP *rune = &r
	var byP *byte = &by
	var stP *struct {
		A int
	} = &st
	var slP []*int = []*int{&sl[0], &sl[1], &sl[2]}

	// Nil pointers
	var bN *bool = nil
	var sN *string = nil
	var iN *int = nil
	var fN *float64 = nil
	var cN *complex128 = nil
	var rN *rune = nil
	var byN *byte = nil
	var stN *struct {
		A int
	} = nil
	var slN []int = nil

	testValues(t, b)
	testValues(t, s)
	testValues(t, i)
	testValues(t, f)
	testValues(t, c)
	testValues(t, r)
	testValues(t, by)
	testValues(t, st)
	testValues(t, sl)

	testValues(t, bP)
	testValues(t, sP)
	testValues(t, iP)
	testValues(t, fP)
	testValues(t, cP)
	testValues(t, rP)
	testValues(t, byP)
	testValues(t, stP)
	testValues(t, slP)

	testValues(t, bN)
	testValues(t, sN)
	testValues(t, iN)
	testValues(t, fN)
	testValues(t, cN)
	testValues(t, rN)
	testValues(t, byN)
	testValues(t, stN)
	testValues(t, slN)
}

func TestValueConcurrency(t *testing.T) {
	// Create a new Value with initial value 0
	m := internal.NewValue[int]()
	m.Store(0)

	// Create a wait group to synchronize goroutines
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			m.Exclusive(func(v int, ok bool) int {
				return v + 1
			})
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			m.Exclusive(func(v int, ok bool) int {
				return v - 1
			})

		}
	}()

	// Wait for goroutines to finish
	wg.Wait()

	// Check if the value is now the string "0"
	v, ok := m.Load()
	fmt.Println("value=", v, "ok=", ok)
}

func TestValueClear(t *testing.T) {
	m := internal.NewValue[int]()
	m.Store(42)
	m.Clear()
	if v, ok := m.Load(); ok {
		t.Errorf("Load(): Expected value to be absent")
	} else if v != 0 {
		t.Errorf("Load(): Expected value to be 0, got %d", v)
	}
}

func TestValueAny(t *testing.T) {
	m := internal.NewValue[any]()
	if m.CompareAndSwap(nil, 43) { // this fails as set flag is false
		t.Errorf("CompareAndSwap(): Expected value not to be swapped")
	}
	m.Store(nil)
	if !m.CompareAndSwap(nil, 42) {
		t.Errorf("CompareAndSwap(): Expected value to be swapped")
	}
}

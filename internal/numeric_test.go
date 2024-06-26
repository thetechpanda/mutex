package internal_test

import (
	"sync"
	"testing"

	"github.com/thetechpanda/mutex/internal"
)

func TestNumeric(t *testing.T) {
	// create a new Numeric with initial value 0
	for _, mv := range []any{internal.NewNumeric[int](), internal.NewNumericWithValue[int](0)} {
		m := mv.(*internal.Numeric[int])
		// test the Add method
		t.Run("Add", func(t *testing.T) {
			m.Store(0)

			// add 5 to the value
			m.Add(5)

			// check if the value is now 5
			if v, _ := m.Load(); v != 5 {
				t.Errorf("Expected value 5, got %v", v)
			}
		})

		// test concurrent access to the Numeric
		t.Run("ConcurrentAccess", func(t *testing.T) {
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

			// check if the value is now 0
			if v, _ := m.Load(); v != 0 {
				t.Errorf("Expected value 0, got %v", v)
			}
		})

		// test adding negative values
		t.Run("AddNegative", func(t *testing.T) {
			m.Store(0)

			// add -5 to the value
			m.Add(-5)

			// check if the value is now -5
			if v, _ := m.Load(); v != -5 {
				t.Errorf("Expected value -5, got %v", v)
			}
		})

	}

}

func TestNumericClear(t *testing.T) {
	m := internal.NewNumeric[int]()
	m.Store(42)
	m.Clear()
	if v, ok := m.Load(); ok {
		t.Errorf("Load(): Expected value to be absent")
	} else if v != 0 {
		t.Errorf("Load(): Expected value to be 0, got %d", v)
	}
}

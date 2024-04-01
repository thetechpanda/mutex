package mutex_test

import (
	"testing"

	"github.com/thetechpanda/mutex"
)

func TestValue(t *testing.T) {

	t.Run("new without value", func(t *testing.T) {
		mv := mutex.NewValue[string]()
		_, ok := mv.Load()
		if ok {
			t.Errorf("Expected ok to be false, got true")
		}
		mv.Store("42")
		v, ok := mv.Load()
		if !ok {
			t.Errorf("Expected ok to be true, got false")
		}
		if v != "42" {
			t.Errorf("Expected value to be 42, got %v", v)
		}
	})
	t.Run("new with value", func(t *testing.T) {
		mv := mutex.NewWithValue[string]("42")
		v, ok := mv.Load()
		if !ok {
			t.Errorf("Expected ok to be true, got false")
		}
		if v != "42" {
			t.Errorf("Expected value to be 42, got %v", v)
		}
	})
}

func TestNumeric(t *testing.T) {
	t.Run("new without value", func(t *testing.T) {
		mv := mutex.NewNumeric[int]()
		_, ok := mv.Load()
		if ok {
			t.Errorf("Expected ok to be false, got true")
		}
		mv.Store(42)
		v, ok := mv.Load()
		if !ok {
			t.Errorf("Expected ok to be true, got false")
		}
		if v != 42 {
			t.Errorf("Expected value to be 42, got %v", v)
		}
	})
	t.Run("new with value", func(t *testing.T) {
		mv := mutex.NewNumericWithValue(42)
		v, ok := mv.Load()
		if !ok {
			t.Errorf("Expected ok to be true, got false")
		}
		if v != 42 {
			t.Errorf("Expected value to be 42, got %v", v)
		}
		mv.Add(1)
		v, ok = mv.Load()
		if !ok {
			t.Errorf("Expected ok to be true, got false")
		}
		if v != 43 {
			t.Errorf("Expected value to be 43, got %v", v)
		}
	})

}

func TestMap(t *testing.T) {
	t.Run("new without value", func(t *testing.T) {
		mv := mutex.NewMap[string, string]()
		_, ok := mv.Load("key")
		if ok {
			t.Errorf("Expected ok to be false, got true")
		}
		mv.Store("key", "42")
		v, ok := mv.Load("key")
		if !ok {
			t.Errorf("Expected ok to be true, got false")
		}
		if v != "42" {
			t.Errorf("Expected value to be 42, got %v", v)
		}
	})

	t.Run("new with value", func(t *testing.T) {
		m := map[string]string{"key": "42"}
		mv := mutex.NewMapWithValue(m)
		v, ok := mv.Load("key")
		if !ok {
			t.Errorf("Expected ok to be true, got false")
		}
		if v != "42" {
			t.Errorf("Expected value to be 42, got %v", v)
		}
	})
}

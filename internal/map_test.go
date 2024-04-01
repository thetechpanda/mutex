package internal_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/thetechpanda/mutex/internal"
)

func TestNewMap(t *testing.T) {
	m := internal.NewMap[string, int](nil)
	if m.Len() != 0 {
		t.Errorf("Len(): Expected a new map, got map with length %d", m.Len())
	}
}
func TestNewWithMap(t *testing.T) {
	data := map[string]int{
		"key1": 42,
		"key2": 43,
	}
	m := internal.NewMap(data)
	v, ok := m.Load("key1")
	if !ok || v != 42 {
		t.Errorf("Load(): Expected value 42 for key %q, got value %d", "key1", v)
	}
	if !m.Has("key2") {
		t.Errorf("Has(): Expected key %q to be present", "key2")
	}

	v, ok = m.Load("key2")
	if !ok || v != 43 {
		t.Errorf("Load(): Expected value 43 for key %q, got value %d", "key2", v)
	}
}

func TestMapLoad(t *testing.T) {
	m := internal.NewMap[string, int](nil)
	key := "key"
	value := 42
	m.Store(key, value)
	v, ok := m.Load(key)
	if !ok || v != value {
		t.Errorf("Load(): Expected value %d for key %q, got value %d", value, key, v)
	}
}

func TestMapStore(t *testing.T) {
	m := internal.NewMap[string, int](nil)
	key := "key"
	value := 42
	m.Store(key, value)
	v, ok := m.Load(key)
	if !ok || v != value {
		t.Errorf("Load(): Expected value %d for key %q, got value %d", value, key, v)
	}
}

func TestMapLoadOrStore(t *testing.T) {
	m := internal.NewMap[string, int](nil)
	key := "key"
	value := 42
	actual, loaded := m.LoadOrStore(key, value)
	if loaded {
		t.Errorf("LoadOrStore(): Expected key %q to be stored", key)
	}
	if actual != value {
		t.Errorf("Expected value %d for key %q, got value %d", value, key, actual)
	}

	actual, loaded = m.LoadOrStore(key, 43)
	if !loaded {
		t.Errorf("LoadOrStore(): Expected key %q to be loaded", key)
	}
	if actual != value {
		t.Errorf("Expected value %d for key %q, got value %d", value, key, actual)
	}
}

func TestMapLoadAndDelete(t *testing.T) {
	m := internal.NewMap[string, int](nil)
	key := "key"
	m.Store(key, 42)
	v, deleted := m.LoadAndDelete(key)
	if !deleted {
		t.Errorf("LoadAndDelete(): Expected key %q to be deleted", key)
	}
	if _, ok := m.Load(key); ok {
		t.Errorf("Load(): Expected key %q to be deleted", key)
	} else if v != 42 {
		t.Errorf("Expected value 42, got %d", v)
	}
}

func TestMapDelete(t *testing.T) {
	m := internal.NewMap[string, int](nil)
	key := "key"
	value := 42
	m.Store(key, value)
	m.Delete(key)
	if _, ok := m.Load(key); ok {
		t.Errorf("Load(): Expected key %q to be deleted", key)
	}
}

func TestMapSwap(t *testing.T) {
	m := internal.NewMap[string, int](nil)
	key := "key"
	value := 42
	previous, loaded := m.Swap(key, value)
	if loaded {
		t.Errorf("Swap(): Key %q was present in an empty map", key)
	}
	if previous != 0 {
		t.Errorf("Expected previous value 0, got %d", previous)
	}

	previous, loaded = m.Swap(key, 43)
	if !loaded {
		t.Errorf("Swap(): Expected key %q to be loaded", key)
	}
	if previous != value {
		t.Errorf("Expected previous value %d, got %d", value, previous)
	}
}

func TestMapCompareAndSwap(t *testing.T) {
	m := internal.NewMap[string, int](nil)
	key := "key"
	current, swap := 42, 43
	m.Store(key, current)
	swapped := m.CompareAndSwap(key, current, swap)
	if !swapped {
		t.Errorf("CompareAndSwap(): Expected key %q to be swapped", key)
	}
	v, ok := m.Load(key)
	if !ok || v != swap {
		t.Errorf("Load(): Expected value %d for key %q, got value %d", swap, key, v)
	}
}

func TestMapCompareAndDelete(t *testing.T) {
	m := internal.NewMap[string, int](nil)
	key := "key"
	value := 42
	m.Store(key, value)
	deleted := m.CompareAndDelete(key, 43)
	if deleted {
		t.Errorf("CompareAndDelete(): Expected key %q to not be deleted", key)
	}
	deleted = m.CompareAndDelete(key, value)
	if !deleted {
		t.Errorf("CompareAndDelete(): Expected key %q to be deleted", key)
	}
	actual, ok := m.Load(key)
	if ok {
		t.Errorf("Load(): Expected key %q to be deleted", key)
	}
	if actual != 0 {
		t.Errorf("Expected value 0, got %d", actual)
	}
}

func TestMapEntries(t *testing.T) {
	m := internal.NewMap[int, int](nil)

	zK, zV := m.Entries()
	if len(zK) != 0 {
		t.Errorf("Entries(): Expected empty keys, got %v", zK)
	}
	if len(zV) != 0 {
		t.Errorf("Entries(): Expected empty values, got %v", zV)
	}

	for i := 0; i < 100; i++ {
		m.Store(i, i)
	}
	sumKeys := 0
	sumValues := 0
	keys, values := m.Entries()
	for i := 0; i < 100; i++ {
		sumKeys += keys[i]
		sumValues += values[i]
	}
	switch {
	case len(keys) != 100:
		t.Errorf("Keys(): Expected 100 keys, got %d", len(keys))
	case sumKeys != 4950:
		t.Errorf("Expected sum 4950, got %d", sumKeys)
	case len(values) != 100:
		t.Errorf("Values(): Expected 100 values, got %d", len(values))
	case sumValues != 4950:
		t.Errorf("Expected sum 4950, got %d", sumValues)
	}
}

func TestMapRange(t *testing.T) {
	m := internal.NewMap[int, int](nil)
	for i := 0; i < 100; i++ {
		m.Store(i, 1)
	}

	var sum int
	m.Range(func(key, value int) bool {
		sum += value
		return true
	})
	if sum != 100 {
		t.Errorf("Range(): Expected sum 4950, got %d", sum)
	}
	sum = 0
	m.Range(func(key, value int) bool {
		sum++
		return false
	})
	if sum != 1 {
		t.Errorf("Range(): Expected sum 1, got %d", sum)
	}
}

func TestMapLen(t *testing.T) {
	m := internal.NewMap[string, int](nil)
	if m.Len() != 0 {
		t.Errorf("Len(): Expected length 0, got %d", m.Len())
	}

	m.Store("key1", 42)
	m.Store("key2", 42)
	if m.Len() != 2 {
		t.Errorf("Len(): Expected length 2, got %d", m.Len())
	}
}

func TestMapUpdate(t *testing.T) {
	m := internal.NewMap[string, int](nil)
	var wg sync.WaitGroup
	n := 100
	loops := 10
	incrementKey := func() {
		for i := 0; i < loops; i++ {
			m.Update("key", func(value int, ok bool) int {
				if !ok {
					panic("Update(): Expected key to be present")
				}
				return value + 1
			})
		}
	}

	m.Store("key", 0)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			incrementKey()
			wg.Done()
		}()
	}

	wg.Wait()

	expectedValue := n * 10
	if value, _ := m.Load("key"); value != expectedValue {
		t.Errorf("Load(): Expected final value to be %d, got %d", expectedValue, value)
	}
}

func TestMapUpdateRange(t *testing.T) {
	m := internal.NewMap[int, int](nil)
	for i := 0; i < 100; i++ {
		m.Store(i, 1)
	}

	var sum int
	wg := sync.WaitGroup{}
	wg.Add(1)
	sig := make(chan any)
	go func() {
		<-sig
		defer wg.Done()
		// update range will be called 100 times
		// each time it will add 1 to the value
		// so the final sum should be 100 + 100
		m.UpdateRange(func(k, i int) (int, bool) {
			sum++
			return 2, true
		})
	}()
	sig <- nil
	wg.Wait()
	if sum != 100 {
		t.Errorf("UpdateRange(): Expected sum 100, got %d", sum)
	}
	sum = 0
	// here all values are 2, so the sum should be 200
	m.Range(func(key, value int) bool {
		sum += value
		return true
	})
	if sum != 200 {
		t.Errorf("Range(): Expected sum 200, got %d", sum)
	}
	sum = 0
	// here all values are 2, so the sum should be 200
	m.Range(func(key, value int) bool {
		sum++
		return false
	})
	if sum != 1 {
		t.Errorf("Range(): Expected sum 1, got %d", sum)
	}

	sum = 0
	var mapKey int
	m.UpdateRange(func(k, i int) (int, bool) {
		sum++
		mapKey = k
		return 0, false
	})

	if sum != 1 {
		t.Errorf("UpdateRange(): Expected sum 1, got %d", sum)
	}

	if v, ok := m.Load(mapKey); !ok {
		t.Errorf("Load(): Expected key %d to be present", mapKey)
	} else if v != 2 {
		t.Errorf("Load(): Expected value 2, got %d", v)
	}
}

func TestMapConcurrentAccessSet(t *testing.T) {
	m := internal.NewMap[int, int](nil)

	// Number of goroutines to spawn
	numGoroutines := 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	ctx, cancel := context.WithCancel(context.Background())
	// starts numGoroutines go routines, the i values is used for content
	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			<-ctx.Done()
			// starts numKeys go routines, the j values is used as key
			// so we will have 1000*100 Get, Set and Delete operations
			for j := 0; j < numGoroutines; j++ {
				// uses context done to have all goroutines start at the same time
				_, ok := m.Load(j)
				if !ok {
					m.Store(j, i*i)
				}
				m.Delete(j)
			}
		}(i)
	}
	cancel()
	wg.Wait()
	if m.Len() != 0 {
		t.Errorf("Len(): Expected length 0, got %d", m.Len())
	}
}

func TestMapConcurrentAccessUpdate(t *testing.T) {
	m := internal.NewMap[int, int](nil)

	// Number of goroutines to spawn
	numGoroutines := 100
	for i := 0; i < numGoroutines; i++ {
		m.Store(i, 0)
	}
	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	ctx, cancel := context.WithCancel(context.Background())
	// starts numGoroutines go routines, the i values is used for content
	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			// uses context done to have all goroutines start at the same time
			<-ctx.Done()
			for j := 0; j < numGoroutines; j++ {
				for j := 0; j < numGoroutines; j++ {
					m.Update(i, func(v int, ok bool) int {
						if !ok {
							panic("Expected key to be present")
						}
						return v + 1
					})
				}
			}
		}(i)
	}
	cancel()
	wg.Wait()
	if m.Len() != numGoroutines {
		t.Errorf("Len(): Expected length 100, got %d", m.Len())
	}
	m.Range(func(k, v int) bool {
		if v != numGoroutines*numGoroutines {
			panic(fmt.Errorf("Expected value %d, got %d", numGoroutines*numGoroutines, v))
		}
		return true
	})

}

func TestMapConcurrentAccessUpdateRange(t *testing.T) {
	m := internal.NewMap[int, int](nil)
	// Number of goroutines to spawn
	numGoroutines := 100
	for i := 0; i < numGoroutines; i++ {
		m.Store(i, 0)
	}

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	ctx, cancel := context.WithCancel(context.Background())
	// starts numGoroutines go routines, the i values is used for content
	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			// uses context done to have all goroutines start at the same time
			<-ctx.Done()
			defer wg.Done()
			for j := 0; j < numGoroutines; j++ {
				m.UpdateRange(func(k, v int) (int, bool) {
					return v + 1, true
				})
			}
		}(i)
	}
	cancel()
	wg.Wait()
	if m.Len() != numGoroutines {
		t.Errorf("Len(): Expected length 100, got %d", m.Len())
	}
	m.Range(func(k, v int) bool {
		if v != numGoroutines*numGoroutines {
			panic(fmt.Errorf("Expected value %d, got %d", numGoroutines*numGoroutines, v))
		}
		return true
	})
}

func TestMapConcurrentAccessRange(t *testing.T) {
	m := internal.NewMap[int, int](nil)

	// Number of goroutines to spawn
	numGoroutines := 100
	for i := 0; i < numGoroutines; i++ {
		m.Store(i, 0)
	}
	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	ctx, cancel := context.WithCancel(context.Background())
	var count int
	var mu sync.Mutex
	increment := func() {
		mu.Lock()
		defer mu.Unlock()
		count++
	}
	// starts numGoroutines go routines, the i values is used for content
	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			// uses context done to have all goroutines start at the same time
			<-ctx.Done()
			// starts numKeys go routines, the j values is used as key
			// so we will have 1000*100 Get, Set and Delete operations
			for j := 0; j < numGoroutines; j++ {
				m.Range(func(k, v int) bool {
					increment()
					return true
				})
			}
		}(i)
	}
	cancel()
	wg.Wait()
	expect := numGoroutines * numGoroutines * numGoroutines
	if count != expect {
		t.Errorf("Expected count %d, got %d", expect, count)
	}
}

func TestMapConcurrentAccessExclusive(t *testing.T) {
	m := internal.NewMap[int, int](nil)
	// Number of goroutines to spawn
	numGoroutines := 100
	for i := 0; i < numGoroutines; i++ {
		m.Store(i, 0)
	}

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	ctx, cancel := context.WithCancel(context.Background())
	// starts numGoroutines go routines, the i values is used for content
	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			// uses context done to have all goroutines start at the same time
			<-ctx.Done()
			defer wg.Done()
			for j := 0; j < numGoroutines; j++ {
				m.Exclusive(func(m map[int]int) {
					for k, v := range m {
						m[k] = v + 1
					}
				})
			}
		}(i)
	}
	cancel()
	wg.Wait()
	if m.Len() != numGoroutines {
		t.Errorf("Len(): Expected length 100, got %d", m.Len())
	}
	expect := numGoroutines * numGoroutines
	m.Range(func(k, v int) bool {
		if v != expect {
			panic(fmt.Errorf("Expected value %d, got %d", expect, v))
		}
		return true
	})
}

func TestMapNotComparableType(t *testing.T) {
	m := internal.NewMap[int, []int](nil)
	m.Store(1, []int{1, 2, 3})
	if !m.CompareAndDelete(1, []int{1, 2, 3}) {
		t.Errorf("CompareAndDelete(): should return true")
	}
	if m.CompareAndSwap(1, []int{1, 2, 3}, []int{1, 2, 3, 4}) {
		t.Errorf("CompareAndSwap(): should return false")
	}
}
func TestMapComparableType(t *testing.T) {
	m := internal.NewMap[int, int](nil)
	m.Store(1, 1)
	if !m.CompareAndDelete(1, 1) {
		t.Errorf("Expected to return true")
	}
	m.Store(1, 1)
	if !m.CompareAndSwap(1, 1, 2) {
		t.Errorf("Expected to return true")
	}

	if m.CompareAndSwap(1, 3, 4) {
		t.Errorf("Expected not to return true")
	}

}

func TestMapClear(t *testing.T) {
	m := internal.NewMap[int, int](nil)
	for i := 0; i < 100; i++ {
		m.Store(i, i)
	}
	if m.Len() != 100 {
		t.Errorf("Len(): Expected length 100, got %d", m.Len())
	}
	m.Clear()
	if m.Len() != 0 {
		t.Errorf("Len(): Expected length 0, got %d", m.Len())
	}
}

func TestMapKeysValues(t *testing.T) {
	m := internal.NewMap[int, int](nil)

	if len(m.Keys()) != 0 {
		t.Errorf("Keys(): Expected empty keys, got %v", m.Keys())
	}

	if len(m.Values()) != 0 {
		t.Errorf("Values(): Expected empty values, got %v", m.Values())
	}

	for i := 0; i < 100; i++ {
		m.Store(i, i)
	}
	sumKeys := 0
	keys := m.Keys()
	for _, key := range keys {
		sumKeys += key
	}
	if len(keys) != 100 {
		t.Errorf("Keys(): Expected 100 keys, got %d", len(keys))
	} else if sumKeys != 4950 {
		t.Errorf("Expected sum 4950, got %d", sumKeys)
	}

	sumValues := 0
	values := m.Values()
	for _, value := range values {
		sumValues += value
	}
	if len(values) != 100 {
		t.Errorf("Values(): Expected 100 values, got %d", len(values))
	}
	if sumValues != 4950 {
		t.Errorf("Expected sum 4950, got %d", sumValues)
	}
}

func testMap[K comparable, V comparable](t *testing.T, key K, value V) {
	m := internal.NewMap[K, V](nil)
	m.Store(key, value)
	v, ok := m.Load(key)
	if !ok {
		t.Errorf("Load(): Expected key %v to be present", key)
	}
	if v != value {
		t.Errorf("Load(): Expected value %v for key %v, got value %v", value, key, v)
	}

	m.Delete(key)

	if _, loaded := m.LoadOrStore(key, value); loaded {
		t.Errorf("LoadOrStore(): Expected key %v not to be stored", key)
	}
	if !m.CompareAndSwap(key, value, value) {
		t.Errorf("CompareAndSwap(): Expected key %v to be swapped", key)
	}
	if !m.CompareAndDelete(key, value) {
		t.Errorf("CompareAndDelete(): Expected key %v to be deleted", key)
	}
	m.CompareAndDelete(key, value)
	m.Swap(key, value)
	m.LoadAndDelete(key)

}

func TestMapValues(t *testing.T) {

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

	testMap(t, b, b)
	testMap(t, s, s)
	testMap(t, i, i)
	testMap(t, f, f)
	testMap(t, c, c)
	testMap(t, r, r)
	testMap(t, by, by)
	testMap(t, st, st)

	testMap(t, b, bP)
	testMap(t, s, sP)
	testMap(t, i, iP)
	testMap(t, f, fP)
	testMap(t, c, cP)
	testMap(t, r, rP)
	testMap(t, by, byP)
	testMap(t, st, stP)

	testMap(t, bP, b)
	testMap(t, sP, s)
	testMap(t, iP, i)
	testMap(t, fP, f)
	testMap(t, cP, c)
	testMap(t, rP, r)
	testMap(t, byP, by)
	testMap(t, stP, st)

	testMap(t, bN, b)
	testMap(t, sN, s)
	testMap(t, iN, i)
	testMap(t, fN, f)
	testMap(t, cN, c)
	testMap(t, rN, r)
	testMap(t, byN, by)
	testMap(t, stN, st)

	testMap(t, b, bN)
	testMap(t, s, sN)
	testMap(t, i, iN)
	testMap(t, f, fN)
	testMap(t, c, cN)
	testMap(t, r, rN)
	testMap(t, by, byN)
	testMap(t, st, stN)

	testMap(t, bP, bN)
	testMap(t, sP, sN)
	testMap(t, iP, iN)
	testMap(t, fP, fN)
	testMap(t, cP, cN)
	testMap(t, rP, rN)
	testMap(t, byP, byN)
	testMap(t, stP, stN)

}

package internal

type Numeric[V uint | uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64 | float32 | float64 | complex64 | complex128] struct {
	_ noCopy // go vet to alert when copying by value.
	*Value[V]
}

// NewNumeric returns a new Numeric.
func NewNumeric[V uint | uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64 | float32 | float64 | complex64 | complex128]() *Numeric[V] {
	return &Numeric[V]{Value: NewValue[V]()}
}

// NewNumericWithValue returns a new Numeric, set to the specified value.
func NewNumericWithValue[V uint | uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64 | float32 | float64 | complex64 | complex128](v V) *Numeric[V] {
	return &Numeric[V]{Value: NewWithValue(v)}
}

// Add is a shortcut to Exclusive that adds delta to the value stored when using a Numeric.
func (m *Numeric[V]) Add(delta V) V {
	return m.Exclusive(func(v V, ok bool) V {
		return v + delta
	})
}

package mutex

import "github.com/thetechpanda/mutex/internal"

// Numeric is an interface that extends Value with an Add method. The value stored must be a numeric type.
type Numeric[V uint | uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64 | float32 | float64 | complex64 | complex128] interface {
	Value[V]
	// Add adds delta to the value stored.
	Add(delta V) V
}

// NewNumeric returns a new Numeric.
func NewNumeric[V uint | uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64 | float32 | float64 | complex64 | complex128]() Numeric[V] {
	return &internal.Numeric[V]{Value: internal.NewValue[V]()}
}

// NewNumericWithValue returns a new Numeric, set to the specified value.
func NewNumericWithValue[V uint | uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64 | float32 | float64 | complex64 | complex128](v V) Numeric[V] {
	return &internal.Numeric[V]{Value: internal.NewWithValue(v)}
}

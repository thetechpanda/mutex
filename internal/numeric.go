package internal

type Numeric[V uint | uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64 | float32 | float64 | complex64 | complex128] struct {
	*Value[V]
}

func (m *Numeric[V]) Add(delta V) V {
	return m.Exclusive(func(v V, ok bool) V {
		return v + delta
	})
}

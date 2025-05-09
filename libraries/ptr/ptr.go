package ptr

func P[T any](v T) *T {
	return &v
}

func ValueOrNil[T comparable](v T) *T {
	var zero T
	if v == zero {
		return nil
	} else {
		return &v
	}
}

func ValueOrZero[T comparable](v T) T {
	var zero T
	if v == zero {
		return zero
	} else {
		return v
	}
}

func ShallowCopy[T any](v *T) *T {
	if v == nil {
		return nil
	} else {
		cpy := new(T)
		*cpy = *v
		return cpy
	}
}

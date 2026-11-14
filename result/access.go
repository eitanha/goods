package result

import "github.com/rotisserie/eris"

var ErrIndexOutOfBounds = eris.New("index out of bounds")

func AccessSlice[T any](s []T, index int) Result[T] {
	if index < 0 || index >= len(s) {
		return Err[T](ErrIndexOutOfBounds).Wrapf("index %d out of bounds for slice of length %d", index, len(s))
	}

	return Ok(s[index])
}

var ErrKeyNotFound = eris.New("key not found")

func AccessMap[K comparable, V any](m map[K]V, key K) Result[V] {
	v, ok := m[key]
	if !ok {
		return Err[V](ErrKeyNotFound).Wrapf("key %v not found in map", key)
	}

	return Ok(v)
}

var ErrTypeMismatch = eris.New("type mismatch")

func AccessAs[T, U any](i U) Result[T] {
	v, ok := any(i).(T)
	if !ok {
		return Err[T](ErrTypeMismatch).Wrapf("value of type %T, not the same underlying type as %T", i, v)
	}
	return Ok(v)
}

var ErrNilPointer = eris.New("pointer is nil")

func AccessPtr[T any](ptr *T) Result[T] {
	if ptr == nil {
		return Err[T](ErrNilPointer)
	}
	return Ok(*ptr)
}

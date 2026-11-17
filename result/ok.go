package result

import (
	"github.com/eitanha/goods/errorutils"
)

func Ok[T any](value T) Result[T] {
	return Result[T]{value: &value}
}

func OkIfNotNil[T any](value *T) Result[T] {
	return Result[T]{value: value}
}

func OkIfNotErr[T any](value T, err error) Result[T] {
	if err != nil {
		return Result[T]{err: err}
	}
	return Result[T]{value: &value}
}

func OkIfNotNilWithErr[T any](value *T, err error) Result[T] {
	if value != nil {
		return Result[T]{value: value}
	}
	return Result[T]{err: err}
}

func OkIfNotPanic[T any](f func() T) (r Result[T]) {
	defer func() {
		r.err = errorutils.FromPanic(recover())
	}()

	value := f()
	r.value = &value
	return
}

func (r Result[T]) IsOk() bool {
	return r.value != nil
}

func (r Result[T]) IsNotOk() bool {
	return r.value == nil
}

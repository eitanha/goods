package result

import "github.com/rotisserie/eris"

func Ok[T any](value T) Result[T] {
	return Result[T]{value: &value}
}

func OkIfNotNil[T any](value *T) Result[T] {
	return Result[T]{value: value}
}

func OkIfNotPanic[T any](f func() T) (r Result[T]) {
	defer func() {
		if v := recover(); v == nil {
		} else if err, ok := v.(error); ok {
			r.err = err
		} else {
			r.err = eris.Errorf("recovered from panic: %v", v)
		}
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

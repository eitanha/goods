package result

func (r Result[T]) OrError() (T, error) {
	if r.IsNotOk() {
		var zero T
		return zero, r.Err()
	}
	return *r.value, nil
}

func (r Result[T]) OrNil() *T {
	if r.IsNotOk() {
		return nil
	}
	return r.value
}

func (r Result[T]) OrPanic() T {
	if r.IsNotOk() {
		panic(r.Err())
	}
	return *r.value
}

func (r Result[T]) Or(default_ T) T {
	if r.IsNotOk() {
		return default_
	}
	return *r.value
}

func (r Result[T]) OrZero() T {
	if r.IsNotOk() {
		var zero T
		return zero
	}
	return *r.value
}

func (r Result[T]) OrElse(recover func(err error) T) T {
	if r.IsNotOk() {
		return recover(r.Err())
	}
	return *r.value
}

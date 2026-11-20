package result

import "github.com/rotisserie/eris"

func (r Result[T]) Do(f func(T)) Result[T] {
	if r.IsOk() {
		f(*r.value)
	}
	return r
}

var ErrFalsePredicate = eris.New("false predicate")

func (r Result[T]) Validate(predicate func(T) error) Result[T] {
	if r.IsOk() {
		if err := predicate(*r.value); err != nil {
			return Err[T](err)
		}
	}
	return r
}

func (r Result[T]) Filter(predicate func(T) bool) Result[T] {
	if r.IsOk() && !predicate(*r.value) {
		return Err[T](ErrFalsePredicate)
	}
	return r
}

func (r Result[T]) Map[U any](transform func(T) U) Result[U] {
	if r.IsNotOk() {
		return Err[U](r.err)
	}
	u := transform(*r.value)
	return Ok[U](u)
}

func (r Result[T]) FlatMap[U any](transform func(T) Result[U]) Result[U] {
	if r.IsNotOk() {
		return Err[U](r.err)
	}
	return transform(*r.value)
}

func (r Result[T]) MapErr(f func(error) error) Result[T] {
	if r.IsNotOk() {
		return Err[T](f(r.err))
	}
	return r
}

var ErrFalseCondition = eris.New("false condition")

func (r Result[T]) If(ok bool) Result[T] {
	if ok || r.IsNotOk() {
		return r
	}
	return Err[T](ErrFalseCondition)
}

func (r Result[T]) IfNotErr(err error) Result[T] {
	if err == nil {
		return r
	}
	return Err[T](err)
}

func (r Result[T]) Cast[U any]() Result[U] {
	return r.FlatMap(AccessAs[U])
}

func Flatten[T any](r Result[Result[T]]) Result[T] {
	if r.value != nil {
		return *r.value
	}

	return Err[T](r.err)
}

func (r Result[T]) Recover(f func(error) T) Result[T] {
	if r.IsNotOk() {
		return Ok(f(r.err))
	}
	return r
}

func (r Result[T]) TryRecover(f func(error) Result[T]) Result[T] {
	if r.IsNotOk() {
		return f(r.err)
	}
	return r
}

func (r Result[T]) Match[U any](ok func(T) U, err func(error) U) U {
	if r.IsOk() {
		return ok(*r.value)
	} else {
		return err(r.err)
	}
}

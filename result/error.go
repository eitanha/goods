package result

import (
	"errors"

	"github.com/rotisserie/eris"
)

var NoValueErr = eris.New("no value")

func Err[T any](err error) Result[T] {
	if err == nil {
		return Err[T](NoValueErr)
	}
	return Result[T]{err: err}
}

func Errm[T any](message string) Result[T] {
	return Err[T](eris.New(message))
}

func Errf[T any](message string, args ...any) Result[T] {
	return Err[T](eris.Errorf(message, args...))
}

func (r Result[T]) IsErr() bool {
	return r.IsNotOk()
}

func (r Result[T]) Err() error {
	if r.IsOk() {
		return nil
	}

	if r.err != nil {
		return r.err
	}

	return NoValueErr
}

func (r Result[T]) Wrap(message string) Result[T] {
	if r.IsNotOk() {
		r.err = eris.Wrap(r.Err(), message)
	}
	return r
}

func (r Result[T]) Wrapf(message string, args ...any) Result[T] {
	if r.IsNotOk() {
		r.err = eris.Wrapf(r.Err(), message, args...)
	}
	return r
}

func (r Result[T]) ErrJoin(errs ...error) Result[T] {
	if r.IsOk() {
		return r
	}
	errs = append(errs, r.Err())
	r.err = errors.Join(errs...)
	return r
}

func (r Result[T]) WithErr(err error) Result[T] {
	if r.IsNotOk() {
		r.err = err
	}

	return r
}

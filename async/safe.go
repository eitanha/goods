package async

import (
	"github.com/eitanha/goods/result"
	"github.com/eitanha/goods/void"
)

// Guard turns a Future[T] into a Result[T] that doesn't panic.
func Guard[T any](future Future[T]) Result[T] {
	return &guardedFuture[T]{future}
}

type guardedFuture[T any] struct {
	Future[T]
}

func (f *guardedFuture[T]) Await() result.Result[T] {
	return result.OkIfNotPanic(f.Future.Await)
}

func GuardResult[T any](future Result[T]) Result[T] {
	return &guardedResult[T]{future}
}

type guardedResult[T any] struct {
	Result[T]
}

func (r *guardedResult[T]) Await() result.Result[T] {
	return result.Flatten(result.OkIfNotPanic(r.Result.Await))
}

func GuardError(future Error) Error {
	return &guardedError[error]{future}
}

type guardedError[T any] struct {
	Error
}

func (r *guardedError[T]) Await() error {
	return result.OkIfNotPanic(r.Error.Await).Match(func(result error) error {
		return result
	}, func(err error) error {
		return err
	})
}

func GuardVoid(future Void) Error {
	return &guardedVoid{future}
}

type guardedVoid struct {
	future Void
}

func (r *guardedVoid) Await() error {
	return result.OkIfNotPanic(r.future.Await).Match(func(result void.Type) error {
		return nil
	}, func(err error) error {
		return err
	})
}

func (r *guardedVoid) Done() <-chan struct{} {
	return r.future.Done()
}

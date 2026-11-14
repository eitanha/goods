package await

import (
	"github.com/eitanha/goods/async"
	"github.com/eitanha/goods/result"
)

func Safely[T any](future async.Future[T]) result.Result[T] {
	return async.Guard(future).Await()
}

func ResultSafely[T any](future async.Result[T]) result.Result[T] {
	return async.GuardResult(future).Await()
}

func ErrorSafely(future async.Error) error {
	return async.GuardError(future).Await()
}

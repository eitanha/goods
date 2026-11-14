package async

import (
	"github.com/eitanha/goods/result"
	"github.com/eitanha/goods/void"
)

type Waitable interface {
	Done() <-chan struct{}
}

type Future[TResult any] interface {
	Waitable
	Await() TResult
}

type Result[T any] = Future[result.Result[T]]
type Void = Future[void.Type]
type Error = Future[error]

func IsDone(waitable Waitable) bool {
	select {
	case <-waitable.Done():
		return true
	default:
		return false
	}
}

package async

import (
	"github.com/eitanha/goods/result"
	"github.com/eitanha/goods/void"
)

func Go[TValue any](fn func() TValue) Future[TValue] {
	return GoFor(MakeTask[TValue](), fn)
}

type routineView[TValue any] struct {
	Future[TValue]
}

func (r routineView[TValue]) Safely() Result[TValue] {
	return Guard(r.Future)
}

func VoidGo(fn func()) Void {
	return Go(void.Func(fn))
}

func GoResult[TValue any](fn func() result.Result[TValue]) Result[TValue] {
	return GuardResult(Go(fn))
}

func GoError(fn func() error) Error {
	return GuardError(Go(fn))
}

func GoFor[TValue any](task Task[TValue], fn func() TValue) Task[TValue] {
	go func() {
		defer RecoverInto(task)
		task.Set(fn())
	}()
	return task
}

func GoAndRecover(fn func(), panicables ...Panicable) {
	go func() {
		defer RecoverInto(panicables...)
		fn()
	}()
}

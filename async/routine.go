package async

import (
	"github.com/eitanha/goods/result"
	"github.com/eitanha/goods/void"
)

type routineView[TValue any] struct {
	Task[TValue]
}

func Go[TValue any](fn func() TValue) routineView[TValue] {
	return GoFor(MakeTask[TValue](), fn)
}

func (r routineView[TValue]) Safely() Result[TValue] {
	return Guard(r)
}

type voidRoutineView struct {
	Task[void.Type]
}

func VoidGo(fn func()) voidRoutineView {
	return voidRoutineView{Go(void.Func(fn))}
}

func (r voidRoutineView) Safely() Error {
	return GuardVoid(r)
}

func GoResult[TValue any](fn func() result.Result[TValue]) Result[TValue] {
	return GuardResult(Go(fn))
}

func GoError(fn func() error) Error {
	return GuardError(Go(fn))
}

func GoFor[TValue any](task Task[TValue], fn func() TValue) routineView[TValue] {
	go func() {
		defer RecoverInto(task)
		task.Set(fn())
	}()
	return routineView[TValue]{task}
}

func GoAndRecover(fn func(), panicables ...Panicable) {
	go func() {
		defer RecoverInto(panicables...)
		fn()
	}()
}

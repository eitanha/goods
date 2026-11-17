package after

import (
	"github.com/eitanha/goods/async"
	"github.com/eitanha/goods/errorutils"
)

func Finishing[T any](futures async.Future[T]) afterFinishing[T] {
	return afterFinishing[T]{futures}
}

type afterFinishing[T any] struct {
	async.Future[T]
}

func (c afterFinishing[T]) Do(fn func(T)) afterFinishing[T] {
	return afterFinishing[T]{async.Go(func() T {
		value := c.Await()
		fn(value)
		return value
	})}
}

func (c afterFinishing[T]) Map[U any](fn func(T) U) afterFinishing[U] {
	return afterFinishing[U]{async.Go(func() U {
		return fn(c.Await())
	})}
}

func (c afterFinishing[T]) Catch(onPanic func(any)) afterFinishing[T] {
	return c.Recover(func(p any) (v T) {
		onPanic(p)
		return
	})
}

func (c afterFinishing[T]) Recover(recover_ func(any) T) afterFinishing[T] {
	task := async.MakeTask[T]()
	return afterFinishing[T]{async.GoFor(&recoveringTask[T]{Task: task, fn: recover_}, func() T {
		return c.Await()
	})}
}

func (c afterFinishing[T]) RecoverError(recover_ func(error) T) afterFinishing[T] {
	return c.Recover(func(p any) T {
		return recover_(errorutils.FromPanic(p))
	})
}

type recoveringTask[T any] struct {
	async.Task[T]
	fn func(any) T
}

func (t *recoveringTask[T]) Panic(p any) {
	defer async.RecoverInto(t.Task)
	t.Task.Set(t.fn(p))
}

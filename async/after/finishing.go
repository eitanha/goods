package after

import "github.com/eitanha/goods/async"

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

func (c afterFinishing[T]) Recover(recover func(any) T) afterFinishing[T] {
	task := async.MakeTask[T]()
	return afterFinishing[T]{async.GoFor(&recoveringTask[T]{Task: task, fn: recover}, func() T {
		return c.Await()
	})}
}

type recoveringTask[T any] struct {
	async.Task[T]
	fn func(any) T
}

func (t *recoveringTask[T]) Panic(p any) {
	defer async.RecoverInto(t.Task)
	t.Task.Set(t.fn(p))
}

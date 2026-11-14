package await

import (
	"github.com/eitanha/goods/async"
	"github.com/eitanha/goods/slices"
	"github.com/eitanha/goods/void"
)

type Group[T any] []async.Future[T]

func (group Group[T]) Add(futures ...async.Future[T]) Group[T] {
	return append(group, futures...)
}

func (group Group[T]) AllSequentially() async.Future[[]T] {
	return async.Go(func() []T {
		return AllSequentially(group...)
	})
}

func (group Group[T]) All() async.Future[[]T] {
	t := async.MakeTask[[]T]()

	wrappers := slices.From(group).Select(func(a async.Future[T]) async.Future[T] {
		wrapper := async.MakeTask[T]()
		async.GoAndRecover(func() {
			select {
			case <-t.Done():
				t.Await()
			case <-a.Done():
				wrapper.Set(a.Await())
			}
		}, wrapper, t)
		return wrapper
	}).Slice()

	return async.GoFor(t, func() []T {
		return AllSequentially(wrappers...)
	})
}

func (group Group[T]) First() async.Future[T] {
	task := async.MakeTask[T]()
	for _, awaitable := range group {
		async.GoFor(task, func() T {
			select {
			case <-task.Done():
				return task.Await()
			case <-awaitable.Done():
				return awaitable.Await()
			}
		})
	}
	return task
}

type DoneGroup []async.Waitable

func (group DoneGroup) Add(waitable async.Waitable) DoneGroup {
	return append(group, waitable)
}

func (group DoneGroup) All() async.Future[void.Type] {
	return async.Go(void.Func(func() {
		UntilAllDone(group...)
	}))
}

func AllSequentially[T any](futures ...async.Future[T]) []T {
	results := make([]T, len(futures))
	for i, a := range futures {
		results[i] = a.Await()
	}
	return results
}

func UntilAllDone(waitables ...async.Waitable) {
	for _, w := range waitables {
		<-w.Done()
	}
}

func All[T any](futures ...async.Future[T]) []T {
	return Group[T](futures).All().Await()
}

func ForFirstOf[T any](futures ...async.Future[T]) T {
	return Group[T](futures).First().Await()
}

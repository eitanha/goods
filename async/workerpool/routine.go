package workerpool

import (
	"github.com/eitanha/goods/async"
	"github.com/eitanha/goods/result"
)

func (workerpool View) Go[T any](job func() T) async.Result[T] {
	return async.Go(func() result.Result[T] {
		task := async.MakeTask[T]()
		if err := workerpool.Manager.Run(func() {
			defer async.RecoverInto(task)
			task.Set(job())
		}).Await(); err != nil {
			return result.Err[T](err)
		}
		return result.Ok(task.Await())
	})
}

func (workerpool View) GoResult[T any](job func() result.Result[T]) async.Result[T] {
	return async.Go(func() result.Result[T] {
		task := async.MakeTask[result.Result[T]]()
		if err := workerpool.Manager.Run(func() {
			defer async.RecoverInto(task)
			task.Set(job())
		}).Await(); err != nil {
			return result.Err[T](err)
		}
		return task.Await()
	})
}

func (workerpool View) GoError(job func() error) async.Error {
	return async.Go(func() error {
		task := async.MakeTask[error]()
		if err := workerpool.Manager.Run(func() {
			defer async.RecoverInto(task)
			task.Set(job())
		}).Await(); err != nil {
			return err
		}
		return task.Await()
	})
}

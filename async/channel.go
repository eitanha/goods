package async

import (
	"github.com/eitanha/goods/void"
)

func Send[T any](channel chan<- T, value T) Task[void.Type] {
	task := MakeTask[void.Type]()
	return GoFor(task, void.Func(func() {
		select {
		case <-task.Done():
			return
		case channel <- value:
		}
	}))
}

func Receive[T any](channel <-chan T) Task[T] {
	task := MakeTask[T]()
	return GoFor(task, func() (v T) {
		select {
		case <-task.Done():
		case value := <-channel:
			v = value
		}
		return
	})
}

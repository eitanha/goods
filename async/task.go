package async

import "sync"

type Task[TValue any] interface {
	Future[TValue]
	Panicable
	Set(value TValue)
}

type TaskView[TValue any] struct {
	Task[TValue]
}

func MakeTask[TValue any]() TaskView[TValue] {
	return TaskView[TValue]{&task[TValue]{
		done: make(chan struct{}),
	}}
}

type task[TValue any] struct {
	value TValue
	panic any
	done  chan struct{}
	once  sync.Once
}

func (t *task[TValue]) Set(value TValue) {
	t.once.Do(func() {
		defer close(t.done)
		t.value = value
	})
}

func (t *task[TValue]) Panic(p any) {
	if p == nil {
		return
	}
	t.once.Do(func() {
		defer close(t.done)
		t.panic = p
	})
}

func (t *task[TValue]) Await() TValue {
	<-t.Done()
	if t.panic != nil {
		panic(t.panic)
	}
	return t.value
}

func (t *task[TValue]) Done() <-chan struct{} {
	return t.done
}

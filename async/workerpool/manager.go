package workerpool

import (
	"context"

	"github.com/eitanha/goods/async"
	"github.com/eitanha/goods/async/workerpool/worker"
)

type Job func()

type Manager interface {
	// Run pushes a job to the workerpool queue and returns an awaitable that completes when the push is complete.
	Run(job Job) async.Error
}

type View struct {
	Manager
}

func Make(ctx context.Context, options ...Option) View {
	config := &WorkerPoolConfig{
		Workers:    16,
		BufferSize: 64,
	}

	for _, option := range options {
		option(config)
	}

	threadpool := &workerpool{
		ctx:     ctx,
		workers: config.Workers,
		jobs:    make(chan Job, config.BufferSize),
	}

	threadpool.startWorkers()

	return View{threadpool}
}

type workerpool struct {
	ctx     context.Context
	workers int
	jobs    chan Job
}

func (t *workerpool) Run(job Job) async.Error {
	return async.Go(func() error {
		select {
		case <-t.ctx.Done():
			return t.ctx.Err()
		case t.jobs <- job:
		}
		return nil
	})
}

func (t *workerpool) worker() {
	for {
		select {
		case <-t.ctx.Done():
			return
		case job := <-t.jobs:
			job()
		}
	}
}

func (t *workerpool) startWorkers() {
	for i := 0; i < t.workers; i++ {
		worker.Go(t.ctx, t.worker)
	}
}

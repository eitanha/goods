package worker

import (
	"context"
	"time"

	"github.com/eitanha/goods/async"
	"github.com/eitanha/goods/async/await"
	"github.com/eitanha/goods/void"
)

// Go runs a worker function until it succeeds with no panics
func Go(ctx context.Context, worker func(), options ...Option) async.Error {
	config := &KeepaliveConfig{
		RetryInterval: 1 * time.Second,
	}

	for _, option := range options {
		option(config)
	}

	return async.Go(func() error {
		for {
			task := async.Go(void.Func(worker))
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-task.Done():
			}

			err := await.Safely(task).Err()
			if err == nil {
				return nil
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(config.RetryInterval):
			}
		}
	})
}

type Option func(*KeepaliveConfig)

type KeepaliveConfig struct {
	RetryInterval time.Duration
}

func WaitBeforeRetry(wait time.Duration) Option {
	return func(config *KeepaliveConfig) {
		config.RetryInterval = wait
	}
}

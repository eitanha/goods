package workerpool

type Option func(*WorkerPoolConfig)

type WorkerPoolConfig struct {
	Workers    int
	BufferSize int
}

func Workers(workers int) Option {
	return func(config *WorkerPoolConfig) {
		config.Workers = workers
	}
}

func BufferSize(bufferSize int) Option {
	return func(config *WorkerPoolConfig) {
		config.BufferSize = bufferSize
	}
}

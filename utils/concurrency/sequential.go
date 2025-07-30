package concurrency

import "context"

// SequentialExecutor runs jobs one by one
func SequentialExecutor() TaskExecutor {
	return func(ctx context.Context, task TaskFunc) error {
		return task(ctx)
	}
}

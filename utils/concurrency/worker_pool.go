package concurrency

import "context"

type Job struct {
	ctx  context.Context
	task TaskFunc
}

func WorkerPoolExecutor(limit int) TaskExecutor {
	jobs := make(chan Job)

	// Start a fixed number of worker goroutines
	for range limit {
		go func() {
			for j := range jobs {
				// Respect context cancellation
				select {
				case <-j.ctx.Done():
					// skip the task if context cancelled
					continue
				default:
					_ = j.task(j.ctx) // optional: collect/log error
				}
			}
		}()
	}

	// Return the TaskExecutor
	return func(ctx context.Context, task TaskFunc) error {
		select {
		case jobs <- Job{ctx: ctx, task: task}:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

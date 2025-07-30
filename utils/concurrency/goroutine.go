package concurrency

import "context"

func GoroutineExecutor(limit int) TaskExecutor {
	sem := make(chan struct{}, limit)

	return func(ctx context.Context, task TaskFunc) error {
		select {
		case sem <- struct{}{}:
			go func() {
				defer func() { <-sem }()
				_ = task(ctx)
			}()
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

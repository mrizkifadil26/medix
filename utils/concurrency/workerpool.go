package concurrency

import "context"

func WorkerPoolExecutor(limit int) TaskExecutor {
	tasks := make(chan TaskFunc)
	ctxs := make(chan context.Context)

	for i := 0; i < limit; i++ {
		go func() {
			for {
				select {
				case task := <-tasks:
					ctx := <-ctxs
					_ = task(ctx)
				}
			}
		}()
	}

	return func(ctx context.Context, task TaskFunc) error {
		select {
		case tasks <- task:
			ctxs <- ctx
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

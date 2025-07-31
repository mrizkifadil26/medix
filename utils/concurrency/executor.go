package concurrency

import (
	"context"
	"fmt"
	"sync"
)

// SelectExecutor returns an Executor based on config (mode and limit)
func SelectExecutor(cfg Config) (TaskExecutor, error) {
	switch cfg.Mode {
	case ModeSequential:
		return SequentialExecutor(), nil
	case ModeGoroutine:
		if cfg.Limit <= 0 {
			return nil, fmt.Errorf("goroutine mode requires limit > 0")
		}

		return GoroutineExecutor(cfg.Limit), nil
	case ModeWorkerPool:
		if cfg.Limit <= 0 {
			return nil, fmt.Errorf("workerpool mode requires limit > 0")
		}

		return WorkerPoolExecutor(cfg.Limit), nil
	default:
		return nil, fmt.Errorf("unknown executor mode: %s", cfg.Mode)
	}
}

func MustExecutor(cfg Config) TaskExecutor {
	exec, err := SelectExecutor(cfg)
	if err != nil {
		return SequentialExecutor()
	}

	return exec
}

func FromTaskExecutor(exec TaskExecutor) BatchExecutor {
	return func(ctx context.Context, tasks []TaskFunc) error {
		var wg sync.WaitGroup
		errCh := make(chan error, len(tasks))
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		for _, task := range tasks {
			wg.Add(1)

			t := task // avoid closure capture issues
			err := exec(ctx, func(ctx context.Context) error {
				defer wg.Done()

				if err := t(ctx); err != nil {
					select {
					case errCh <- err:
						cancel() // propagate cancel to all workers
					default:
						// error already captured
					}
				}

				return nil
			})

			if err != nil {
				// Executor rejected task (e.g., context canceled or channel full)
				select {
				case errCh <- err:
					cancel()
				default:
				}
			}
		}

		wg.Wait()
		close(errCh)

		// Return the first error (or nil)
		for err := range errCh {
			return err
		}

		return nil
	}
}

package concurrency

import (
	"context"
	"log"
	"time"
)

func WithRetry(base TaskExecutor, attempts int) TaskExecutor {
	return func(ctx context.Context, task TaskFunc) error {
		var err error
		for i := 0; i < attempts; i++ {
			if err = task(ctx); err == nil {
				return nil
			}
		}
		return err
	}
}

func WithTimeout(base TaskExecutor, timeout time.Duration) TaskExecutor {
	return func(parent context.Context, task TaskFunc) error {
		ctx, cancel := context.WithTimeout(parent, timeout)
		defer cancel()
		return base(ctx, task)
	}
}

func WithLogger(base TaskExecutor) TaskExecutor {
	return func(ctx context.Context, task TaskFunc) error {
		log.Println("Task started")
		err := base(ctx, task)
		log.Printf("Task ended, err=%v", err)
		return err
	}
}

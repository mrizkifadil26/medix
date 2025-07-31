package concurrency

import "context"

type TaskFunc func(ctx context.Context) error
type TaskExecutor func(ctx context.Context, task TaskFunc) error
type BatchExecutor func(ctx context.Context, tasks []TaskFunc) error

type Mode string

const (
	ModeSequential Mode = "sequential"
	ModeGoroutine  Mode = "goroutine"
	ModeWorkerPool Mode = "workerpool"
)

type Config struct {
	Mode  Mode `json:"mode"`  // "sequential", "goroutine", or "workerpool"
	Limit int  `json:"limit"` // for goroutine or workerpool mode
}

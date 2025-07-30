package concurrency

import "fmt"

// SelectExecutor returns an Executor based on config (mode and limit)
func SelectExecutor(cfg Config) (Executor, error) {
	switch cfg.Mode {
	case ModeSequential:
		return SequentialExecutor, nil
	case ModeGoroutine:
		if cfg.Limit <= 0 {
			return nil, fmt.Errorf("goroutine mode requires limit > 0")
		}
		return ConcurrentExecutor(cfg.Limit), nil
	case ModeWorkerPool:
		if cfg.Limit <= 0 {
			return nil, fmt.Errorf("workerpool mode requires limit > 0")
		}
		return WorkerPoolExecutor(cfg.Limit), nil
	default:
		return nil, fmt.Errorf("unknown executor mode: %s", cfg.Mode)
	}
}

func MustExecutor(cfg Config) Executor {
	exec, err := SelectExecutor(cfg)
	if err != nil {
		return SequentialExecutor
	}

	return exec
}

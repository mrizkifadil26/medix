package scanner

import "github.com/mrizkifadil26/medix/utils/concurrency"

func SelectExecutor(concurrencyLimit int) (concurrency.TaskExecutor, error) {
	if concurrencyLimit <= 1 {
		return concurrency.SequentialExecutor(), nil
	}

	return concurrency.WorkerPoolExecutor(concurrencyLimit), nil
}

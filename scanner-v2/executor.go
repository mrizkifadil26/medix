package scannerV2

import "github.com/mrizkifadil26/medix/utils/concurrency"

func SelectExecutor(concurrencyLimit int) (concurrency.TaskExecutor, error) {
	if concurrencyLimit <= 1 {
		return concurrency.SequentialExecutor(), nil
	}

	return concurrency.GoroutineExecutor(concurrencyLimit), nil
}

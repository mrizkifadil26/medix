package concurrency

import "sync"

func WorkerPoolExecutor(limit int) Executor {
	return func(jobs []func()) {
		var wg sync.WaitGroup
		jobChan := make(chan func())

		// Start workers
		for i := 0; i < limit; i++ {
			go func() {
				for job := range jobChan {
					job()
					wg.Done()
				}
			}()
		}

		// Submit jobs
		for _, job := range jobs {
			wg.Add(1)
			jobChan <- job
		}

		close(jobChan)
		wg.Wait()
	}
}

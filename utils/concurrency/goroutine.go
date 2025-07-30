package concurrency

import "sync"

func ConcurrentExecutor(limit int) Executor {
	return func(jobs []func()) {
		var wg sync.WaitGroup
		sem := make(chan struct{}, limit)

		for _, job := range jobs {
			wg.Add(1)
			sem <- struct{}{}
			go func(job func()) {
				defer wg.Done()
				job()
				<-sem
			}(job)
		}

		wg.Wait()
	}
}

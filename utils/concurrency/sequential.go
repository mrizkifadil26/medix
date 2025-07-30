package concurrency

// SequentialExecutor runs jobs one by one
func SequentialExecutor(jobs []func()) {
	for _, job := range jobs {
		job()
	}
}

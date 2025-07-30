package concurrency

type Executor func(jobs []func())

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

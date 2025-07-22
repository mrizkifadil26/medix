package scanner

import (
	"sync"

	"github.com/schollz/progressbar/v3"
)

type ProgressTracker struct {
	bar  *progressbar.ProgressBar
	lock sync.Mutex
}

// NewProgress initializes a new progress bar if enabled.
func NewProgress(total int, enabled bool, desc string) *ProgressTracker {
	if !enabled {
		return nil
	}

	if desc == "" {
		desc = "Scanning"
	}

	return &ProgressTracker{
		bar: progressbar.NewOptions(total,
			progressbar.OptionSetDescription(desc),
			progressbar.OptionShowCount(),
			progressbar.OptionSetWidth(20),
			progressbar.OptionClearOnFinish(),
		),
	}
}

// Add increments the progress bar by n steps, thread-safe.
func (p *ProgressTracker) Add(n int) {
	if p == nil || p.bar == nil {
		return
	}

	p.lock.Lock()
	defer p.lock.Unlock()
	_ = p.bar.Add(n)
}

// Finish finalizes the bar display. Optional.
func (p *ProgressTracker) Finish() {
	if p == nil || p.bar == nil {
		return
	}

	_ = p.bar.Finish()
}

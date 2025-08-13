package scanner

import (
	"sync"

	"github.com/schollz/progressbar/v3"
)

type ProgressTracker struct {
	bar  *progressbar.ProgressBar
	mu   sync.Mutex
	done bool
}

// NewProgress initializes a new progress bar if enabled.
func NewProgressTracker(
	total int,
	enabled bool,
	desc string,
) *ProgressTracker {
	if !enabled {
		return nil
	}

	if desc == "" {
		desc = "Processing"
	}

	bar := progressbar.NewOptions(total,
		progressbar.OptionSetDescription(desc),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(20),
		// progressbar.OptionClearOnFinish(),
		progressbar.OptionFullWidth(),
	)

	return &ProgressTracker{bar: bar}
}

// Add increments the progress bar by n steps, thread-safe.
func (p *ProgressTracker) Increment(n int) {
	if p == nil || p.bar == nil || p.done {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	_ = p.bar.Add(n)
}

func (p *ProgressTracker) Set(value int) {
	if p == nil || p.bar == nil || p.done {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	_ = p.bar.Set(value)
}

// Finish finalizes the bar display. Optional.
func (p *ProgressTracker) Finish() {
	if p == nil || p.bar == nil || p.done {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	_ = p.bar.Finish()
	p.done = true
}

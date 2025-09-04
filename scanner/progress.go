package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Spinner struct {
	mu      sync.Mutex
	count   int
	desc    string
	current string
	done    bool
	stopCh  chan struct{}
}

// NewInlineSpinner starts a spinner that shows "desc n"
func NewSpinner(desc string) *Spinner {
	s := &Spinner{
		desc:   desc,
		stopCh: make(chan struct{}),
	}

	go s.loop()
	return s
}

func (s *Spinner) loop() {
	frames := []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}
	i := 0
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.mu.Lock()
			prefix := s.current
			// keep path short
			if len(prefix) > 50 {
				prefix = "..." + prefix[len(prefix)-47:]
			}

			fmt.Fprintf(os.Stderr, "\r%s %s %d files: %s", string(frames[i]), s.desc, s.count, prefix)
			s.mu.Unlock()
			i = (i + 1) % len(frames)
		}
	}
}

// Increment increases the count.
func (s *Spinner) Increment(path string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.count++
	s.current = filepath.ToSlash(path)
}

// Finish stops the spinner and prints final count.
func (s *Spinner) Finish() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.done {
		return
	}
	s.done = true
	close(s.stopCh)

	prefix := s.current
	if len(prefix) > 50 {
		prefix = "..." + prefix[len(prefix)-47:]
	}
	fmt.Fprintf(os.Stderr, "\r✓ %s %d files (last: %s)\n", s.desc, s.count, prefix)
}

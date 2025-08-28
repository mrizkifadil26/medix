package tmdb

import (
	"fmt"
	"sync/atomic"
)

type Progress struct {
	total   int32
	current int32
}

func (p *Progress) Inc(title string, errMessage string) {
	newVal := atomic.AddInt32(&p.current, 1)
	status := "✅"
	displayTitle := title
	if errMessage != "" {
		status = "❌"
		displayTitle = fmt.Sprintf("%s (%s)", title, errMessage)
	}

	percent := float64(newVal) / float64(p.total) * 100

	fmt.Printf("[%d/%d %.1f%%] %s %s\n",
		newVal, p.total, percent, status, displayTitle)
}

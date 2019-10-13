package scheduler

import (
	"fmt"
	"time"
)

// ScheduleJob starts scheduled job.
// To end job send empty struct (struct{}{}) on done chan
// Example duration: 500 * time.Millisecond
func ScheduleJob(job func(), d time.Duration) chan struct{} {
	ticker := time.NewTicker(d)
	done := make(chan struct{})

	go func(t *time.Ticker) {
		for {
			select {
			case <-done:
				ticker.Stop()
				fmt.Println("Ticker stopped")
				return
			case <-ticker.C:
				job()
			}
		}
	}(ticker)

	return done
}

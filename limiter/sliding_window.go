package limiter

import (
	"sync"
	"time"
)

type SlidingWindow struct {
	windowSize time.Duration
	maxReq     int
	requests   []time.Time
	mu         sync.Mutex
}

func NewSlidingWindow(windowSize time.Duration, maxReq int) *SlidingWindow {
	return &SlidingWindow{
		windowSize: windowSize,
		maxReq:     maxReq,
		requests:   []time.Time{},
	}
}

func (sw *SlidingWindow) Allow() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-sw.windowSize)

	// Remove old requests
	valid := []time.Time{}
	for _, t := range sw.requests {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	sw.requests = valid

	if len(sw.requests) >= sw.maxReq {
		return false
	}

	sw.requests = append(sw.requests, now)
	return true
}
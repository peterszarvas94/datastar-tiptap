package main

import (
	"sync"
	"time"
)

type rateLimiter struct {
	mu     sync.Mutex
	max    int
	window time.Duration
	count  int
	reset  time.Time
	active bool
}

func newRateLimiter(max int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		max:    max,
		window: window,
	}
}

func (rl *rateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if !rl.active || now.After(rl.reset) {
		rl.count = 1
		rl.reset = now.Add(rl.window)
		rl.active = true
		return true
	}

	if rl.count >= rl.max {
		return false
	}

	rl.count++
	return true
}

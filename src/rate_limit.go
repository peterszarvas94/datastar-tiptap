package main

import (
	"sync"
	"time"
)

type rateLimiter struct {
	mu      sync.Mutex
	max     int
	window  time.Duration
	entries map[string]*rateLimitEntry
}

type rateLimitEntry struct {
	count int
	reset time.Time
}

func newRateLimiter(max int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		max:     max,
		window:  window,
		entries: make(map[string]*rateLimitEntry),
	}
}

func (rl *rateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	entry, ok := rl.entries[key]
	if !ok || now.After(entry.reset) {
		rl.entries[key] = &rateLimitEntry{
			count: 1,
			reset: now.Add(rl.window),
		}
		return true
	}

	if entry.count >= rl.max {
		return false
	}

	entry.count++
	return true
}

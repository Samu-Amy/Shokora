package ratelimiter

import (
	"sync"
	"time"
)

// TODO: attiva anche rate limiter nginx

type FixedWindowLimiter struct {
	sync.RWMutex
	clients map[string]int
	limit   int
	window  time.Duration
}

func NewFixedWindowLimiter(limit int, window time.Duration) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		clients: make(map[string]int),
		limit:   limit,
		window:  window,
	}
}

// TODO: sistema questo codice (implementa un rate limiter migliore - magari Sliding Window ed usando channels o sync.Map invece che Mutex)
func (rLimiter *FixedWindowLimiter) Allow(ip string) (bool, time.Duration) {
	rLimiter.RLock()
	count, exists := rLimiter.clients[ip]
	rLimiter.RUnlock()

	if !exists || count < rLimiter.limit {
		rLimiter.Lock()
		if !exists {
			go rLimiter.resetCount(ip) // TODO: perché elimina se non esiste?
		}

		rLimiter.clients[ip]++
		rLimiter.Unlock()
		return true, 0
	}

	return false, rLimiter.window
}

func (rLimiter *FixedWindowLimiter) resetCount(ip string) {
	time.Sleep(rLimiter.window)
	rLimiter.Lock()
	delete(rLimiter.clients, ip)
	rLimiter.Unlock()
}

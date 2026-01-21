package ratelimiter

import (
	"sync"
	"time"
)

// TODO: con redis (invece che in-memory con locks) sarebbe più veloce

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

// TODO: sistema questo codice (implementa un rate limiter migliore)
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

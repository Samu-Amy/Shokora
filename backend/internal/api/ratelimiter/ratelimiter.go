package ratelimiter

import "time"

type RateLimiterI interface {
	Allow(ip string) (bool, time.Duration)
}

type RateLimiterConfig struct {
	RequestsPerTimeFrame int
	TimeFrame            time.Duration
	Enabled              bool
}

// TODO: sfruttare il rate limiter anche per deboucing (evitare richieste troppo vicine tra loro, soprattutto se uguali o simili (modificano la stessa risorsa))

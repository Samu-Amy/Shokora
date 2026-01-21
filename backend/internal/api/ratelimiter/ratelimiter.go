package ratelimiter

import "time"

type RateLimiter interface {
	Allow(ip string) (bool, time.Duration)
}

type RateLimiterConfig struct {
	RequestsPerTimeFrame int
	TimeFrame            time.Duration
	Enabled              bool
}

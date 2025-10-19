package services

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter provides rate limiting functionality
type RateLimiter struct {
	limiter *rate.Limiter
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerMinute int) *RateLimiter {
	// Convert requests per minute to requests per second
	requestsPerSecond := float64(requestsPerMinute) / 60.0

	return &RateLimiter{
		limiter: rate.NewLimiter(rate.Limit(requestsPerSecond), 1), // Allow 1 burst
	}
}

// Wait blocks until the rate limiter allows the request
func (rl *RateLimiter) Wait(ctx context.Context) error {
	return rl.limiter.Wait(ctx)
}

// Allow returns true if the request is allowed without waiting
func (rl *RateLimiter) Allow() bool {
	return rl.limiter.Allow()
}

// Reserve reserves a token and returns a reservation
func (rl *RateLimiter) Reserve() *rate.Reservation {
	return rl.limiter.Reserve()
}

// ReserveN reserves n tokens and returns a reservation
func (rl *RateLimiter) ReserveN(now time.Time, n int) *rate.Reservation {
	return rl.limiter.ReserveN(now, n)
}

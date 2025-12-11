package limiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	tokens   uint64
	fillRate float64
	capacity uint64
	lastTime time.Time
	mu       sync.Mutex
}

var _ RateLimiter = (*TokenBucket)(nil)

func NewTokenBucket(rate float64, burst uint64) RateLimiter {
	return &TokenBucket{
		tokens:   burst,
		fillRate: rate,
		capacity: burst,
		lastTime: time.Now(),
		mu:       sync.Mutex{},
	}
}

// Allow implements RateLimiter.
func (t *TokenBucket) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()

	timePassed := now.Sub(t.lastTime).Seconds()

	toAdd := uint64(timePassed * t.fillRate)

	// We need to refill the bucket with toAdd tokens but not exceed the capacity
	if toAdd > 0 {
		t.tokens = min(t.capacity, t.tokens+toAdd)
		t.lastTime = now
	}

	// If we have enough tokens to serve the requests
	if t.tokens > 0 {
		t.tokens -= 1
		return true
	}

	return false
}

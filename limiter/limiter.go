// Package limiter provides rate limiting implementations
// like Sliding Window and Token Bucket
package limiter

type RateLimiter interface {
	Allow() bool
}

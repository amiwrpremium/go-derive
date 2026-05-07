// Package transport — see http.go for the overview.
package transport

import (
	"context"
	"sync"
	"time"
)

// RateLimiter is a token-bucket rate limiter.
//
// Derive's documented per-IP limits are 10 TPS sustained with a 5× burst;
// the SDK installs a [RateLimiter] with those defaults on every transport.
// Construct via [NewRateLimiter] (the zero value is not usable).
//
// A nil *RateLimiter is treated as "limiting disabled" — every operation
// is a no-op. This makes wiring optional limiters easy:
//
//	limiter := NewRateLimiter(0, 0) // returns nil
//	limiter.Wait(ctx)              // no-op, no panic
type RateLimiter struct {
	mu         sync.Mutex
	tokens     float64
	max        float64
	refillRate float64 // tokens per second
	last       time.Time
}

// NewRateLimiter returns a token-bucket limiter with the given sustained
// rate (tokens per second) and burst multiplier (capacity = tps × burst).
//
// Returns nil when tps <= 0, signalling "no limiting". burst <= 0 is
// clamped to 1 so the bucket always has at least 1 token of capacity.
func NewRateLimiter(tps float64, burst float64) *RateLimiter {
	if tps <= 0 {
		return nil
	}
	if burst <= 0 {
		burst = 1
	}
	return &RateLimiter{
		max:        tps * burst,
		refillRate: tps,
		tokens:     tps * burst,
		last:       time.Now(),
	}
}

// Wait blocks until a token is available, then consumes one.
//
// Returns ctx.Err() if ctx is cancelled before a token can be acquired.
// Calling Wait on a nil receiver is a no-op (returns nil immediately) — see
// the type doc for the rationale.
func (r *RateLimiter) Wait(ctx context.Context) error {
	if r == nil {
		return nil
	}
	for {
		r.mu.Lock()
		now := time.Now()
		elapsed := now.Sub(r.last).Seconds()
		r.last = now
		r.tokens += elapsed * r.refillRate
		if r.tokens > r.max {
			r.tokens = r.max
		}
		if r.tokens >= 1 {
			r.tokens--
			r.mu.Unlock()
			return nil
		}
		need := (1 - r.tokens) / r.refillRate
		r.mu.Unlock()

		t := time.NewTimer(time.Duration(need * float64(time.Second)))
		select {
		case <-ctx.Done():
			t.Stop()
			return ctx.Err()
		case <-t.C:
		}
	}
}

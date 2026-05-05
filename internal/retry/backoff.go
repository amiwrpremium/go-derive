// Package retry contains backoff helpers used by the WebSocket reconnect
// loop. Kept tiny and dependency-free so it can be reused anywhere.
package retry

import (
	"math/rand/v2"
	"time"
)

// Backoff is an exponential backoff with jitter. The Backoff struct
// itself is not safe for concurrent use because of the `current` field
// — give each retrying goroutine its own. The jitter source uses
// math/rand/v2's goroutine-safe global PRNG.
type Backoff struct {
	Initial time.Duration
	Max     time.Duration
	Factor  float64
	Jitter  float64

	current time.Duration
}

// NewBackoff constructs a Backoff with sensible defaults (initial 500ms,
// max 30s, factor 2, jitter 0.2).
func NewBackoff() *Backoff {
	return &Backoff{
		Initial: 500 * time.Millisecond,
		Max:     30 * time.Second,
		Factor:  2.0,
		Jitter:  0.2,
	}
}

// Reset returns the backoff to its initial state.
func (b *Backoff) Reset() { b.current = 0 }

// Next returns the duration to sleep before the next attempt. It walks up
// exponentially, capped at Max, and applies +/- Jitter*current jitter.
func (b *Backoff) Next() time.Duration {
	if b.current == 0 {
		b.current = b.Initial
	} else {
		next := time.Duration(float64(b.current) * b.Factor)
		if next > b.Max {
			next = b.Max
		}
		b.current = next
	}
	if b.Jitter <= 0 {
		return b.current
	}
	delta := b.Jitter * float64(b.current)
	jitter := (rand.Float64()*2 - 1) * delta
	out := b.current + time.Duration(jitter)
	if out < 0 {
		out = 0
	}
	return out
}

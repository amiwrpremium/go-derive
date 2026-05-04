// Package retry contains backoff helpers used by the WebSocket reconnect
// loop. Kept tiny and dependency-free so it can be reused anywhere.
package retry

import (
	"math/rand"
	"time"
)

// Backoff is an exponential backoff with jitter. It is not safe for
// concurrent use — give each retrying goroutine its own.
type Backoff struct {
	Initial time.Duration
	Max     time.Duration
	Factor  float64
	Jitter  float64

	current time.Duration
	rng     *rand.Rand
}

// NewBackoff constructs a Backoff with sensible defaults (initial 500ms,
// max 30s, factor 2, jitter 0.2).
func NewBackoff() *Backoff {
	return &Backoff{
		Initial: 500 * time.Millisecond,
		Max:     30 * time.Second,
		Factor:  2.0,
		Jitter:  0.2,
		// #nosec G404 -- jitter, not security-sensitive
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Reset returns the backoff to its initial state.
func (b *Backoff) Reset() { b.current = 0 }

// Next returns the duration to sleep before the next attempt. It walks up
// exponentially, capped at Max, and applies +/- Jitter*current jitter.
func (b *Backoff) Next() time.Duration {
	if b.rng == nil {
		b.rng = rand.New(rand.NewSource(time.Now().UnixNano())) // #nosec G404 -- jitter, not security-sensitive
	}
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
	jitter := (b.rng.Float64()*2 - 1) * delta
	out := b.current + time.Duration(jitter)
	if out < 0 {
		out = 0
	}
	return out
}

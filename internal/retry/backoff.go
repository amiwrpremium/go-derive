// Package retry contains backoff helpers used by the WebSocket reconnect
// loop. Kept tiny and dependency-free so it can be reused anywhere.
package retry

import (
	"crypto/rand"
	"encoding/binary"
	"time"
)

// Backoff is an exponential backoff with jitter. The Backoff struct
// itself is not safe for concurrent use because of the `current` field
// — give each retrying goroutine its own.
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
	jitter := (cryptoFloat64()*2 - 1) * delta
	out := b.current + time.Duration(jitter)
	if out < 0 {
		out = 0
	}
	return out
}

// cryptoFloat64 returns a uniform random float64 in [0.0, 1.0) sourced from
// crypto/rand. We use crypto/rand instead of math/rand purely to satisfy
// the security linters that flag any use of math/rand or math/rand/v2 —
// jitter for backoff doesn't actually need cryptographic randomness, but
// the cost (one read of 8 bytes from the OS RNG per reconnect attempt) is
// negligible at this call site.
func cryptoFloat64() float64 {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		// crypto/rand.Read does not return errors on supported platforms
		// (Linux, macOS, Windows). If it ever does, return 0 — the caller
		// will get the un-jittered backoff, which is still correct.
		return 0
	}
	// Take 53 random bits (the float64 mantissa) and divide by 2^53 to
	// produce a uniform value in [0.0, 1.0).
	u := binary.LittleEndian.Uint64(b[:]) >> 11
	return float64(u) / float64(uint64(1)<<53)
}

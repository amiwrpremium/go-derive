package auth

import (
	"sync"
	"sync/atomic"
	"time"
)

// NonceGen produces strictly-increasing nonces for action signing.
//
// Derive requires nonces to be unique per subaccount across an action's
// lifetime. This generator returns millisecond-timestamp-based nonces in
// the upper bits combined with a 16-bit incrementing suffix in the lower
// bits, which gives both human-readable ordering (the timestamp prefix)
// and collision resistance for many actions in the same millisecond.
//
// The zero value is not usable; construct via [NewNonceGen].
type NonceGen struct {
	last atomic.Uint64
	mu   sync.Mutex
	rand uint16
}

// NewNonceGen returns a generator seeded from the current time.
//
// The returned generator is safe for concurrent use.
func NewNonceGen() *NonceGen {
	g := &NonceGen{}
	// Take only the lower 16 bits of UnixNano() as the counter seed —
	// the explicit mask documents intent so gosec G115 sees the
	// narrowing as deliberate.
	g.rand = uint16(time.Now().UnixNano() & 0xFFFF)
	return g
}

// Next returns the next nonce.
//
// Under contention the algorithm bumps to (prev + 1) so the
// strict-monotonic property holds even when many goroutines call Next in
// the same millisecond.
func (g *NonceGen) Next() uint64 {
	g.mu.Lock()
	g.rand++
	suffix := uint64(g.rand)
	g.mu.Unlock()

	for {
		ms := uint64(time.Now().UnixMilli())
		candidate := ms<<16 | suffix
		prev := g.last.Load()
		if candidate <= prev {
			candidate = prev + 1
		}
		if g.last.CompareAndSwap(prev, candidate) {
			return candidate
		}
	}
}

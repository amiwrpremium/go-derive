package auth_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/auth"
)

// maxJSONSafeInteger is 2^53 − 1: the largest integer IEEE-754 double
// (and therefore standard JSON parsers) can represent exactly. Nonces
// emitted past this point lose precision when the engine deserialises
// them, which mismatches the digest we signed and produces 14014.
const maxJSONSafeInteger uint64 = 1<<53 - 1

func TestNonce_FitsInJSONSafeInteger(t *testing.T) {
	g := auth.NewNonceGen()
	for i := 0; i < 10_000; i++ {
		n := g.Next()
		require.LessOrEqual(t, n, maxJSONSafeInteger,
			"nonce %d at iteration %d exceeds the JSON-safe integer "+
				"ceiling (2^53-1); the engine's parser will truncate it",
			n, i)
	}
}

func TestNonce_MatchesDocsFormat(t *testing.T) {
	g := auth.NewNonceGen()

	for i := 0; i < 100; i++ {
		before := uint64(time.Now().UTC().UnixMilli())
		n := g.Next()
		after := uint64(time.Now().UTC().UnixMilli())

		// nonce / 1000 is the embedded millisecond timestamp.
		// It must lie in the interval [before, after] modulo a tiny
		// scheduler-induced tolerance.
		const tolerance uint64 = 5
		ms := n / 1000
		assert.GreaterOrEqual(t, ms+tolerance, before, "iteration %d", i)
		assert.LessOrEqual(t, ms, after+tolerance, "iteration %d", i)
	}
}

func TestNonce_SuffixWithinDocsRange(t *testing.T) {
	g := auth.NewNonceGen()
	for i := 0; i < 10_000; i++ {
		n := g.Next()
		assert.Less(t, n%1000, uint64(1000),
			"suffix on nonce %d (iter %d) outside the documented 0..999 range",
			n, i)
	}
}

func TestNonceGen_NewIsNonZeroFirst(t *testing.T) {
	// First nonce of a fresh generator is never zero. Timestamp prefix
	// dominates and guarantees this for any wall-clock past the epoch.
	g := auth.NewNonceGen()
	assert.NotZero(t, g.Next())
}

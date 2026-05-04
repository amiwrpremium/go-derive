package retry_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/internal/retry"
)

func TestBackoff_GrowsExponentially(t *testing.T) {
	b := &retry.Backoff{
		Initial: 100 * time.Millisecond,
		Max:     1 * time.Second,
		Factor:  2.0,
		Jitter:  0, // deterministic
	}
	got := []time.Duration{b.Next(), b.Next(), b.Next(), b.Next(), b.Next()}
	want := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		400 * time.Millisecond,
		800 * time.Millisecond,
		1 * time.Second, // capped
	}
	assert.Equal(t, want, got)
}

func TestBackoff_Reset(t *testing.T) {
	b := &retry.Backoff{Initial: 50 * time.Millisecond, Max: 1 * time.Second, Factor: 2}
	_ = b.Next()
	_ = b.Next()
	b.Reset()
	assert.Equal(t, 50*time.Millisecond, b.Next())
}

func TestNewBackoff_Defaults(t *testing.T) {
	b := retry.NewBackoff()
	assert.Equal(t, 500*time.Millisecond, b.Initial)
	assert.Equal(t, 30*time.Second, b.Max)
	assert.Equal(t, 2.0, b.Factor)
	assert.Equal(t, 0.2, b.Jitter)

	// First Next() returns Initial under jitter; assert it's within
	// [Initial*(1-Jitter), Initial*(1+Jitter)].
	first := b.Next()
	assert.GreaterOrEqual(t, first, 400*time.Millisecond)
	assert.LessOrEqual(t, first, 600*time.Millisecond)
}

// TestBackoff_Next_LazyRNG covers the rng==nil branch — when callers
// allocate a Backoff via struct literal without going through NewBackoff.
func TestBackoff_Next_LazyRNG(t *testing.T) {
	b := &retry.Backoff{Initial: 10 * time.Millisecond, Max: time.Second, Factor: 2, Jitter: 0.5}
	d := b.Next()
	// With jitter 0.5 and current=10ms, output is in [5ms, 15ms].
	assert.GreaterOrEqual(t, d, 5*time.Millisecond)
	assert.LessOrEqual(t, d, 15*time.Millisecond)
}

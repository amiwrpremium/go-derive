package transport_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/internal/transport"
)

func TestRateLimiter_Disabled(t *testing.T) {
	rl := transport.NewRateLimiter(0, 1)
	assert.Nil(t, rl, "tps<=0 should disable the limiter")
	// Wait must be safe on a nil receiver.
	require.NoError(t, rl.Wait(context.Background()))
}

func TestRateLimiter_AllowsBurst(t *testing.T) {
	// 100 TPS with 5x burst = 500 tokens immediately available.
	rl := transport.NewRateLimiter(100, 5)
	require.NotNil(t, rl)
	start := time.Now()
	for i := 0; i < 50; i++ {
		require.NoError(t, rl.Wait(context.Background()))
	}
	// Burst should not have required any waiting.
	assert.Less(t, time.Since(start), 50*time.Millisecond)
}

func TestRateLimiter_BlocksWhenExhausted(t *testing.T) {
	rl := transport.NewRateLimiter(50, 1)
	for i := 0; i < 50; i++ {
		require.NoError(t, rl.Wait(context.Background()))
	}
	start := time.Now()
	require.NoError(t, rl.Wait(context.Background()))
	// 51st token must wait at least ~20ms (1 / 50 = 0.02s) at this rate.
	assert.GreaterOrEqual(t, time.Since(start), 10*time.Millisecond)
}

func TestRateLimiter_RespectsContextCancel(t *testing.T) {
	rl := transport.NewRateLimiter(1, 1)
	require.NoError(t, rl.Wait(context.Background())) // burn the only token

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	err := rl.Wait(ctx)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestRateLimiter_NegativeBurstClampedToOne(t *testing.T) {
	// Burst <= 0 should be clamped to 1 internally — must still permit at least one call.
	rl := transport.NewRateLimiter(10, -1)
	require.NotNil(t, rl)
	require.NoError(t, rl.Wait(context.Background()))
}

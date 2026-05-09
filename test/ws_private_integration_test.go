//go:build integration

package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/pkg/channels/private"
	"github.com/amiwrpremium/go-derive/pkg/ws"
)

func TestWS_Login(t *testing.T) {
	env := loadEnv(t)
	c := newAuthWSClient(t, env) // Login already happened in newAuthWSClient.
	assert.True(t, c.IsConnected())
}

func TestWS_OrdersSubscribe(t *testing.T) {
	env := loadEnv(t)
	c := newAuthWSClient(t, env)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sub, err := ws.Subscribe[[]derive.Order](ctx, c, private.Orders{SubaccountID: env.subaccount})
	require.NoError(t, err)
	defer sub.Close()

	// Orders channel sends only on lifecycle events; subscribing without
	// error is the success criterion.
	assert.Equal(t, "subaccount."+itoa(env.subaccount)+".orders", sub.Channel())
}

func TestWS_BalancesSubscribe(t *testing.T) {
	env := loadEnv(t)
	c := newAuthWSClient(t, env)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	sub, err := ws.Subscribe[derive.Balance](ctx, c, private.Balances{SubaccountID: env.subaccount})
	require.NoError(t, err)
	defer sub.Close()

	// Balances often emit an initial snapshot.
	select {
	case bal, ok := <-sub.Updates():
		require.True(t, ok)
		assert.Equal(t, env.subaccount, bal.SubaccountID)
	case <-ctx.Done():
		t.Skip("no balance snapshot within 15s — channel may only emit on changes")
	}
}

// Derive does not expose a `subaccount.{id}.positions` subscription
// channel (verified against `derivexyz/cockpit`'s canonical schema list).
// Position state must be polled via `private/get_positions` or derived
// from the trades feed.

// itoa is a tiny strconv.Itoa wrapper for int64 to avoid importing strconv
// inline in every check.
func itoa(n int64) string {
	const digits = "0123456789"
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = digits[n%10]
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

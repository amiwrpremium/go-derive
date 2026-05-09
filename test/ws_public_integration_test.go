//go:build integration

package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive"
)

func TestWS_PublicConnect(t *testing.T) {
	env := loadEnv(t)
	c := newWSClient(t, env)
	assert.True(t, c.IsConnected())
}

func TestWS_OrderBookSubscribe(t *testing.T) {
	env := loadEnv(t)
	c := newWSClient(t, env)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	sub, err := derive.Subscribe[derive.OrderBook](ctx, c, derive.PublicOrderBook{
		Instrument: env.instrument, Depth: 5,
	})
	require.NoError(t, err)
	defer sub.Close()

	select {
	case ob, ok := <-sub.Updates():
		require.True(t, ok, "channel closed: %v", sub.Err())
		assert.Equal(t, env.instrument, ob.InstrumentName)
	case <-ctx.Done():
		t.Fatalf("no orderbook update within %v", 15*time.Second)
	}
}

func TestWS_TickerSubscribe(t *testing.T) {
	env := loadEnv(t)
	c := newWSClient(t, env)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sub, err := derive.Subscribe[derive.TickerSlim](ctx, c, derive.PublicTickerSlim{
		Instrument: env.instrument, Interval: "1000",
	})
	require.NoError(t, err)
	defer sub.Close()

	select {
	case tk, ok := <-sub.Updates():
		require.True(t, ok)
		assert.Equal(t, 1, tk.Ticker.MarkPrice.Sign())
	case <-ctx.Done():
		t.Fatalf("no ticker_slim update within %v", 10*time.Second)
	}
}

func TestWS_TradesSubscribe(t *testing.T) {
	env := loadEnv(t)
	c := newWSClient(t, env)

	// Trades channels only fire when prints happen — give them more time.
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	sub, err := derive.Subscribe[[]derive.Trade](ctx, c, derive.PublicTrades{Instrument: env.instrument})
	require.NoError(t, err)
	defer sub.Close()

	select {
	case trades, ok := <-sub.Updates():
		require.True(t, ok)
		// On testnet the channel may publish empty heartbeats — accept either.
		for _, tr := range trades {
			assert.Equal(t, env.instrument, tr.InstrumentName)
		}
	case <-ctx.Done():
		t.Skip("no trade prints within 60s — skipping (testnet may be quiet)")
	}
}

package ws_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/internal/testutil"
	"github.com/amiwrpremium/go-derive/pkg/channels/public"
	"github.com/amiwrpremium/go-derive/pkg/ws"
)

func TestSubscribe_TypedDelivery(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := ws.Subscribe[derive.OrderBook](context.Background(), c, public.OrderBook{Instrument: "BTC-PERP"})
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()

	require.True(t, srv.WaitSubscribed("orderbook.BTC-PERP.1.10", 1*time.Second))
	srv.Notify("orderbook.BTC-PERP.1.10", map[string]any{
		"instrument_name": "BTC-PERP",
		"bids":            [][]string{{"100", "1"}},
		"asks":            [][]string{{"101", "1"}},
		"timestamp":       1700000000000,
	})

	select {
	case ob, ok := <-sub.Updates():
		require.True(t, ok)
		assert.Equal(t, "BTC-PERP", ob.InstrumentName)
	case <-time.After(2 * time.Second):
		t.Fatal("update never delivered")
	}
}

func TestSubscribe_Channel_Method(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := ws.Subscribe[derive.OrderBook](context.Background(), c, public.OrderBook{Instrument: "X"})
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()

	assert.Equal(t, "orderbook.X.1.10", sub.Channel())
}

func TestSubscribe_Close_StopsUpdates(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := ws.Subscribe[derive.OrderBook](context.Background(), c, public.OrderBook{Instrument: "X"})
	require.NoError(t, err)
	require.NoError(t, sub.Close())
	// Second Close is harmless.
	require.NoError(t, sub.Close())
}

func TestSubscribe_TypeMismatch_Drops(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	// Subscribe with the wrong T (string) for an OrderBook channel.
	sub, err := ws.Subscribe[string](context.Background(), c, public.OrderBook{Instrument: "BTC-PERP"})
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()
	require.True(t, srv.WaitSubscribed("orderbook.BTC-PERP.1.10", 1*time.Second))

	srv.Notify("orderbook.BTC-PERP.1.10", map[string]any{
		"instrument_name": "BTC-PERP",
		"bids":            [][]string{}, "asks": [][]string{},
		"timestamp": 1700000000000,
	})

	select {
	case <-sub.Updates():
		t.Fatal("type mismatch should have prevented delivery")
	case <-time.After(150 * time.Millisecond):
		// expected
	}
}

func TestSubscribeFunc_DriverCallback(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	got := make(chan derive.OrderBook, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = ws.SubscribeFunc(ctx, c, public.OrderBook{Instrument: "ETH-PERP"}, func(ob derive.OrderBook) {
			got <- ob
		})
	}()

	require.True(t, srv.WaitSubscribed("orderbook.ETH-PERP.1.10", 1*time.Second))
	srv.Notify("orderbook.ETH-PERP.1.10", map[string]any{
		"instrument_name": "ETH-PERP",
		"bids":            [][]string{}, "asks": [][]string{},
		"timestamp": 1700000000000,
	})

	select {
	case ob := <-got:
		assert.Equal(t, "ETH-PERP", ob.InstrumentName)
	case <-time.After(2 * time.Second):
		t.Fatal("callback never fired")
	}
}

func TestSubscribeFunc_ContextCancelReturnsErr(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- ws.SubscribeFunc(ctx, c, public.OrderBook{Instrument: "X"}, func(derive.OrderBook) {})
	}()
	require.True(t, srv.WaitSubscribed("orderbook.X.1.10", 1*time.Second))
	cancel()

	select {
	case err := <-done:
		assert.True(t, errors.Is(err, context.Canceled))
	case <-time.After(2 * time.Second):
		t.Fatal("SubscribeFunc didn't return after cancel")
	}
}

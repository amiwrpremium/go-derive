package derive_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/internal/jsonrpc"
	"github.com/amiwrpremium/go-derive/internal/testutil"
)

func newWSTestClient(t *testing.T, srv *testutil.MockWSServer, signed bool) *derive.WsClient {
	t.Helper()
	cfg := derive.Testnet()
	cfg.WSURL = srv.URL()
	opts := []derive.Option{
		derive.WithCustomNetwork(cfg),
		derive.WithReconnect(false),
		derive.WithPingInterval(50 * time.Millisecond),
	}
	if signed {
		s, err := derive.NewLocalSigner(testKey)
		require.NoError(t, err)
		opts = append(opts, derive.WithSigner(s), derive.WithSubaccount(1))
	}
	c, err := derive.NewWsClient(opts...)
	require.NoError(t, err)
	return c
}

func TestWSClient_RequiresNetwork(t *testing.T) {
	_, err := derive.NewWsClient()
	assert.ErrorIs(t, err, derive.ErrInvalidConfig)
}

func TestWSClient_ConnectAndCall(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	srv.HandleResult("public/get_time", 1700000000000)

	c := newWSTestClient(t, srv, false)
	ctx := context.Background()
	require.NoError(t, c.Connect(ctx))
	defer func() { _ = c.Close() }()

	got, err := c.GetTime(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(1700000000000), got)
	assert.True(t, c.IsConnected())
}

func TestWSClient_Login(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	srv.HandleResult("public/login", nil)

	c := newWSTestClient(t, srv, true)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()
	require.NoError(t, c.Login(context.Background()))
}

func TestWSClient_Login_RequiresSigner(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSTestClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()
	err := c.Login(context.Background())
	assert.ErrorIs(t, err, derive.ErrUnauthorized)
}

func TestWSClient_Login_PropagatesAPIError(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	srv.Handle("public/login", func(json.RawMessage) (any, *jsonrpc.Error) {
		return nil, &jsonrpc.Error{Code: derive.CodeInvalidSignature, Message: "no"}
	})

	c := newWSTestClient(t, srv, true)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()
	err := c.Login(context.Background())
	assert.True(t, derive.Is(err, derive.ErrUnauthorized))
}

func TestWSClient_SubscribeTyped(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSTestClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := derive.Subscribe[derive.OrderBook](context.Background(), c, derive.PublicOrderBook{
		Instrument: "BTC-PERP", Group: "1", Depth: 5,
	})
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()

	require.True(t, srv.WaitSubscribed("orderbook.BTC-PERP.1.5", 1*time.Second))

	srv.Notify("orderbook.BTC-PERP.1.5", map[string]any{
		"instrument_name": "BTC-PERP",
		"bids":            [][]string{{"100", "1"}},
		"asks":            [][]string{{"101", "1"}},
		"timestamp":       1700000000000,
	})

	select {
	case ob, ok := <-sub.Updates():
		require.True(t, ok)
		assert.Equal(t, "BTC-PERP", ob.InstrumentName)
		assert.Equal(t, "orderbook.BTC-PERP.1.5", sub.Channel())
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for orderbook update")
	}
	assert.NoError(t, sub.Err())
}

func TestWSClient_SubscribeFunc(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSTestClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	got := make(chan derive.OrderBook, 1)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		_ = derive.SubscribeFunc(ctx, c, derive.PublicOrderBook{Instrument: "ETH-PERP"}, func(ob derive.OrderBook) {
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
	cancel()
}

func TestWSClient_SubscribeMethodTypeMismatch(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSTestClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := derive.Subscribe[string](context.Background(), c, derive.PublicOrderBook{Instrument: "BTC-PERP"})
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
	}
}

func TestWSClient_NetworkAccessor(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSTestClient(t, srv, false)
	defer func() { _ = c.Close() }()
	assert.Equal(t, derive.NetworkTestnet, c.Network().Network)
}

func TestWS_RequiresNetwork(t *testing.T) {
	_, err := derive.NewWsClient()
	assert.ErrorIs(t, err, derive.ErrInvalidConfig)
}

func TestWS_WithMainnet(t *testing.T) {
	c, err := derive.NewWsClient(derive.WithMainnet())
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
	assert.Equal(t, derive.NetworkMainnet, c.Network().Network)
}

func TestWS_WithTestnet(t *testing.T) {
	c, err := derive.NewWsClient(derive.WithTestnet())
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
	assert.Equal(t, derive.NetworkTestnet, c.Network().Network)
}

func TestWS_WithCustomNetwork(t *testing.T) {
	custom := derive.Testnet()
	custom.WSURL = "ws://example.invalid/ws"
	c, err := derive.NewWsClient(derive.WithCustomNetwork(custom))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
	assert.Equal(t, "ws://example.invalid/ws", c.Network().WSURL)
}

func TestWS_AllOptionsCompose(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()

	signer, err := derive.NewLocalSigner(testKey)
	require.NoError(t, err)

	cfg := derive.Testnet()
	cfg.WSURL = srv.URL()

	c, err := derive.NewWsClient(
		derive.WithCustomNetwork(cfg),
		derive.WithSigner(signer),
		derive.WithSubaccount(7),
		derive.WithUserAgent("custom/1"),
		derive.WithRateLimit(50, 2),
		derive.WithPingInterval(100*time.Millisecond),
		derive.WithReconnect(false),
		derive.WithSignatureExpiry(120),
	)
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
	assert.NotNil(t, c)
}

func TestWS_DefaultsApplied(t *testing.T) {
	c, err := derive.NewWsClient(derive.WithTestnet())
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
}

func TestSubscribe_TypedDelivery(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSTestClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := derive.Subscribe[derive.OrderBook](context.Background(), c, derive.PublicOrderBook{Instrument: "BTC-PERP"})
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
	c := newWSTestClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := derive.Subscribe[derive.OrderBook](context.Background(), c, derive.PublicOrderBook{Instrument: "X"})
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()

	assert.Equal(t, "orderbook.X.1.10", sub.Channel())
}

func TestSubscribe_Close_StopsUpdates(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSTestClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := derive.Subscribe[derive.OrderBook](context.Background(), c, derive.PublicOrderBook{Instrument: "X"})
	require.NoError(t, err)
	require.NoError(t, sub.Close())
	require.NoError(t, sub.Close())
}

func TestSubscribe_TypeMismatch_Drops(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSTestClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := derive.Subscribe[string](context.Background(), c, derive.PublicOrderBook{Instrument: "BTC-PERP"})
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
	}
}

func TestSubscribeFunc_DriverCallback(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSTestClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	got := make(chan derive.OrderBook, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = derive.SubscribeFunc(ctx, c, derive.PublicOrderBook{Instrument: "ETH-PERP"}, func(ob derive.OrderBook) {
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
	c := newWSTestClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- derive.SubscribeFunc(ctx, c, derive.PublicOrderBook{Instrument: "X"}, func(derive.OrderBook) {})
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

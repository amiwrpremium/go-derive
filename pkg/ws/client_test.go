package ws_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/internal/jsonrpc"
	"github.com/amiwrpremium/go-derive/internal/netconf"
	"github.com/amiwrpremium/go-derive/internal/testutil"
	"github.com/amiwrpremium/go-derive/pkg/channels/public"
	"github.com/amiwrpremium/go-derive/pkg/ws"
)

const testKey = "0x4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318"

func newWSClient(t *testing.T, srv *testutil.MockWSServer, signed bool) *ws.Client {
	t.Helper()
	cfg := netconf.Testnet()
	cfg.WSURL = srv.URL()
	opts := []ws.Option{ws.WithCustomNetwork(cfg), ws.WithReconnect(false), ws.WithPingInterval(50 * time.Millisecond)}
	if signed {
		s, err := derive.NewLocalSigner(testKey)
		require.NoError(t, err)
		opts = append(opts, ws.WithSigner(s), ws.WithSubaccount(1))
	}
	c, err := ws.New(opts...)
	require.NoError(t, err)
	return c
}

func TestWSClient_RequiresNetwork(t *testing.T) {
	_, err := ws.New()
	assert.ErrorIs(t, err, derive.ErrInvalidConfig)
}

func TestWSClient_ConnectAndCall(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	srv.HandleResult("public/get_time", 1700000000000)

	c := newWSClient(t, srv, false)
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

	c := newWSClient(t, srv, true)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()
	require.NoError(t, c.Login(context.Background()))
}

func TestWSClient_Login_RequiresSigner(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false) // no signer
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

	c := newWSClient(t, srv, true)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()
	err := c.Login(context.Background())
	assert.True(t, derive.Is(err, derive.ErrUnauthorized))
}

func TestWSClient_SubscribeTyped(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := ws.Subscribe[derive.OrderBook](context.Background(), c, public.OrderBook{
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
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	got := make(chan derive.OrderBook, 1)
	ctx, cancel := context.WithCancel(context.Background())

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
	cancel()
}

func TestWSClient_SubscribeMethodTypeMismatch(t *testing.T) {
	// Channel decodes into derive.OrderBook but the test asks for a different type;
	// the generic glue must reject the mismatch without crashing.
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := ws.Subscribe[string](context.Background(), c, public.OrderBook{Instrument: "BTC-PERP"})
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()
	require.True(t, srv.WaitSubscribed("orderbook.BTC-PERP.1.10", 1*time.Second))

	// Push an orderbook payload — the typed pump should drop it (mismatch),
	// not deliver it.
	srv.Notify("orderbook.BTC-PERP.1.10", map[string]any{
		"instrument_name": "BTC-PERP",
		"bids":            [][]string{}, "asks": [][]string{},
		"timestamp": 1700000000000,
	})
	select {
	case <-sub.Updates():
		t.Fatal("type mismatch should have prevented delivery")
	case <-time.After(150 * time.Millisecond):
		// Expected: nothing arrived.
	}
}

func TestWSClient_NetworkAccessor(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	defer func() { _ = c.Close() }()
	assert.Equal(t, netconf.NetworkTestnet, c.Network().Network)
}

// silence the "imported and not used" lint when building before tests run.
var _ = json.Marshal
var _ = jsonrpc.Version

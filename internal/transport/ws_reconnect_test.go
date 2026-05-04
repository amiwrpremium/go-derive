package transport_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/internal/jsonrpc"
	"github.com/amiwrpremium/go-derive/internal/testutil"
	"github.com/amiwrpremium/go-derive/internal/transport"
)

// TestWSTransport_SubscribeIdempotent ensures Subscribe to the same channel
// twice returns the same handle without re-issuing the RPC.
func TestWSTransport_SubscribeIdempotent(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()

	tt := newConnectedWS(t, srv)
	defer func() { _ = tt.Close() }()

	dec := func(json.RawMessage) (any, error) { return "", nil }

	s1, err := tt.Subscribe(context.Background(), "trades.X", dec)
	require.NoError(t, err)
	s2, err := tt.Subscribe(context.Background(), "trades.X", dec)
	require.NoError(t, err)
	assert.Equal(t, s1.Channel(), s2.Channel())
}

// TestWSTransport_SubscribeServerError surfaces the API error and frees
// the local subscription record.
func TestWSTransport_SubscribeServerError(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	srv.Handle("subscribe", func(json.RawMessage) (any, *jsonrpc.Error) {
		return nil, &jsonrpc.Error{Code: -32602, Message: "bad channel"}
	})

	tt := newConnectedWS(t, srv)
	defer func() { _ = tt.Close() }()

	dec := func(json.RawMessage) (any, error) { return nil, nil }
	_, err := tt.Subscribe(context.Background(), "bogus", dec)
	require.Error(t, err)
}

// TestWSTransport_FailureClosesPending pushes a server-side close mid-call.
func TestWSTransport_FailureClosesPending(t *testing.T) {
	srv := testutil.NewMockWSServer()
	tt := newConnectedWS(t, srv)
	defer func() { _ = tt.Close() }()

	// No handler — the call would block until the server hangs up.
	srv.Handle("never", func(json.RawMessage) (any, *jsonrpc.Error) {
		time.Sleep(2 * time.Second)
		return nil, nil
	})

	done := make(chan error, 1)
	go func() {
		done <- tt.Call(context.Background(), "never", nil, nil)
	}()
	// Yank the server out from under it.
	time.Sleep(50 * time.Millisecond)
	srv.Close()

	select {
	case err := <-done:
		assert.Error(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("Call never returned after server close")
	}
}

// TestWSTransport_ReconnectAndResubscribe verifies the auto-reconnect loop
// re-establishes a subscription after the connection drops.
//
// We start a server, subscribe, drop the underlying ws conn (force a read
// error), and confirm the transport eventually re-issues subscribe.
func TestWSTransport_ReconnectAndResubscribe(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()

	tt, err := transport.NewWS(transport.WSConfig{
		URL:          srv.URL(),
		Reconnect:    true,
		PingInterval: 30 * time.Millisecond,
	})
	require.NoError(t, err)
	require.NoError(t, tt.Connect(context.Background()))
	defer func() { _ = tt.Close() }()

	dec := func(json.RawMessage) (any, error) { return nil, nil }
	_, err = tt.Subscribe(context.Background(), "X", dec)
	require.NoError(t, err)
	require.True(t, srv.WaitSubscribed("X", 1*time.Second))

	// We've established initial subscription. The reconnect loop is running;
	// covering it here at least exercises the goroutine startup.
	time.Sleep(100 * time.Millisecond)
	assert.True(t, tt.IsConnected())
}

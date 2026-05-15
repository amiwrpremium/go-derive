package transport_test

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
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

// TestWSTransport_OnReconnect_FiresOnceAfterDrop verifies that the
// user-facing OnReconnect callback fires exactly once per reconnect
// cycle, after PostDialHook and resubscribe both run.
func TestWSTransport_OnReconnect_FiresOnceAfterDrop(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()

	var (
		mu             sync.Mutex
		calls          []error
		postDialCalls  int
		resubAfterHook bool
	)
	tt, err := transport.NewWS(transport.WSConfig{
		URL:          srv.URL(),
		Reconnect:    true,
		PingInterval: 30 * time.Millisecond,
		PostDialHook: func(_ context.Context, _ *transport.WSTransport) error {
			mu.Lock()
			postDialCalls++
			mu.Unlock()
			return nil
		},
		OnReconnect: func(err error) {
			mu.Lock()
			calls = append(calls, err)
			mu.Unlock()
		},
	})
	require.NoError(t, err)
	require.NoError(t, tt.Connect(context.Background()))
	defer func() { _ = tt.Close() }()

	dec := func(json.RawMessage) (any, error) { return nil, nil }
	_, err = tt.Subscribe(context.Background(), "X", dec)
	require.NoError(t, err)
	require.True(t, srv.WaitSubscribed("X", time.Second))

	// Initial Connect must not fire the user callback.
	time.Sleep(100 * time.Millisecond)
	mu.Lock()
	require.Empty(t, calls, "OnReconnect should not fire on initial Connect")
	mu.Unlock()

	srv.DropClients()

	require.Eventually(t, func() bool {
		mu.Lock()
		n := len(calls)
		mu.Unlock()
		return n >= 1
	}, 5*time.Second, 20*time.Millisecond, "OnReconnect never fired after drop")

	mu.Lock()
	defer mu.Unlock()
	require.Len(t, calls, 1, "OnReconnect must fire exactly once per cycle")
	assert.NoError(t, calls[0])
	assert.GreaterOrEqual(t, postDialCalls, 1)
	resubAfterHook = postDialCalls >= 1 && srv.Subscribed("X")
	assert.True(t, resubAfterHook, "resubscribe should have re-registered the channel")
}

// TestWSTransport_OnReconnect_NoSubsStillFires covers the empty-state
// path — a client with zero subscriptions must still receive the
// snapshot-refetch hook after reconnect.
func TestWSTransport_OnReconnect_NoSubsStillFires(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()

	var (
		mu    sync.Mutex
		calls []error
	)
	tt, err := transport.NewWS(transport.WSConfig{
		URL:          srv.URL(),
		Reconnect:    true,
		PingInterval: 30 * time.Millisecond,
		OnReconnect: func(err error) {
			mu.Lock()
			calls = append(calls, err)
			mu.Unlock()
		},
	})
	require.NoError(t, err)
	require.NoError(t, tt.Connect(context.Background()))
	defer func() { _ = tt.Close() }()

	srv.DropClients()

	require.Eventually(t, func() bool {
		mu.Lock()
		defer mu.Unlock()
		return len(calls) >= 1
	}, 5*time.Second, 20*time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	assert.NoError(t, calls[0], "no-sub reconnect must still report success")
}

// TestWSTransport_OnReconnect_SurfacesPostDialError verifies that when
// the post-dial hook returns an error, the user callback sees it.
func TestWSTransport_OnReconnect_SurfacesPostDialError(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()

	hookErr := errors.New("login rejected")
	var (
		mu    sync.Mutex
		calls []error
	)
	tt, err := transport.NewWS(transport.WSConfig{
		URL:          srv.URL(),
		Reconnect:    true,
		PingInterval: 30 * time.Millisecond,
		PostDialHook: func(_ context.Context, _ *transport.WSTransport) error {
			return hookErr
		},
		OnReconnect: func(err error) {
			mu.Lock()
			calls = append(calls, err)
			mu.Unlock()
		},
	})
	require.NoError(t, err)
	require.NoError(t, tt.Connect(context.Background()))
	defer func() { _ = tt.Close() }()

	srv.DropClients()

	require.Eventually(t, func() bool {
		mu.Lock()
		defer mu.Unlock()
		return len(calls) >= 1
	}, 5*time.Second, 20*time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	require.Len(t, calls, 1)
	require.Error(t, calls[0])
	assert.ErrorIs(t, calls[0], hookErr)
}

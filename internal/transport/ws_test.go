package transport_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/internal/jsonrpc"
	"github.com/amiwrpremium/go-derive/internal/testutil"
	"github.com/amiwrpremium/go-derive/internal/transport"
	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
)

func newConnectedWS(t *testing.T, srv *testutil.MockWSServer) *transport.WSTransport {
	t.Helper()
	cfg := transport.WSConfig{URL: srv.URL(), PingInterval: 50 * time.Millisecond}
	tt, err := transport.NewWS(cfg)
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	require.NoError(t, tt.Connect(ctx))
	return tt
}

func TestWSTransport_RequiresURL(t *testing.T) {
	_, err := transport.NewWS(transport.WSConfig{})
	assert.Error(t, err)
}

func TestWSTransport_CallSuccess(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	srv.HandleResult("public/get_time", 1700000000000)

	tt := newConnectedWS(t, srv)
	defer func() { _ = tt.Close() }()

	var got int64
	require.NoError(t, tt.Call(context.Background(), "public/get_time", nil, &got))
	assert.Equal(t, int64(1700000000000), got)
}

func TestWSTransport_CallAPIError(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	srv.Handle("x", func(json.RawMessage) (any, *jsonrpc.Error) {
		return nil, &jsonrpc.Error{Code: 10002, Message: "rate"}
	})

	tt := newConnectedWS(t, srv)
	defer func() { _ = tt.Close() }()

	err := tt.Call(context.Background(), "x", nil, nil)
	require.Error(t, err)
	var apiErr *derrors.APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, 10002, apiErr.Code)
}

func TestWSTransport_CallContextCancelled(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	// No handler — server returns method-not-found, but we cancel before
	// the response can arrive.
	srv.Handle("slow", func(json.RawMessage) (any, *jsonrpc.Error) {
		time.Sleep(200 * time.Millisecond)
		return "ok", nil
	})
	tt := newConnectedWS(t, srv)
	defer func() { _ = tt.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	err := tt.Call(ctx, "slow", nil, nil)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestWSTransport_CallNotConnected(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	tt, err := transport.NewWS(transport.WSConfig{URL: srv.URL()})
	require.NoError(t, err)
	err = tt.Call(context.Background(), "x", nil, nil)
	assert.ErrorIs(t, err, derrors.ErrNotConnected)
}

func TestWSTransport_DoubleConnect(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	tt := newConnectedWS(t, srv)
	defer func() { _ = tt.Close() }()
	err := tt.Connect(context.Background())
	assert.ErrorIs(t, err, derrors.ErrAlreadyConnected)
}

func TestWSTransport_Subscribe_NotificationDispatch(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()

	tt := newConnectedWS(t, srv)
	defer func() { _ = tt.Close() }()

	dec := func(raw json.RawMessage) (any, error) {
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			return nil, err
		}
		return s, nil
	}
	sub, err := tt.Subscribe(context.Background(), "trades.BTC-PERP", dec)
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()

	require.True(t, srv.WaitSubscribed("trades.BTC-PERP", 1*time.Second))
	srv.Notify("trades.BTC-PERP", "first")
	srv.Notify("trades.BTC-PERP", "second")

	got := []string{}
	deadline := time.After(1 * time.Second)
	for len(got) < 2 {
		select {
		case <-deadline:
			t.Fatalf("got only %v before timeout", got)
		case v, ok := <-sub.Updates():
			require.True(t, ok)
			got = append(got, v.(string))
		}
	}
	assert.Equal(t, []string{"first", "second"}, got)
}

func TestWSTransport_IsConnected(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	tt := newConnectedWS(t, srv)
	assert.True(t, tt.IsConnected())
	require.NoError(t, tt.Close())
	assert.False(t, tt.IsConnected())
}

func TestWSTransport_CloseFinishesPending(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	tt := newConnectedWS(t, srv)
	require.NoError(t, tt.Close())
	// A second Close is a no-op.
	require.NoError(t, tt.Close())
}

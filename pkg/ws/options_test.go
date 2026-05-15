package ws_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/internal/netconf"
	"github.com/amiwrpremium/go-derive/internal/testutil"
	"github.com/amiwrpremium/go-derive/pkg/auth"
	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
	"github.com/amiwrpremium/go-derive/pkg/ws"
)

func TestWS_RequiresNetwork(t *testing.T) {
	_, err := ws.New()
	assert.ErrorIs(t, err, derrors.ErrInvalidConfig)
}

func TestWS_WithMainnet(t *testing.T) {
	c, err := ws.New(ws.WithMainnet())
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
	assert.Equal(t, netconf.NetworkMainnet, c.Network().Network)
}

func TestWS_WithTestnet(t *testing.T) {
	c, err := ws.New(ws.WithTestnet())
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
	assert.Equal(t, netconf.NetworkTestnet, c.Network().Network)
}

func TestWS_WithCustomNetwork(t *testing.T) {
	custom := netconf.Testnet()
	custom.WSURL = "ws://example.invalid/ws"
	c, err := ws.New(ws.WithCustomNetwork(custom))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
	assert.Equal(t, "ws://example.invalid/ws", c.Network().WSURL)
}

func TestWS_AllOptionsCompose(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()

	signer, err := auth.NewLocalSigner(testKey)
	require.NoError(t, err)

	cfg := netconf.Testnet()
	cfg.WSURL = srv.URL()

	c, err := ws.New(
		ws.WithCustomNetwork(cfg),
		ws.WithSigner(signer),
		ws.WithSubaccount(7),
		ws.WithUserAgent("custom/1"),
		ws.WithRateLimit(50, 2),
		ws.WithPingInterval(100*time.Millisecond),
		ws.WithReconnect(false),
		ws.WithSignatureExpiry(120),
	)
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
	assert.NotNil(t, c)
}

func TestWS_DefaultsApplied(t *testing.T) {
	// New with only WithTestnet still produces a valid client — verifies
	// the default tps/burst/ping/reconnect/expiry values don't panic.
	c, err := ws.New(ws.WithTestnet())
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
}

// TestWS_WithOnReconnect_FiresAfterDrop confirms the option plumbs
// through to the transport — when the server severs every client
// connection the callback fires exactly once with a nil error.
func TestWS_WithOnReconnect_FiresAfterDrop(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	cfg := netconf.Testnet()
	cfg.WSURL = srv.URL()

	var (
		mu    sync.Mutex
		calls []error
	)
	c, err := ws.New(
		ws.WithCustomNetwork(cfg),
		ws.WithReconnect(true),
		ws.WithPingInterval(50*time.Millisecond),
		ws.WithOnReconnect(func(err error) {
			mu.Lock()
			calls = append(calls, err)
			mu.Unlock()
		}),
	)
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
	require.NoError(t, c.Connect(context.Background()))

	srv.DropClients()

	require.Eventually(t, func() bool {
		mu.Lock()
		defer mu.Unlock()
		return len(calls) >= 1
	}, 5*time.Second, 20*time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	require.Len(t, calls, 1)
	assert.NoError(t, calls[0])
}

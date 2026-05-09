package ws_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/internal/testutil"
	"github.com/amiwrpremium/go-derive/pkg/ws"
)

func TestWS_RequiresNetwork(t *testing.T) {
	_, err := ws.New()
	assert.ErrorIs(t, err, derive.ErrInvalidConfig)
}

func TestWS_WithMainnet(t *testing.T) {
	c, err := ws.New(ws.WithMainnet())
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
	assert.Equal(t, derive.NetworkMainnet, c.Network().Network)
}

func TestWS_WithTestnet(t *testing.T) {
	c, err := ws.New(ws.WithTestnet())
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
	assert.Equal(t, derive.NetworkTestnet, c.Network().Network)
}

func TestWS_WithCustomNetwork(t *testing.T) {
	custom := derive.Testnet()
	custom.WSURL = "ws://example.invalid/ws"
	c, err := ws.New(ws.WithCustomNetwork(custom))
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

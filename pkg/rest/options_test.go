package rest_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/internal/jsonrpc"
	"github.com/amiwrpremium/go-derive/internal/netconf"
	"github.com/amiwrpremium/go-derive/internal/testutil"
	"github.com/amiwrpremium/go-derive/pkg/auth"
	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
	"github.com/amiwrpremium/go-derive/pkg/rest"
)

const testKey = "0x4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318"

func TestRest_WithMainnet(t *testing.T) {
	c, err := rest.New(rest.WithMainnet())
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
	assert.Equal(t, netconf.NetworkMainnet, c.Network().Network)
}

func TestRest_WithTestnet(t *testing.T) {
	c, err := rest.New(rest.WithTestnet())
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
	assert.Equal(t, netconf.NetworkTestnet, c.Network().Network)
}

func TestRest_RequiresNetwork(t *testing.T) {
	_, err := rest.New()
	assert.ErrorIs(t, err, derrors.ErrInvalidConfig)
}

func TestRest_AllOptionsApplied(t *testing.T) {
	srv := testutil.NewMockServer()
	defer srv.Close()
	srv.Handle("public/get_time", func(_ testutil.MockRequest) (any, *jsonrpc.Error) {
		return 1700000000000, nil
	})

	cfg := netconf.Testnet()
	cfg.HTTPURL = srv.URL()

	signer, err := auth.NewLocalSigner(testKey)
	require.NoError(t, err)

	customClient := &http.Client{Timeout: 10 * time.Second}

	c, err := rest.New(
		rest.WithCustomNetwork(cfg),
		rest.WithSigner(signer),
		rest.WithSubaccount(99),
		rest.WithHTTPClient(customClient),
		rest.WithUserAgent("custom-agent/1.0"),
		rest.WithRateLimit(50, 2),
		rest.WithSignatureExpiry(60),
	)
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	// One real call against the mock, asserting the User-Agent option flowed through.
	_, err = c.GetTime(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, srv.Requests())
	assert.Equal(t, "custom-agent/1.0", srv.Requests()[0].Headers.Get("User-Agent"))
}

func TestRest_WithHTTPTimeout_FiresOnSlowResponse(t *testing.T) {
	srv := testutil.NewMockServer()
	defer srv.Close()
	srv.Handle("public/get_time", func(_ testutil.MockRequest) (any, *jsonrpc.Error) {
		// Outlast the configured timeout — caller should see a
		// deadline-exceeded error before this returns.
		time.Sleep(200 * time.Millisecond)
		return 1700000000000, nil
	})
	cfg := netconf.Testnet()
	cfg.HTTPURL = srv.URL()

	c, err := rest.New(
		rest.WithCustomNetwork(cfg),
		rest.WithHTTPTimeout(50*time.Millisecond),
	)
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	_, err = c.GetTime(context.Background())
	require.Error(t, err)
	// http.Client surfaces a *url.Error wrapping a timeout. We don't
	// care about the exact wrapper — just that it fired.
	assert.Contains(t, err.Error(), "deadline exceeded",
		"expected timeout error, got: %v", err)
}

func TestRest_WithHTTPClient_OverridesWithHTTPTimeout(t *testing.T) {
	// Build a *http.Client whose timeout is longer than the slow
	// handler so a "fast enough" response wins. If WithHTTPTimeout
	// shadowed WithHTTPClient, the 50ms cap would fire and the call
	// would fail — instead it succeeds because WithHTTPClient wins.
	srv := testutil.NewMockServer()
	defer srv.Close()
	srv.Handle("public/get_time", func(_ testutil.MockRequest) (any, *jsonrpc.Error) {
		time.Sleep(100 * time.Millisecond)
		return 1700000000000, nil
	})
	cfg := netconf.Testnet()
	cfg.HTTPURL = srv.URL()

	c, err := rest.New(
		rest.WithCustomNetwork(cfg),
		rest.WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
		rest.WithHTTPTimeout(50*time.Millisecond), // ignored
	)
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	got, err := c.GetTime(context.Background())
	require.NoError(t, err, "WithHTTPClient must take precedence over WithHTTPTimeout")
	assert.Equal(t, int64(1700000000000), got)
}

func TestRest_SignerAttachesAuthHeaders(t *testing.T) {
	srv := testutil.NewMockServer()
	defer srv.Close()
	srv.Handle("public/get_time", func(_ testutil.MockRequest) (any, *jsonrpc.Error) {
		return 1700000000000, nil
	})

	cfg := netconf.Testnet()
	cfg.HTTPURL = srv.URL()
	signer, err := auth.NewLocalSigner(testKey)
	require.NoError(t, err)

	c, err := rest.New(rest.WithCustomNetwork(cfg), rest.WithSigner(signer))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
	_, err = c.GetTime(context.Background())
	require.NoError(t, err)

	require.NotEmpty(t, srv.Requests())
	headers := srv.Requests()[0].Headers
	assert.Equal(t, signer.Owner().Hex(), headers.Get("X-LyraWallet"))
	assert.NotEmpty(t, headers.Get("X-LyraTimestamp"))
	assert.NotEmpty(t, headers.Get("X-LyraSignature"))
}

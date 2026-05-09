package derive_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/internal/jsonrpc"
	"github.com/amiwrpremium/go-derive/internal/testutil"
)

// withRESTMock configures a RestClient pointed at a mock server.
// Network-aware helpers like WithTestnet would otherwise resolve to the
// real Derive URL.
func withRESTMock(t *testing.T, srv *testutil.MockServer) *derive.RestClient {
	t.Helper()
	cfg := derive.Testnet()
	cfg.HTTPURL = srv.URL()
	c, err := derive.NewRestClient(derive.WithCustomNetwork(cfg))
	require.NoError(t, err)
	return c
}

func TestGetInstruments_DecodesPayload(t *testing.T) {
	srv := testutil.NewMockServer()
	defer srv.Close()

	srv.Handle("public/get_instruments", func(_ testutil.MockRequest) (any, *jsonrpc.Error) {
		return []map[string]any{
			{
				"instrument_name": "BTC-PERP",
				"base_currency":   "BTC",
				"quote_currency":  "USDC",
				"instrument_type": "perp",
				"is_active":       true,
				"tick_size":       "0.5",
				"minimum_amount":  "0.001",
				"maximum_amount":  "1000",
				"amount_step":     "0.001",
				"mark_price":      "65000.5",
				"index_price":     "64999",
			},
		}, nil
	})

	c := withRESTMock(t, srv)
	defer c.Close()

	insts, err := c.GetInstruments(context.Background(), "BTC", "perp")
	require.NoError(t, err)
	require.Len(t, insts, 1)
	assert.Equal(t, "BTC-PERP", insts[0].Name)
	assert.Equal(t, "65000.5", insts[0].MarkPrice.String())
}

func TestAPIError_MapsToSentinel(t *testing.T) {
	srv := testutil.NewMockServer()
	defer srv.Close()

	srv.Handle("public/get_time", func(_ testutil.MockRequest) (any, *jsonrpc.Error) {
		return nil, &jsonrpc.Error{Code: derive.CodeRateLimitExceeded, Message: "rate limited"}
	})

	c := withRESTMock(t, srv)
	defer c.Close()

	_, err := c.GetTime(context.Background())
	require.Error(t, err)

	assert.True(t, derive.Is(err, derive.ErrRateLimited),
		"expected rate-limit code %d to map to ErrRateLimited; got %v",
		derive.CodeRateLimitExceeded, err)
}

func TestPrivateMethod_RequiresSubaccount(t *testing.T) {
	srv := testutil.NewMockServer()
	defer srv.Close()

	c := withRESTMock(t, srv)
	defer c.Close()

	_, err := c.GetPositions(context.Background())
	assert.True(t, derive.Is(err, derive.ErrSubaccountRequired), "got %v", err)
}

func TestRest_WithMainnet(t *testing.T) {
	c, err := derive.NewRestClient(derive.WithMainnet())
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
	assert.Equal(t, derive.NetworkMainnet, c.Network().Network)
}

func TestRest_WithTestnet(t *testing.T) {
	c, err := derive.NewRestClient(derive.WithTestnet())
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
	assert.Equal(t, derive.NetworkTestnet, c.Network().Network)
}

func TestRest_RequiresNetwork(t *testing.T) {
	_, err := derive.NewRestClient()
	assert.ErrorIs(t, err, derive.ErrInvalidConfig)
}

func TestRest_AllOptionsApplied(t *testing.T) {
	srv := testutil.NewMockServer()
	defer srv.Close()
	srv.Handle("public/get_time", func(_ testutil.MockRequest) (any, *jsonrpc.Error) {
		return 1700000000000, nil
	})

	cfg := derive.Testnet()
	cfg.HTTPURL = srv.URL()

	signer, err := derive.NewLocalSigner(testKey)
	require.NoError(t, err)

	customClient := &http.Client{Timeout: 10 * time.Second}

	c, err := derive.NewRestClient(
		derive.WithCustomNetwork(cfg),
		derive.WithSigner(signer),
		derive.WithSubaccount(99),
		derive.WithHTTPClient(customClient),
		derive.WithUserAgent("custom-agent/1.0"),
		derive.WithRateLimit(50, 2),
		derive.WithSignatureExpiry(60),
	)
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	_, err = c.GetTime(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, srv.Requests())
	assert.Equal(t, "custom-agent/1.0", srv.Requests()[0].Headers.Get("User-Agent"))
}

func TestRest_SignerAttachesAuthHeaders(t *testing.T) {
	srv := testutil.NewMockServer()
	defer srv.Close()
	srv.Handle("public/get_time", func(_ testutil.MockRequest) (any, *jsonrpc.Error) {
		return 1700000000000, nil
	})

	cfg := derive.Testnet()
	cfg.HTTPURL = srv.URL()
	signer, err := derive.NewLocalSigner(testKey)
	require.NoError(t, err)

	c, err := derive.NewRestClient(derive.WithCustomNetwork(cfg), derive.WithSigner(signer))
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

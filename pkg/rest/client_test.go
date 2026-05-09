package rest_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/internal/jsonrpc"
	"github.com/amiwrpremium/go-derive/internal/testutil"
	"github.com/amiwrpremium/go-derive/pkg/rest"
)

// withMock configures a Client pointed at a mock server. Network-aware
// helpers like WithTestnet would otherwise resolve to the real Derive URL.
func withMock(t *testing.T, srv *testutil.MockServer) *rest.Client {
	t.Helper()
	cfg := derive.Testnet()
	cfg.HTTPURL = srv.URL()
	c, err := rest.New(rest.WithCustomNetwork(cfg))
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

	c := withMock(t, srv)
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

	c := withMock(t, srv)
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

	c := withMock(t, srv) // no signer, no subaccount
	defer c.Close()

	_, err := c.GetPositions(context.Background())
	assert.True(t, derive.Is(err, derive.ErrSubaccountRequired), "got %v", err)
}

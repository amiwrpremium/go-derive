//go:build integration

package integration_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestPublic_GetTime(t *testing.T) {
	env := loadEnv(t)
	c := newRESTClient(t, env)
	ctx, cancel := withTimeout(t)
	defer cancel()

	got, err := c.GetTime(ctx)
	require.NoError(t, err)

	delta := time.Now().UnixMilli() - got
	if delta < 0 {
		delta = -delta
	}
	assert.LessOrEqual(t, delta, int64(60_000), "server time within ±60s of local")
}

func TestPublic_GetCurrencies(t *testing.T) {
	env := loadEnv(t)
	c := newRESTClient(t, env)
	ctx, cancel := withTimeout(t)
	defer cancel()

	got, err := c.GetCurrencies(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, got)
}

func TestPublic_GetInstruments_BTCPerps(t *testing.T) {
	env := loadEnv(t)
	c := newRESTClient(t, env)
	ctx, cancel := withTimeout(t)
	defer cancel()

	insts, err := c.GetInstruments(ctx, "BTC", enums.InstrumentTypePerp)
	require.NoError(t, err)
	require.NotEmpty(t, insts, "Derive testnet should have at least one BTC perp")
	for _, in := range insts {
		assert.Equal(t, enums.InstrumentTypePerp, in.Type)
	}
}

func TestPublic_GetTicker(t *testing.T) {
	env := loadEnv(t)
	c := newRESTClient(t, env)
	ctx, cancel := withTimeout(t)
	defer cancel()

	tk, err := c.GetTicker(ctx, env.instrument)
	require.NoError(t, err, "ticker for %s", env.instrument)

	assert.Equal(t, env.instrument, tk.InstrumentName)
	if !tk.BestBidPrice.IsZero() && !tk.BestAskPrice.IsZero() {
		bid := tk.BestBidPrice.Inner()
		ask := tk.BestAskPrice.Inner()
		assert.True(t, bid.Cmp(ask) <= 0, "bid <= ask: %s vs %s", bid, ask)
	}
	assert.Equal(t, 1, tk.MarkPrice.Sign(), "mark price should be positive")
}

// `public/get_orderbook` was removed from Derive's REST surface as of
// 2024+; full L2 orderbook is now WebSocket-subscription-only via the
// `orderbook.<inst>.<group>.<depth>` channel. Top-of-book + 5 % depth
// lives in `public/get_ticker`. The WS subscription is exercised by
// `TestWS_OrderBookSubscribe` in `ws_public_integration_test.go`.

func TestPublic_GetTradeHistory(t *testing.T) {
	env := loadEnv(t)
	c := newRESTClient(t, env)
	ctx, cancel := withTimeout(t)
	defer cancel()

	trades, page, err := c.GetPublicTradeHistory(ctx, env.instrument, types.PageRequest{PageSize: 10})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, page.PageSize, 0)
	for _, tr := range trades {
		assert.Equal(t, env.instrument, tr.InstrumentName)
	}
}

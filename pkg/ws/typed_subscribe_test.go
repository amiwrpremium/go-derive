package ws_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/internal/testutil"
	"github.com/amiwrpremium/go-derive/pkg/enums"
)

func TestClient_SubscribeTrades_Wires(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := c.SubscribeTrades(context.Background(), "BTC-PERP")
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()
	assert.True(t, srv.WaitSubscribed("trades.BTC-PERP", time.Second))
}

func TestClient_SubscribeOrderBook_DefaultsAndOverrides(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub1, err := c.SubscribeOrderBook(context.Background(), "BTC-PERP", "", 0)
	require.NoError(t, err)
	defer func() { _ = sub1.Close() }()
	assert.True(t, srv.WaitSubscribed("orderbook.BTC-PERP.1.10", time.Second))

	sub2, err := c.SubscribeOrderBook(context.Background(), "ETH-PERP", "10", 25)
	require.NoError(t, err)
	defer func() { _ = sub2.Close() }()
	assert.True(t, srv.WaitSubscribed("orderbook.ETH-PERP.10.25", time.Second))
}

func TestClient_SubscribeTicker_DefaultInterval(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := c.SubscribeTicker(context.Background(), "BTC-PERP", "")
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()
	assert.True(t, srv.WaitSubscribed("ticker.BTC-PERP.1000", time.Second))
}

func TestClient_SubscribeTickerSlim_ExplicitInterval(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := c.SubscribeTickerSlim(context.Background(), "BTC-PERP", "100")
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()
	assert.True(t, srv.WaitSubscribed("ticker_slim.BTC-PERP.100", time.Second))
}

func TestClient_SubscribeSpotFeed(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := c.SubscribeSpotFeed(context.Background(), "BTC")
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()
	assert.True(t, srv.WaitSubscribed("spot_feed.BTC", time.Second))
}

func TestClient_SubscribeMarginWatch(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := c.SubscribeMarginWatch(context.Background())
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()
	assert.True(t, srv.WaitSubscribed("margin.watch", time.Second))
}

func TestClient_SubscribeAuctionsWatch(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := c.SubscribeAuctionsWatch(context.Background())
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()
	assert.True(t, srv.WaitSubscribed("auctions.watch", time.Second))
}

func TestClient_SubscribeTradesByType(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := c.SubscribeTradesByType(context.Background(), enums.InstrumentTypePerp, "BTC")
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()
	assert.True(t, srv.WaitSubscribed("trades.perp.BTC", time.Second))
}

func TestClient_SubscribeTradesByTypeWithStatus(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := c.SubscribeTradesByTypeWithStatus(context.Background(), enums.InstrumentTypePerp, "BTC", enums.TxStatusSettled)
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()
	assert.True(t, srv.WaitSubscribed("trades.perp.BTC.settled", time.Second))
}

func TestClient_SubscribeOrders(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := c.SubscribeOrders(context.Background(), 7)
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()
	assert.True(t, srv.WaitSubscribed("7.orders", time.Second))
}

func TestClient_SubscribeBalances(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := c.SubscribeBalances(context.Background(), 7)
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()
	assert.True(t, srv.WaitSubscribed("7.balances", time.Second))
}

func TestClient_SubscribeBestQuotes(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := c.SubscribeBestQuotes(context.Background(), 9)
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()
	assert.True(t, srv.WaitSubscribed("9.best.quotes", time.Second))
}

func TestClient_SubscribeQuotes(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := c.SubscribeQuotes(context.Background(), 7)
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()
	assert.True(t, srv.WaitSubscribed("7.quotes", time.Second))
}

func TestClient_SubscribeRFQs(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := c.SubscribeRFQs(context.Background(), "0xabc")
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()
	assert.True(t, srv.WaitSubscribed("0xabc.rfqs", time.Second))
}

func TestClient_SubscribeSubaccountTrades(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := c.SubscribeSubaccountTrades(context.Background(), 7)
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()
	assert.True(t, srv.WaitSubscribed("7.trades", time.Second))
}

func TestClient_SubscribeSubaccountTradesByStatus(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := c.SubscribeSubaccountTradesByStatus(context.Background(), 7, enums.TxStatusSettled)
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()
	assert.True(t, srv.WaitSubscribed("7.trades.settled", time.Second))
}

// End-to-end decode through one typed method to confirm the type
// parameter is wired through correctly.
func TestClient_SubscribeOrders_DeliversTyped(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := c.SubscribeOrders(context.Background(), 7)
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()
	require.True(t, srv.WaitSubscribed("7.orders", time.Second))

	srv.Notify("7.orders", []map[string]any{{
		"order_id": "O1", "subaccount_id": 7, "instrument_name": "BTC-PERP",
		"direction": "buy", "order_type": "limit", "time_in_force": "gtc",
		"order_status": "open", "amount": "0.1", "filled_amount": "0",
		"limit_price": "65000", "max_fee": "10", "nonce": 1,
		"signer":             "0x0000000000000000000000000000000000000000",
		"creation_timestamp": 1700000000000, "last_update_timestamp": 1700000000000,
	}})

	select {
	case batch, ok := <-sub.Updates():
		require.True(t, ok)
		require.Len(t, batch, 1)
		assert.Equal(t, "O1", batch[0].OrderID)
	case <-time.After(2 * time.Second):
		t.Fatal("typed subscribe never delivered")
	}
}

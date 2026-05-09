//go:build integration

package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive"
)

// requireLiveOrders skips the test unless DERIVE_RUN_LIVE_ORDERS=1 AND
// private creds AND a base asset are configured.
//
// Live order placement creates real (testnet) state; we double-gate it.
func requireLiveOrders(t *testing.T, env integrationEnv) {
	t.Helper()
	if !env.liveOrders {
		t.Skip("DERIVE_RUN_LIVE_ORDERS=1 not set; skipping live-order test")
	}
	if env.baseAsset == (common.Address{}) {
		t.Skip("DERIVE_BASE_ASSET not set; cannot sign trade module data")
	}
	if !env.hasPrivateCreds() {
		t.Skip("private creds not configured")
	}
}

// farFromMarkBuy returns a limit price 5% below mark — far enough away that
// the order should not fill in the test window.
func farFromMarkBuy(mark derive.Decimal) derive.Decimal {
	priced := mark.Inner().Mul(decimal.RequireFromString("0.95"))
	d, _ := derive.NewDecimal(priced.String())
	return d
}

func TestPrivate_PlaceAndCancelOrder_REST(t *testing.T) {
	env := loadEnv(t)
	requireLiveOrders(t, env)
	c := newAuthRESTClient(t, env)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	tk, err := c.GetTicker(ctx, env.instrument)
	require.NoError(t, err)

	in := derive.PlaceOrderInput{
		InstrumentName: env.instrument,
		Asset:          env.baseAsset,
		SubID:          0,
		Direction:      derive.DirectionBuy,
		OrderType:      derive.OrderTypeLimit,
		TimeInForce:    derive.TimeInForceGTC,
		Amount:         derive.MustDecimal("0.001"),
		LimitPrice:     farFromMarkBuy(tk.MarkPrice),
		MaxFee:         derive.MustDecimal("10"),
		Label:          "go-derive-integration",
	}
	order, err := c.PlaceOrder(ctx, in)
	require.NoError(t, err, "PlaceOrder")
	require.NotEmpty(t, order.OrderID)
	assert.Equal(t, derive.DirectionBuy, order.Direction)

	require.NoError(t, c.CancelOrder(ctx, env.instrument, order.OrderID))
}

func TestPrivate_PlaceAndCancelOrder_WS(t *testing.T) {
	env := loadEnv(t)
	requireLiveOrders(t, env)
	c := newAuthWSClient(t, env)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	tk, err := c.GetTicker(ctx, env.instrument)
	require.NoError(t, err)

	in := derive.PlaceOrderInput{
		InstrumentName: env.instrument,
		Asset:          env.baseAsset,
		Direction:      derive.DirectionBuy,
		OrderType:      derive.OrderTypeLimit,
		TimeInForce:    derive.TimeInForceGTC,
		Amount:         derive.MustDecimal("0.001"),
		LimitPrice:     farFromMarkBuy(tk.MarkPrice),
		MaxFee:         derive.MustDecimal("10"),
		Label:          "go-derive-integration-ws",
	}
	order, err := c.PlaceOrder(ctx, in)
	require.NoError(t, err)
	require.NotEmpty(t, order.OrderID)

	require.NoError(t, c.CancelOrder(ctx, env.instrument, order.OrderID))
}

func TestPrivate_OrderEventsArrive(t *testing.T) {
	env := loadEnv(t)
	requireLiveOrders(t, env)
	c := newAuthWSClient(t, env)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	sub, err := derive.Subscribe[[]derive.Order](ctx, c, derive.PrivateOrders{SubaccountID: env.subaccount})
	require.NoError(t, err)
	defer sub.Close()

	tk, err := c.GetTicker(ctx, env.instrument)
	require.NoError(t, err)

	in := derive.PlaceOrderInput{
		InstrumentName: env.instrument,
		Asset:          env.baseAsset,
		Direction:      derive.DirectionBuy,
		OrderType:      derive.OrderTypeLimit,
		TimeInForce:    derive.TimeInForceGTC,
		Amount:         derive.MustDecimal("0.001"),
		LimitPrice:     farFromMarkBuy(tk.MarkPrice),
		MaxFee:         derive.MustDecimal("10"),
		Label:          "go-derive-integration-events",
	}
	order, err := c.PlaceOrder(ctx, in)
	require.NoError(t, err)
	defer func() { _ = c.CancelOrder(ctx, env.instrument, order.OrderID) }()

	// Wait for the order event to flow through the subscription.
	deadline := time.After(20 * time.Second)
	for {
		select {
		case batch, ok := <-sub.Updates():
			if !ok {
				t.Fatalf("subscription closed: %v", sub.Err())
			}
			for _, o := range batch {
				if o.OrderID == order.OrderID {
					return // success
				}
			}
		case <-deadline:
			t.Fatal("placed order did not appear on the orders channel within 20s")
		}
	}
}

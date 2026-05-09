// Places one limit order over WebSocket. Requires DERIVE_BASE_ASSET and
// DERIVE_RUN_LIVE_ORDERS=1.
package main

import "github.com/amiwrpremium/go-derive"

import (
	"log"
	"os"

	"github.com/shopspring/decimal"

	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/internal/methods"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	if os.Getenv("DERIVE_RUN_LIVE_ORDERS") != "1" {
		log.Fatal("set DERIVE_RUN_LIVE_ORDERS=1")
	}
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	tk, err := c.GetTicker(ctx, example.Instrument())
	example.Fatal(err)
	limit := tk.MarkPrice.Inner().Mul(decimal.RequireFromString("0.95"))
	price, _ := types.NewDecimal(limit.String())

	o, err := c.PlaceOrder(ctx, methods.PlaceOrderInput{
		InstrumentName: example.Instrument(),
		Asset:          example.BaseAsset(),
		Direction:      derive.DirectionBuy,
		OrderType:      derive.OrderTypeLimit,
		TimeInForce:    derive.TimeInForceGTC,
		Amount:         types.MustDecimal("0.001"),
		LimitPrice:     price,
		MaxFee:         types.MustDecimal("10"),
	})
	example.Fatal(err)
	example.Print("placed", o.OrderID)
}

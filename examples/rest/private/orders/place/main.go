// Places one limit buy 5% below mark (won't fill on a normal book).
//
// Requires: DERIVE_BASE_ASSET, DERIVE_RUN_LIVE_ORDERS=1.
package main

import "github.com/amiwrpremium/go-derive"

import (
	"log"
	"os"

	"github.com/shopspring/decimal"

	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/internal/methods"
)

func main() {
	if os.Getenv("DERIVE_RUN_LIVE_ORDERS") != "1" {
		log.Fatal("set DERIVE_RUN_LIVE_ORDERS=1 to actually place an order")
	}
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	tk, err := c.GetTicker(ctx, example.Instrument())
	example.Fatal(err)
	limit := tk.MarkPrice.Inner().Mul(decimal.RequireFromString("0.95"))
	price, _ := derive.NewDecimal(limit.String())

	o, err := c.PlaceOrder(ctx, methods.PlaceOrderInput{
		InstrumentName: example.Instrument(),
		Asset:          example.BaseAsset(),
		Direction:      derive.DirectionBuy,
		OrderType:      derive.OrderTypeLimit,
		TimeInForce:    derive.TimeInForceGTC,
		Amount:         derive.MustDecimal("0.001"),
		LimitPrice:     price,
		MaxFee:         derive.MustDecimal("10"),
	})
	example.Fatal(err)
	example.Print("placed", o.OrderID)
	example.Print("status", o.OrderStatus)
}

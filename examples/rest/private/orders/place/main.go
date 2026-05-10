// Places one limit buy 5% below mark (won't fill on a normal book).
//
// Requires: DERIVE_BASE_ASSET, DERIVE_RUN_LIVE_ORDERS=1.
package main

import (
	"log"
	"os"

	"github.com/shopspring/decimal"

	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
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
	price, _ := types.NewDecimal(limit.String())

	o, err := c.PlaceOrder(ctx, types.PlaceOrderInput{
		InstrumentName: example.Instrument(),
		Asset:          types.Address(example.BaseAsset()),
		Direction:      enums.DirectionBuy,
		OrderType:      enums.OrderTypeLimit,
		TimeInForce:    enums.TimeInForceGTC,
		Amount:         types.MustDecimal("0.001"),
		LimitPrice:     price,
		MaxFee:         types.MustDecimal("10"),
	})
	example.Fatal(err)
	example.Print("placed", o.OrderID)
	example.Print("status", o.OrderStatus)
}

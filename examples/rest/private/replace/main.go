// Replaces (cancel + place) one outstanding order in a single round trip.
//
// Replace is the maker-friendly way to re-price without racing the engine.
// Set DERIVE_ORDER_ID to the order you want to cancel; the replacement
// uses the configured BaseAsset / Instrument.
//
// Double-gate the live submission behind DERIVE_RUN_LIVE_ORDERS=1.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	id := os.Getenv("DERIVE_ORDER_ID")
	if id == "" {
		log.Fatal("DERIVE_ORDER_ID required")
	}
	if os.Getenv("DERIVE_RUN_LIVE_ORDERS") != "1" {
		log.Fatal("set DERIVE_RUN_LIVE_ORDERS=1 to actually submit")
	}
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	res, err := c.Replace(ctx, types.ReplaceOrderInput{
		PlaceOrderInput: types.PlaceOrderInput{
			InstrumentName: example.Instrument(),
			Asset:          types.Address(example.BaseAsset()),
			Direction:      enums.DirectionBuy,
			OrderType:      enums.OrderTypeLimit,
			TimeInForce:    enums.TimeInForceGTC,
			Amount:         types.MustDecimal("0.01"),
			LimitPrice:     types.MustDecimal("50000"),
			MaxFee:         types.MustDecimal("10"),
		},
		OrderIDToCancel: id,
	})
	example.Fatal(err)
	example.Print("cancelled_order", res.CancelledOrder.OrderID)
	if res.Order != nil {
		example.Print("new_order", res.Order.OrderID)
	}
	if res.CreateOrderError != nil {
		example.Print("create_order_error code", res.CreateOrderError.Code)
		example.Print("create_order_error message", res.CreateOrderError.Message)
	}
	example.Print("trade count", len(res.Trades))
}

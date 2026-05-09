// Replaces (cancel + place) one outstanding order in a single round trip.
//
// Replace is the maker-friendly way to re-price without racing the engine.
// Pass `order_id_to_cancel` plus the same fields PlaceOrder would take.
//
// This example is illustrative — set DERIVE_RUN_LIVE_ORDERS=1 only when
// you actually want it to run; the SDK doesn't gate Replace itself.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	res, err := c.Replace(ctx, map[string]any{
		"order_id_to_cancel": "<order-id>",
		"instrument_name":    example.Instrument(),
		"direction":          "buy",
		"order_type":         "limit",
		"time_in_force":      "gtc",
		"amount":             "0.01",
		"limit_price":        "50000",
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

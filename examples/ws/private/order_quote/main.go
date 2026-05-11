// Pre-flights a hypothetical order for the configured subaccount
// over WebSocket. Private counterpart to public/order_quote.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	res, err := c.OrderQuote(ctx, map[string]any{
		"instrument_name": example.Instrument(),
		"direction":       "buy",
		"order_type":      "limit",
		"time_in_force":   "gtc",
		"amount":          "0.01",
		"limit_price":     "50000",
		"max_fee":         "10",
	})
	example.Fatal(err)
	example.Print("is_valid", res.IsValid)
	example.Print("estimated_fee", res.EstimatedFee.String())
}

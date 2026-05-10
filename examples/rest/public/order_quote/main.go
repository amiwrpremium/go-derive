// Pre-flights a hypothetical order through the matching engine
// without auth. Public counterpart to private/order_quote.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	res, err := c.OrderQuotePublic(ctx, map[string]any{
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
	example.Print("estimated_fill_price", res.EstimatedFillPrice.String())
	example.Print("estimated_order_status", string(res.EstimatedOrderStatus))
	example.Print("max_amount", res.MaxAmount.String())
}

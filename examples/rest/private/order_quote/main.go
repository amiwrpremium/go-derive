// Pre-flights a hypothetical order through the matching engine for
// the configured subaccount. Private counterpart to
// public/order_quote — same response shape, but the engine accounts
// for the subaccount's actual margin / collateral.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

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
	example.Print("estimated_fill_price", res.EstimatedFillPrice.String())
	example.Print("max_amount", res.MaxAmount.String())
}

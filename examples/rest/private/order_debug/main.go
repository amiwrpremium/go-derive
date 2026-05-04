// Previews an order without submitting it. Returns the engine's view of
// fees and margin impact, useful for sanity-checking signed payloads in
// CI or pre-flight.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	raw, err := c.OrderDebug(ctx, map[string]any{
		"instrument_name": example.Instrument(),
		"direction":       "buy",
		"order_type":      "limit",
		"time_in_force":   "gtc",
		"amount":          "0.01",
		"limit_price":     "50000",
	})
	example.Fatal(err)
	example.Print("order_debug bytes", len(raw))
}

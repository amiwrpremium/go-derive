// Lists historical oracle spot-feed values for one currency.
// Required env: DERIVE_CURRENCY.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	currency := os.Getenv("DERIVE_CURRENCY")
	if currency == "" {
		log.Fatal("DERIVE_CURRENCY required")
	}
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	cur, items, err := c.GetSpotFeedHistory(ctx, map[string]any{
		"currency": currency,
	})
	example.Fatal(err)
	example.Print("currency", cur)
	example.Print("items", len(items))
}

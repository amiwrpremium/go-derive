// Fetches the per-expiry settlement prices for one currency's
// option market. Pre-settlement entries return Price as zero.
package main

import (
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	currency := os.Getenv("DERIVE_CURRENCY")
	if currency == "" {
		currency = "BTC"
	}
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	prices, err := c.GetOptionSettlementPrices(ctx, currency)
	example.Fatal(err)
	example.Print("expiry count", len(prices))
	for i, p := range prices {
		if i >= 5 {
			break
		}
		example.Print("expiry", p.ExpiryDate)
		example.Print("  utc_expiry_sec", p.UTCExpirySec)
		example.Print("  price", p.Price.String())
	}
}

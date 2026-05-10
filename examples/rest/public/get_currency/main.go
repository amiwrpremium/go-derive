// Fetches per-asset margin parameters and protocol-asset addresses
// for one underlying currency.
package main

import (
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	currency := os.Getenv("DERIVE_CURRENCY")
	if currency == "" {
		currency = "ETH"
	}
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	cur, err := c.GetCurrency(ctx, currency)
	example.Fatal(err)
	example.Print("currency", cur.Currency)
	example.Print("spot_price", cur.SpotPrice.String())
	example.Print("market_type", cur.MarketType)
	example.Print("manager count", len(cur.Managers))
	example.Print("instrument_types", cur.InstrumentTypes)
}

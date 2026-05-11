// Executes (accepts) a quote received against the configured
// subaccount's RFQ. Requires DERIVE_QUOTE_ID and
// DERIVE_RUN_LIVE_ORDERS=1 since execution fills the trade
// immediately.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	quoteID := os.Getenv("DERIVE_QUOTE_ID")
	if quoteID == "" {
		log.Fatal("DERIVE_QUOTE_ID required")
	}
	if os.Getenv("DERIVE_RUN_LIVE_ORDERS") != "1" {
		log.Fatal("set DERIVE_RUN_LIVE_ORDERS=1 to actually execute")
	}
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	res, err := c.ExecuteQuote(ctx, map[string]any{"quote_id": quoteID})
	example.Fatal(err)
	example.Print("rfq filled pct", res.RFQFilledPct.String())
}

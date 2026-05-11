// Polls quotes received against the configured subaccount's RFQs
// over WebSocket. Requires DERIVE_RFQ_ID.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	rfqID := os.Getenv("DERIVE_RFQ_ID")
	if rfqID == "" {
		log.Fatal("DERIVE_RFQ_ID required")
	}
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	quotes, _, err := c.PollQuotes(ctx, map[string]any{"rfq_id": rfqID})
	example.Fatal(err)
	example.Print("quotes", len(quotes))
}

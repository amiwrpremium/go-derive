// Polls quotes received against the configured subaccount's RFQs —
// the maker-side view of what other makers are quoting on an
// active RFQ. Requires DERIVE_RFQ_ID to scope the poll.
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
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	quotes, _, err := c.PollQuotes(ctx, map[string]any{"rfq_id": rfqID})
	example.Fatal(err)
	example.Print("quotes", len(quotes))
}

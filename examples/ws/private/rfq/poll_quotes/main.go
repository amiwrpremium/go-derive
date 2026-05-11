// Polls quotes received against the configured subaccount's RFQs
// over WebSocket. Requires DERIVE_RFQ_ID.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/types"
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

	quotes, _, err := c.PollQuotes(ctx, types.PollQuotesQuery{RFQID: rfqID}, types.PageRequest{})
	example.Fatal(err)
	example.Print("quotes", len(quotes))
}

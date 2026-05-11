// Cancels one outstanding quote by id over WebSocket.
// Requires DERIVE_QUOTE_ID.
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
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	q, err := c.CancelQuote(ctx, map[string]any{"quote_id": quoteID})
	example.Fatal(err)
	example.Print("quote id", q.QuoteID)
}

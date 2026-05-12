// Cancels one outstanding quote by id. Requires DERIVE_QUOTE_ID.
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
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	q, err := c.CancelQuote(ctx, quoteID)
	example.Fatal(err)
	example.Print("quote id", q.QuoteID)
	example.Print("status", q.Status)
}

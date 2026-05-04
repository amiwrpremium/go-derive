// Cancels one RFQ over WebSocket. Set DERIVE_RFQ_ID.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	id := os.Getenv("DERIVE_RFQ_ID")
	if id == "" {
		log.Fatal("DERIVE_RFQ_ID required")
	}
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	example.Fatal(c.CancelRFQ(ctx, id))
	example.Print("cancelled", id)
}

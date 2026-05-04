// Cancels one outstanding RFQ. Set DERIVE_RFQ_ID.
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
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	example.Fatal(c.CancelRFQ(ctx, id))
	example.Print("cancelled", id)
}

// Cancels one order by id over WebSocket. Set DERIVE_ORDER_ID.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	id := os.Getenv("DERIVE_ORDER_ID")
	if id == "" {
		log.Fatal("DERIVE_ORDER_ID required")
	}
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	example.Fatal(c.CancelOrder(ctx, example.Instrument(), id))
	example.Print("cancelled", id)
}

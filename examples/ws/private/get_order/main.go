// Fetches one order by id over WebSocket. Set DERIVE_ORDER_ID.
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

	o, err := c.GetOrder(ctx, id)
	example.Fatal(err)
	example.Print("status", o.OrderStatus)
}

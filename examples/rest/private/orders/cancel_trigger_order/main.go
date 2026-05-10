// Cancels one untriggered trigger order by id. Set DERIVE_ORDER_ID.
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
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	o, err := c.CancelTriggerOrder(ctx, id)
	example.Fatal(err)
	example.Print("cancelled", o.OrderID)
	example.Print("status", string(o.OrderStatus))
}

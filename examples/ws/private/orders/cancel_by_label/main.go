// Cancels orders by label over WebSocket. Set DERIVE_LABEL.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	label := os.Getenv("DERIVE_LABEL")
	if label == "" {
		log.Fatal("DERIVE_LABEL required")
	}
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	n, err := c.CancelByLabel(ctx, label)
	example.Fatal(err)
	example.Print("cancelled", n)
}

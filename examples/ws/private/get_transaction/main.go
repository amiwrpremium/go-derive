// Fetches one transaction record by id over WebSocket.
// Required env: DERIVE_TX_ID.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	txID := os.Getenv("DERIVE_TX_ID")
	if txID == "" {
		log.Fatal("DERIVE_TX_ID required")
	}
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	tx, err := c.GetTransaction(ctx, txID)
	example.Fatal(err)
	example.Print("status", string(tx.Status))
}

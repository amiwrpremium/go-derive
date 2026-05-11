// Fetches one transaction record by id. Required env: DERIVE_TX_ID.
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
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	tx, err := c.GetTransaction(ctx, txID)
	example.Fatal(err)
	example.Print("status", string(tx.Status))
	example.Print("tx_hash", tx.TransactionHash)
}

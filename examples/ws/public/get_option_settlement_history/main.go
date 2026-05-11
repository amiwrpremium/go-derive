// Lists platform-wide option settlements over WebSocket.
// Required env: DERIVE_CURRENCY.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	currency := os.Getenv("DERIVE_CURRENCY")
	if currency == "" {
		log.Fatal("DERIVE_CURRENCY required")
	}
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPublic(ctx)
	defer c.Close()

	settlements, _, err := c.GetPublicOptionSettlementHistory(ctx, map[string]any{
		"currency": currency,
	})
	example.Fatal(err)
	example.Print("count", len(settlements))
}

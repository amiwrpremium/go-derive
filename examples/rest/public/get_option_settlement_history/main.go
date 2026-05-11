// Lists platform-wide option settlements. Public — no auth needed.
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
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	settlements, page, err := c.GetPublicOptionSettlementHistory(ctx, map[string]any{
		"currency": currency,
	})
	example.Fatal(err)
	example.Print("count", len(settlements))
	example.Print("page count", page.Count)
}

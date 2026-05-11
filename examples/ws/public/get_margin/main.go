// Simulates a margin calculation for an arbitrary subaccount over
// WebSocket. Public — no auth required.
// Required env: DERIVE_SUBACCOUNT.
package main

import (
	"log"
	"os"
	"strconv"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	subStr := os.Getenv("DERIVE_SUBACCOUNT")
	if subStr == "" {
		log.Fatal("DERIVE_SUBACCOUNT required")
	}
	sub, err := strconv.ParseInt(subStr, 10, 64)
	example.Fatal(err)

	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPublic(ctx)
	defer c.Close()

	res, err := c.GetPublicMargin(ctx, map[string]any{
		"subaccount_id":    sub,
		"simulated_trades": []any{},
	})
	example.Fatal(err)
	example.Print("post_initial_margin", res.PostInitialMargin.String())
}

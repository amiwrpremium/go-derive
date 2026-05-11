// Simulates a margin calculation for an arbitrary (possibly
// hypothetical) subaccount. Public — no auth required. Pass the
// portfolio shape via params per the docs at
// https://docs.derive.xyz/reference/public-get_margin.
//
// Required env: DERIVE_SUBACCOUNT (the subaccount to simulate
// against).
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

	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	res, err := c.GetPublicMargin(ctx, map[string]any{
		"subaccount_id":    sub,
		"simulated_trades": []any{},
	})
	example.Fatal(err)
	example.Print("pre_initial_margin", res.PreInitialMargin.String())
	example.Print("post_initial_margin", res.PostInitialMargin.String())
	example.Print("is_valid_trade", res.IsValidTrade)
}

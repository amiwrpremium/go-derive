// Simulates a margin calculation for the configured subaccount over
// WebSocket.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	res, err := c.GetMargin(ctx)
	example.Fatal(err)
	example.Print("subaccount_id", res.SubaccountID)
	example.Print("pre_initial_margin", res.PreInitialMargin.String())
	example.Print("is_valid_trade", res.IsValidTrade)
}

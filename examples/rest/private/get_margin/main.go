// Simulates a margin calculation against the configured subaccount
// and reports pre/post initial- and maintenance-margin values.
// No simulated trade — calling without params reports the current
// subaccount's margin snapshot.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	res, err := c.GetMargin(ctx)
	example.Fatal(err)
	example.Print("subaccount_id", res.SubaccountID)
	example.Print("pre_initial_margin", res.PreInitialMargin.String())
	example.Print("post_initial_margin", res.PostInitialMargin.String())
	example.Print("is_valid_trade", res.IsValidTrade)
}

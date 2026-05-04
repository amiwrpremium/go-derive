// Fetches the configured subaccount's past liquidation events.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	raw, err := c.GetLiquidationHistory(ctx, nil)
	example.Fatal(err)
	example.Print("liquidation history bytes", len(raw))
}

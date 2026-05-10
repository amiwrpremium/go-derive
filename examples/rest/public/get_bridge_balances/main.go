// Lists every bridge / cross-chain balance the engine tracks.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	balances, err := c.GetBridgeBalances(ctx)
	example.Fatal(err)
	example.Print("bridges", len(balances))
	for i, b := range balances {
		if i >= 5 {
			break
		}
		example.Print("bridge", b.Name)
		example.Print("  chain_id", b.ChainID)
		example.Print("  balance", b.Balance.String())
	}
}

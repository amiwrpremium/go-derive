// Lists every registered vault ERC-20 pool.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	pools, err := c.GetVaultPools(ctx)
	example.Fatal(err)
	example.Print("pools", len(pools))
	for i, p := range pools {
		if i >= 5 {
			break
		}
		example.Print("pool", p.Name)
		example.Print("  chain_id", p.ChainID)
	}
}

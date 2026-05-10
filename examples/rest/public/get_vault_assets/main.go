// Lists every ERC-20 asset tracked by Derive's vault orderbook.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	assets, err := c.GetVaultAssets(ctx)
	example.Fatal(err)
	example.Print("vault assets", len(assets))
	for i, a := range assets {
		if i >= 5 {
			break
		}
		example.Print("asset", a.Name)
		example.Print("  chain_id", a.ChainID)
	}
}

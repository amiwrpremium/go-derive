// Lists summary statistics (TVL, total supply, last-trade
// subaccount value) for every Derive vault.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	stats, err := c.GetVaultStatistics(ctx)
	example.Fatal(err)
	example.Print("vault count", len(stats))
	for i, v := range stats {
		if i >= 5 {
			break
		}
		example.Print("vault", v.VaultName)
		example.Print("  usd_tvl", v.USDTVL.String())
		example.Print("  total_supply", v.TotalSupply.String())
	}
}

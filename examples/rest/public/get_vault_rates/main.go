// Returns the engine's current view of one basis vault's rate
// components. Optional env: DERIVE_VAULT_TYPE (e.g. "lbtc",
// "weeth").
package main

import (
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	rates, err := c.GetVaultRates(ctx, os.Getenv("DERIVE_VAULT_TYPE"))
	example.Fatal(err)
	example.Print("rate", rates.Rate.String())
	example.Print("total_rate", rates.TotalRate.String())
	example.Print("funding_rate", rates.FundingRate.String())
	example.Print("interest_rate", rates.InterestRate.String())
}

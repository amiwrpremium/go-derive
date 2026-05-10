// Lists one wallet's vault-token holdings across every Derive vault.
//
// Set DERIVE_WALLET to the smart-contract wallet address you want to
// query; the result is empty when the wallet has no vault deposits.
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

	params := map[string]any{}
	if w := os.Getenv("DERIVE_WALLET"); w != "" {
		params["wallet"] = w
	}

	balances, err := c.GetVaultBalances(ctx, params)
	example.Fatal(err)
	example.Print("balance count", len(balances))
	for i, b := range balances {
		if i >= 5 {
			break
		}
		example.Print("vault", b.Name)
		example.Print("  amount", b.Amount.String())
		example.Print("  chain_id", b.ChainID)
	}
}

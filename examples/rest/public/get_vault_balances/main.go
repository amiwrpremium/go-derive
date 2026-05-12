// Lists one wallet's vault-token holdings across every Derive vault.
//
// Set DERIVE_WALLET to the smart-contract wallet address (or
// DERIVE_SMART_CONTRACT_OWNER to the EOA that owns it). At least
// one must be set; the result is empty when the wallet has no
// vault deposits.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	wallet := os.Getenv("DERIVE_WALLET")
	owner := os.Getenv("DERIVE_SMART_CONTRACT_OWNER")
	if wallet == "" && owner == "" {
		log.Fatal("DERIVE_WALLET or DERIVE_SMART_CONTRACT_OWNER required")
	}
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	balances, err := c.GetVaultBalances(ctx, wallet, owner)
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

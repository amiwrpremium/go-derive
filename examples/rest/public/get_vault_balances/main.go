// Lists one wallet's vault-token holdings across every Derive vault.
//
// Set DERIVE_WALLET to the smart-contract wallet address (or
// DERIVE_SMART_CONTRACT_OWNER to the EOA that owns it). At least
// one must be set; the result is empty when the wallet has no
// vault deposits.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amiwrpremium/go-derive/pkg/rest"
)

func main() {
	restNetwork := rest.WithTestnet()
	if os.Getenv("DERIVE_NETWORK") == "mainnet" {
		restNetwork = rest.WithMainnet()
	}
	c, err := rest.New(restNetwork)
	if err != nil {
		log.Fatalf("rest.New: %v", err)
	}
	defer c.Close()
	wallet := os.Getenv("DERIVE_WALLET")
	owner := os.Getenv("DERIVE_SMART_CONTRACT_OWNER")
	if wallet == "" && owner == "" {
		log.Fatal("DERIVE_WALLET or DERIVE_SMART_CONTRACT_OWNER required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	balances, err := c.GetVaultBalances(ctx, wallet, owner)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "balance count:", len(balances))
	for i, b := range balances {
		if i >= 5 {
			break
		}
		fmt.Printf("%-30s %v\n", "vault:", b.Name)
		fmt.Printf("%-30s %v\n", "  amount:", b.Amount.String())
		fmt.Printf("%-30s %v\n", "  chain_id:", b.ChainID)
	}
}

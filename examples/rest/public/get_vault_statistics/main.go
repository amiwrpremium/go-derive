// Lists summary statistics (TVL, total supply, last-trade
// subaccount value) for every Derive vault.
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stats, err := c.GetVaultStatistics(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "vault count:", len(stats))
	for i, v := range stats {
		if i >= 5 {
			break
		}
		fmt.Printf("%-30s %v\n", "vault:", v.VaultName)
		fmt.Printf("%-30s %v\n", "  usd_tvl:", v.USDTVL.String())
		fmt.Printf("%-30s %v\n", "  total_supply:", v.TotalSupply.String())
	}
}

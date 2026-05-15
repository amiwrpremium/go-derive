// Returns the engine's current view of one basis vault's rate
// components. Optional env: DERIVE_VAULT_TYPE (e.g. "lbtc",
// "weeth").
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amiwrpremium/go-derive/pkg/rest"
	"github.com/amiwrpremium/go-derive/pkg/types"
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

	rates, err := c.GetVaultRates(ctx, types.VaultRatesQuery{VaultType: os.Getenv("DERIVE_VAULT_TYPE")})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "rate:", rates.Rate.String())
	fmt.Printf("%-30s %v\n", "total_rate:", rates.TotalRate.String())
	fmt.Printf("%-30s %v\n", "funding_rate:", rates.FundingRate.String())
	fmt.Printf("%-30s %v\n", "interest_rate:", rates.InterestRate.String())
}

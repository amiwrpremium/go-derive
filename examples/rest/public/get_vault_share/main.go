// Fetches per-block snapshots of one vault token's price-per-share
// over the last 24 hours.
//
// Required: DERIVE_VAULT_NAME (run get_vault_statistics first to
// discover available vault names).
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
	name := os.Getenv("DERIVE_VAULT_NAME")
	if name == "" {
		log.Fatal("DERIVE_VAULT_NAME required (run get_vault_statistics to discover names)")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	end := time.Now().Unix()
	start := end - 24*3600

	shares, page, err := c.GetVaultShare(ctx, types.VaultShareQuery{
		VaultName: name,
		FromSec:   start,
		ToSec:     end,
	}, types.PageRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "snapshot count:", len(shares))
	fmt.Printf("%-30s %v\n", "total pages:", page.NumPages)
	for i, s := range shares {
		if i >= 3 {
			break
		}
		fmt.Printf("%-30s %v\n", "block:", s.BlockNumber)
		fmt.Printf("%-30s %v\n", "  usd_value:", s.USDValue.String())
		fmt.Printf("%-30s %v\n", "  base_value:", s.BaseValue.String())
	}
}

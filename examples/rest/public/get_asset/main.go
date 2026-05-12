// Fetches one Asset record by name. Required env:
// DERIVE_ASSET_NAME.
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
	name := os.Getenv("DERIVE_ASSET_NAME")
	if name == "" {
		log.Fatal("DERIVE_ASSET_NAME required (e.g. USDC)")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	a, err := c.GetAsset(ctx, name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "asset_name:", a.AssetName)
	fmt.Printf("%-30s %v\n", "asset_type:", string(a.AssetType))
	fmt.Printf("%-30s %v\n", "currency:", a.Currency)
	fmt.Printf("%-30s %v\n", "is_collateral:", a.IsCollateral)
	fmt.Printf("%-30s %v\n", "is_position:", a.IsPosition)
}

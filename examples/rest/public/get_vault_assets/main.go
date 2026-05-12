// Lists every ERC-20 asset tracked by Derive's vault orderbook.
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

	assets, err := c.GetVaultAssets(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "vault assets:", len(assets))
	for i, a := range assets {
		if i >= 5 {
			break
		}
		fmt.Printf("%-30s %v\n", "asset:", a.Name)
		fmt.Printf("%-30s %v\n", "  chain_id:", a.ChainID)
	}
}

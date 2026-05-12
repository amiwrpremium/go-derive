// Lists every bridge / cross-chain balance the engine tracks.
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

	balances, err := c.GetBridgeBalances(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "bridges:", len(balances))
	for i, b := range balances {
		if i >= 5 {
			break
		}
		fmt.Printf("%-30s %v\n", "bridge:", b.Name)
		fmt.Printf("%-30s %v\n", "  chain_id:", b.ChainID)
		fmt.Printf("%-30s %v\n", "  balance:", b.Balance.String())
	}
}

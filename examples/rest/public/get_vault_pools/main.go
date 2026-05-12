// Lists every registered vault ERC-20 pool.
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

	pools, err := c.GetVaultPools(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "pools:", len(pools))
	for i, p := range pools {
		if i >= 5 {
			break
		}
		fmt.Printf("%-30s %v\n", "pool:", p.Name)
		fmt.Printf("%-30s %v\n", "  chain_id:", p.ChainID)
	}
}

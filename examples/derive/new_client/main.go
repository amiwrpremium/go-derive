// Constructs the top-level facade with the configured network.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/pkg/derive"
)

func main() {
	network := derive.WithTestnet()
	if os.Getenv("DERIVE_NETWORK") == "mainnet" {
		network = derive.WithMainnet()
	}
	c, err := derive.NewClient(network)
	if err != nil {
		log.Fatalf("derive.NewClient: %v", err)
	}
	defer c.Close()
	fmt.Printf("%-30s %v\n", "network:", c.Network().Network)
	fmt.Printf("%-30s %v\n", "chain id:", c.Network().ChainID)
}

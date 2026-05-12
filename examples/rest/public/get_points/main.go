// One wallet's points record for one program. Required env:
// DERIVE_PROGRAM_NAME and DERIVE_WALLET.
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
	name := os.Getenv("DERIVE_PROGRAM_NAME")
	wallet := os.Getenv("DERIVE_WALLET")
	if name == "" || wallet == "" {
		log.Fatal("DERIVE_PROGRAM_NAME and DERIVE_WALLET required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := c.GetPoints(ctx, name, wallet)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "total_notional_volume:", res.TotalNotionalVolume.String())
	fmt.Printf("%-30s %v\n", "points (raw bytes):", len(res.Points))
}

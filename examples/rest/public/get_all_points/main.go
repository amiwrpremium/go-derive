// Program-wide points snapshot. Required env: DERIVE_PROGRAM_NAME.
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
	if name == "" {
		log.Fatal("DERIVE_PROGRAM_NAME required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := c.GetAllPoints(ctx, name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "total_users:", res.TotalUsers)
	fmt.Printf("%-30s %v\n", "total_notional_volume:", res.TotalNotionalVolume.String())
	fmt.Printf("%-30s %v\n", "points (raw bytes):", len(res.Points))
}

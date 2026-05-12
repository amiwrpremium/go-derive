// Uses only the c.WS client from the facade — connects and fetches a ticker.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := c.WS.Connect(ctx); err != nil {
		log.Fatal(err)
	}

	instrument := os.Getenv("DERIVE_INSTRUMENT")
	if instrument == "" {
		instrument = "BTC-PERP"
	}
	tk, err := c.WS.GetTicker(ctx, instrument)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "mark:", tk.MarkPrice)
}

// Fetches a ticker (top-of-book + marks).
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
	instrument := os.Getenv("DERIVE_INSTRUMENT")
	if instrument == "" {
		instrument = "BTC-PERP"
	}

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

	t, err := c.GetTicker(ctx, instrument)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "instrument:", t.InstrumentName)
	fmt.Printf("%-30s %v\n", "bid:", t.BestBidPrice)
	fmt.Printf("%-30s %v\n", "ask:", t.BestAskPrice)
	fmt.Printf("%-30s %v\n", "mark:", t.MarkPrice)
}

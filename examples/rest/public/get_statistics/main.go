// Fetches rolling daily / lifetime statistics (volume, fees, trades, OI)
// for one instrument and prints the headline numbers.
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

	s, err := c.GetStatistics(ctx, instrument)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%-30s %v\n", "daily notional volume:", s.DailyNotionalVolume.String())
	fmt.Printf("%-30s %v\n", "daily trades:", s.DailyTrades)
	fmt.Printf("%-30s %v\n", "daily fees:", s.DailyFees.String())
	fmt.Printf("%-30s %v\n", "total trades:", s.TotalTrades)
	fmt.Printf("%-30s %v\n", "open interest:", s.OpenInterest.String())
}

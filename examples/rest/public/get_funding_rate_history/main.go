// Fetches historical funding-rate prints for one perpetual instrument
// and prints the most recent few entries.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amiwrpremium/go-derive/pkg/rest"
	"github.com/amiwrpremium/go-derive/pkg/types"
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

	history, err := c.GetFundingRateHistory(ctx, types.FundingRateHistoryQuery{
		InstrumentName: instrument,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%-30s %v\n", "entries returned:", len(history))
	for i, e := range history {
		if i >= 5 {
			break
		}
		fmt.Printf("%-30s %v\n", "rate at ms:", e.Timestamp.Millis())
		fmt.Printf("%-30s %v\n", "  funding_rate:", e.FundingRate.String())
	}
}

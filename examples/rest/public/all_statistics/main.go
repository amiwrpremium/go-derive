// Aggregate (currency, instrument_type) statistics across every
// instrument.
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

	stats, err := c.GetAllStatistics(ctx, types.AllStatisticsQuery{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "rows:", len(stats))
	for i, s := range stats {
		if i >= 5 {
			break
		}
		fmt.Printf("%-30s %v\n", "tuple:", s.Currency+"/"+s.InstrumentType)
		fmt.Printf("%-30s %v\n", "  daily_notional_volume:", s.DailyNotionalVolume.String())
		fmt.Printf("%-30s %v\n", "  open_interest:", s.OpenInterest.String())
	}
}

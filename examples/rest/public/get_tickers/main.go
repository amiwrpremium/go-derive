// Bulk-fetches the per-instrument ticker snapshot for every perp
// instrument. Returns a map keyed by instrument name. Counterpart
// to GetTicker (single-instrument).
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amiwrpremium/go-derive/pkg/enums"
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

	tickers, err := c.GetTickers(ctx, types.TickersQuery{InstrumentType: enums.InstrumentTypePerp})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "ticker count:", len(tickers))
	i := 0
	for name, t := range tickers {
		if i >= 3 {
			break
		}
		fmt.Printf("%-30s %v\n", "instrument:", name)
		fmt.Printf("%-30s %v\n", "  mark:", t.MarkPrice.String())
		fmt.Printf("%-30s %v\n", "  best bid:", t.BestBidPrice.String())
		fmt.Printf("%-30s %v\n", "  best ask:", t.BestAskPrice.String())
		i++
	}
}

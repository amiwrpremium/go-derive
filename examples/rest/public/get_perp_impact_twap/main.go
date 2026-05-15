// Fetches the time-weighted average impact price for one currency's
// perpetual book over the last hour.
//
// Derive returns the TWAPs of three differences (mid, ask-impact,
// bid-impact each minus spot) as decimal strings. When the book is
// quiet, all three are "0".
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

	now := time.Now()
	twap, err := c.GetPerpImpactTWAP(ctx, types.PerpImpactTWAPQuery{Currency: "BTC", StartTime: now.Add(-time.Hour).UnixMilli(), EndTime: now.UnixMilli()})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%-30s %v\n", "currency:", twap.Currency)
	fmt.Printf("%-30s %v\n", "mid diff TWAP:", twap.MidPriceDiffTWAP.String())
	fmt.Printf("%-30s %v\n", "ask impact TWAP:", twap.AskImpactDiffTWAP.String())
	fmt.Printf("%-30s %v\n", "bid impact TWAP:", twap.BidImpactDiffTWAP.String())
}

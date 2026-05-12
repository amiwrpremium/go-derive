// Fetches OHLC candles for one currency's spot feed over the last
// hour at one-minute resolution. Same per-bar shape as
// get_index_chart_data, sourced from the spot feed instead.
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
	currency := os.Getenv("DERIVE_CURRENCY")
	if currency == "" {
		currency = "BTC"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	end := time.Now()
	start := end.Add(-time.Hour)

	cur, candles, err := c.GetSpotFeedHistoryCandles(ctx, types.SpotFeedHistoryCandlesQuery{
		HistoryWindow: types.HistoryWindow{
			StartTimestamp: types.NewMillisTime(start),
			EndTimestamp:   types.NewMillisTime(end),
		},
		Currency:  currency,
		PeriodSec: 60,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "currency:", cur)
	fmt.Printf("%-30s %v\n", "candles:", len(candles))
	for i, k := range candles {
		if i >= 3 {
			break
		}
		fmt.Printf("%-30s %v\n", "at ms:", k.Timestamp.Millis())
		fmt.Printf("%-30s %v\n", "  spot price:", k.Price.String())
	}
}

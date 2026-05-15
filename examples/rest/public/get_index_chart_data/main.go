// Fetches OHLC candles for one currency's index price feed over the
// last hour at one-minute resolution.
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

	candles, err := c.GetIndexChartData(ctx, types.IndexChartQuery{
		HistoryWindow: types.HistoryWindow{
			StartTimestamp: types.MillisTimeFromTime(start),
			EndTimestamp:   types.MillisTimeFromTime(end),
		},
		Currency:  currency,
		PeriodSec: 60,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "candles:", len(candles))
	for i, k := range candles {
		if i >= 3 {
			break
		}
		fmt.Printf("%-30s %v\n", "at ms:", k.Timestamp.Millis())
		fmt.Printf("%-30s %v\n", "  open / close:", k.OpenPrice.String()+" / "+k.ClosePrice.String())
	}
}

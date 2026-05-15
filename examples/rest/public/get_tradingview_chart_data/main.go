// Fetches TradingView-format OHLC bars for one instrument over the
// last hour at one-minute resolution. Each bar carries volume in
// both contracts and USD.
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

	end := time.Now()
	start := end.Add(-time.Hour)

	bars, err := c.GetTradingViewChartData(ctx, types.TradingViewChartQuery{
		HistoryWindow: types.HistoryWindow{
			StartTimestamp: types.MillisTimeFromTime(start),
			EndTimestamp:   types.MillisTimeFromTime(end),
		},
		InstrumentName: instrument,
		PeriodSec:      60,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "bars:", len(bars))
	for i, b := range bars {
		if i >= 3 {
			break
		}
		fmt.Printf("%-30s %v\n", "at ms:", b.Timestamp.Millis())
		fmt.Printf("%-30s %v\n", "  ohlc:", b.OpenPrice.String()+" "+b.HighPrice.String()+" "+b.LowPrice.String()+" "+b.ClosePrice.String())
		fmt.Printf("%-30s %v\n", "  volume_usd:", b.VolumeUSD.String())
	}
}

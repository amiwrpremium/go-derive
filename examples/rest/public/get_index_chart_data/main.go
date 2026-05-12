// Fetches OHLC candles for one currency's index price feed over the
// last hour at one-minute resolution.
package main

import (
	"os"
	"time"

	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	currency := os.Getenv("DERIVE_CURRENCY")
	if currency == "" {
		currency = "BTC"
	}
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	end := time.Now()
	start := end.Add(-time.Hour)

	candles, err := c.GetIndexChartData(ctx, types.IndexChartQuery{
		HistoryWindow: types.HistoryWindow{
			StartTimestamp: types.NewMillisTime(start),
			EndTimestamp:   types.NewMillisTime(end),
		},
		Currency:  currency,
		PeriodSec: 60,
	})
	example.Fatal(err)
	example.Print("candles", len(candles))
	for i, k := range candles {
		if i >= 3 {
			break
		}
		example.Print("at ms", k.Timestamp.Millis())
		example.Print("  open / close", k.OpenPrice.String()+" / "+k.ClosePrice.String())
	}
}

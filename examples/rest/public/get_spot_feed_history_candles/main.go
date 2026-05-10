// Fetches OHLC candles for one currency's spot feed over the last
// hour at one-minute resolution. Same per-bar shape as
// get_index_chart_data, sourced from the spot feed instead.
package main

import (
	"os"
	"time"

	"github.com/amiwrpremium/go-derive/examples/example"
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

	end := time.Now().Unix()
	start := end - 3600

	cur, candles, err := c.GetSpotFeedHistoryCandles(ctx, map[string]any{
		"currency":        currency,
		"start_timestamp": start,
		"end_timestamp":   end,
		"period":          "60",
	})
	example.Fatal(err)
	example.Print("currency", cur)
	example.Print("candles", len(candles))
	for i, k := range candles {
		if i >= 3 {
			break
		}
		example.Print("at ms", k.Timestamp.Millis())
		example.Print("  spot price", k.Price.String())
	}
}

// Fetches OHLC candles for one currency's index price feed over the
// last hour at one-minute resolution.
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

	candles, err := c.GetIndexChartData(ctx, map[string]any{
		"currency":        currency,
		"start_timestamp": start,
		"end_timestamp":   end,
		"period":          "60",
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

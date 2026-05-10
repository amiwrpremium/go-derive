// Fetches TradingView-format OHLC bars for one instrument over the
// last hour at one-minute resolution. Each bar carries volume in
// both contracts and USD.
package main

import (
	"time"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	end := time.Now().Unix()
	start := end - 3600

	bars, err := c.GetTradingViewChartData(ctx, map[string]any{
		"instrument_name": example.Instrument(),
		"start_timestamp": start,
		"end_timestamp":   end,
		"period":          "60",
	})
	example.Fatal(err)
	example.Print("bars", len(bars))
	for i, b := range bars {
		if i >= 3 {
			break
		}
		example.Print("at ms", b.Timestamp.Millis())
		example.Print("  ohlc",
			b.OpenPrice.String()+" "+b.HighPrice.String()+" "+b.LowPrice.String()+" "+b.ClosePrice.String())
		example.Print("  volume_usd", b.VolumeUSD.String())
	}
}

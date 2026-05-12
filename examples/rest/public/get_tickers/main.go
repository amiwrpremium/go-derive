// Bulk-fetches the per-instrument ticker snapshot for every perp
// instrument. Returns a map keyed by instrument name. Counterpart
// to GetTicker (single-instrument).
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/enums"
)

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	tickers, err := c.GetTickers(ctx, enums.InstrumentTypePerp, "", 0)
	example.Fatal(err)
	example.Print("ticker count", len(tickers))
	i := 0
	for name, t := range tickers {
		if i >= 3 {
			break
		}
		example.Print("instrument", name)
		example.Print("  mark", t.MarkPrice.String())
		example.Print("  best bid", t.BestBidPrice.String())
		example.Print("  best ask", t.BestAskPrice.String())
		i++
	}
}

// Fetches the time-weighted average impact price for one currency's
// perpetual book over the last hour.
//
// Derive returns the TWAPs of three differences (mid, ask-impact,
// bid-impact each minus spot) as decimal strings. When the book is
// quiet, all three are "0".
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

	now := time.Now()
	twap, err := c.GetPerpImpactTWAP(ctx, "BTC", now.Add(-time.Hour).UnixMilli(), now.UnixMilli())
	example.Fatal(err)

	example.Print("currency", twap.Currency)
	example.Print("mid diff TWAP", twap.MidPriceDiffTWAP.String())
	example.Print("ask impact TWAP", twap.AskImpactDiffTWAP.String())
	example.Print("bid impact TWAP", twap.BidImpactDiffTWAP.String())
}

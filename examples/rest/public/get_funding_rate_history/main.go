// Fetches historical funding-rate prints for one perpetual instrument
// and prints the most recent few entries.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	history, err := c.GetFundingRateHistory(ctx, types.FundingRateHistoryQuery{
		InstrumentName: example.Instrument(),
	})
	example.Fatal(err)

	example.Print("entries returned", len(history))
	for i, e := range history {
		if i >= 5 {
			break
		}
		example.Print("rate at ms", e.Timestamp.Millis())
		example.Print("  funding_rate", e.FundingRate.String())
	}
}

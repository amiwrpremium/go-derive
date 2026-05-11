// Fetches funding payments received / paid by the configured subaccount.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	events, page, err := c.GetFundingHistory(ctx, types.FundingHistoryQuery{}, types.PageRequest{})
	example.Fatal(err)
	example.Print("event count", len(events))
	example.Print("page count", page.Count)
	if len(events) > 0 {
		e := events[0]
		example.Print("first instrument", e.InstrumentName)
		example.Print("first funding", e.Funding.String())
	}
}

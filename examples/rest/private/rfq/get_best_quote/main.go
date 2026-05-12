// Returns the best quote currently available for a hypothetical
// RFQ shape — the SDK's pre-flight for "what could I get if I
// asked for this?" without actually submitting the RFQ.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	res, err := c.RFQGetBestQuote(ctx, types.BestQuoteInput{
		Direction: enums.DirectionBuy,
		Legs: []types.RFQLeg{
			{
				InstrumentName: example.Instrument(),
				Direction:      enums.DirectionBuy,
				Amount:         types.MustDecimal("0.1"),
			},
		},
	})
	example.Fatal(err)
	example.Print("is_valid", res.IsValid)
	example.Print("estimated_total_cost", res.EstimatedTotalCost.String())
	if res.BestQuote != nil {
		example.Print("best_quote_id", res.BestQuote.QuoteID)
	}
}

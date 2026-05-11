// Returns the best quote available for a hypothetical RFQ shape
// over WebSocket. Pre-flight without submitting the RFQ.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

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
}

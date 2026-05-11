// Returns the best quote available for a hypothetical RFQ shape
// over WebSocket. Pre-flight without submitting the RFQ.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/enums"
)

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	res, err := c.RFQGetBestQuote(ctx, map[string]any{
		"direction": string(enums.DirectionBuy),
		"legs": []map[string]any{
			{
				"instrument_name": example.Instrument(),
				"direction":       string(enums.DirectionBuy),
				"amount":          "0.1",
			},
		},
	})
	example.Fatal(err)
	example.Print("is_valid", res.IsValid)
	example.Print("estimated_total_cost", res.EstimatedTotalCost.String())
}

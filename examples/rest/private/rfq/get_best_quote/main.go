// Returns the best quote currently available for a hypothetical
// RFQ shape — the SDK's pre-flight for "what could I get if I
// asked for this?" without actually submitting the RFQ.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/enums"
)

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

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
	if res.BestQuote != nil {
		example.Print("best_quote_id", res.BestQuote.QuoteID)
	}
}

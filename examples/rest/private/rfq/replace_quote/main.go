// Replaces (cancel + send) one outstanding maker quote in a single
// round trip — the quote-side counterpart to private/replace for
// orders.
//
// This example is illustrative — set DERIVE_RUN_LIVE_ORDERS=1 only
// when you actually want it to run; the SDK doesn't gate
// ReplaceQuote itself. Required fields (legs, signing fields, etc.)
// must be filled in for the call to succeed against the real engine.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	if os.Getenv("DERIVE_RUN_LIVE_ORDERS") != "1" {
		log.Fatal("set DERIVE_RUN_LIVE_ORDERS=1 to actually run replace_quote")
	}
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	res, err := c.ReplaceQuote(ctx, types.ReplaceQuoteInput{
		SendQuoteInput: types.SendQuoteInput{
			RFQID:              "<rfq-id>",
			Direction:          enums.DirectionBuy,
			Legs:               nil,
			MaxFee:             types.MustDecimal("10"),
			Nonce:              0,
			Signature:          "",
			Signer:             example.MustSigner().Owner().Hex(),
			SignatureExpirySec: 0,
		},
		QuoteIDToCancel: "<quote-id>",
	})
	example.Fatal(err)
	example.Print("cancelled_quote", res.CancelledQuote.QuoteID)
	if res.Quote != nil {
		example.Print("new_quote", res.Quote.QuoteID)
	}
	if res.CreateQuoteError != nil {
		example.Print("create_quote_error code", res.CreateQuoteError.Code)
		example.Print("create_quote_error message", res.CreateQuoteError.Message)
	}
}

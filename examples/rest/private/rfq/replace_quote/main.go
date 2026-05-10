// Replaces (cancel + send) one outstanding maker quote in a single
// round trip — the quote-side counterpart to private/replace for
// orders.
//
// This example is illustrative — set DERIVE_RUN_LIVE_ORDERS=1 only
// when you actually want it to run; the SDK doesn't gate
// ReplaceQuote itself. Required params (rfq_id, direction, legs,
// max_fee, signing fields) must be filled in for the call to
// succeed against the real engine.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	if os.Getenv("DERIVE_RUN_LIVE_ORDERS") != "1" {
		log.Fatal("set DERIVE_RUN_LIVE_ORDERS=1 to actually run replace_quote")
	}
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	res, err := c.ReplaceQuote(ctx, map[string]any{
		"rfq_id":             "<rfq-id>",
		"quote_id_to_cancel": "<quote-id>",
		// fill in: direction, legs, max_fee, nonce, signature,
		// signature_expiry_sec, signer
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

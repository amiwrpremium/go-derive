// Submits a quote in response to an existing RFQ — the maker side of
// the RFQ flow. private/send_quote is a low-level pass-through:
// the caller supplies a fully-signed quote payload (signature,
// signer, nonce, signature_expiry_sec) per the docs at
// https://docs.derive.xyz/reference/private-send_quote.
//
// Requires DERIVE_RFQ_ID and DERIVE_RUN_LIVE_ORDERS=1 (since the
// quote, once accepted, may fill).
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	rfqID := os.Getenv("DERIVE_RFQ_ID")
	if rfqID == "" {
		log.Fatal("DERIVE_RFQ_ID required")
	}
	if os.Getenv("DERIVE_RUN_LIVE_ORDERS") != "1" {
		log.Fatal("set DERIVE_RUN_LIVE_ORDERS=1 to actually send a quote")
	}
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	// Fill in legs / signature / signer / nonce / signature_expiry_sec
	// from your signing pipeline; this example only demonstrates the
	// call shape.
	q, err := c.SendQuote(ctx, map[string]any{
		"rfq_id":               rfqID,
		"direction":            "buy",
		"max_total_cost":       "1000",
		"legs":                 []any{},
		"signature":            "",
		"signer":               example.MustSigner().Owner().Hex(),
		"nonce":                0,
		"signature_expiry_sec": 0,
	})
	example.Fatal(err)
	example.Print("quote id", q.QuoteID)
	example.Print("status", q.Status)
}

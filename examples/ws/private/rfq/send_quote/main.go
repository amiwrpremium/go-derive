// Submits a quote in response to an existing RFQ over WebSocket —
// the maker side of the RFQ flow. private/send_quote is a low-level
// pass-through: the caller supplies a fully-signed quote payload.
//
// Requires DERIVE_RFQ_ID and DERIVE_RUN_LIVE_ORDERS=1.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	rfqID := os.Getenv("DERIVE_RFQ_ID")
	if rfqID == "" {
		log.Fatal("DERIVE_RFQ_ID required")
	}
	if os.Getenv("DERIVE_RUN_LIVE_ORDERS") != "1" {
		log.Fatal("set DERIVE_RUN_LIVE_ORDERS=1 to actually send a quote")
	}
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	q, err := c.SendQuote(ctx, types.SendQuoteInput{
		RFQID:              rfqID,
		Direction:          enums.DirectionBuy,
		Legs:               nil,
		MaxFee:             types.MustDecimal("10"),
		Nonce:              0,
		Signature:          "",
		Signer:             example.MustSigner().Owner().Hex(),
		SignatureExpirySec: 0,
	})
	example.Fatal(err)
	example.Print("quote id", q.QuoteID)
}

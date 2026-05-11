// Executes (accepts) a quote received against the configured
// subaccount's RFQ. Requires DERIVE_QUOTE_ID, DERIVE_RFQ_ID and
// DERIVE_RUN_LIVE_ORDERS=1 since execution fills the trade
// immediately.
//
// private/execute_quote takes a fully-signed payload — the caller
// must supply the signature material from their own signing flow.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	quoteID := os.Getenv("DERIVE_QUOTE_ID")
	if quoteID == "" {
		log.Fatal("DERIVE_QUOTE_ID required")
	}
	rfqID := os.Getenv("DERIVE_RFQ_ID")
	if rfqID == "" {
		log.Fatal("DERIVE_RFQ_ID required")
	}
	if os.Getenv("DERIVE_RUN_LIVE_ORDERS") != "1" {
		log.Fatal("set DERIVE_RUN_LIVE_ORDERS=1 to actually execute")
	}
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	res, err := c.ExecuteQuote(ctx, types.ExecuteQuoteInput{
		RFQID:              rfqID,
		QuoteID:            quoteID,
		Direction:          enums.DirectionBuy,
		Legs:               nil,
		MaxFee:             types.MustDecimal("10"),
		Nonce:              0,
		Signature:          "",
		Signer:             example.MustSigner().Owner().Hex(),
		SignatureExpirySec: 0,
	})
	example.Fatal(err)
	example.Print("rfq filled pct", res.RFQFilledPct.String())
}

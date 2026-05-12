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
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/amiwrpremium/go-derive/pkg/auth"
	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/rest"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	subStr := os.Getenv("DERIVE_SUBACCOUNT")
	if subStr == "" {
		log.Fatal("DERIVE_SUBACCOUNT required")
	}
	subaccount, err := strconv.ParseInt(subStr, 10, 64)
	if err != nil {
		log.Fatalf("DERIVE_SUBACCOUNT=%q: %v", subStr, err)
	}
	key := os.Getenv("DERIVE_SESSION_KEY")
	if key == "" {
		log.Fatal("DERIVE_SESSION_KEY required")
	}
	var signer auth.Signer
	if owner := os.Getenv("DERIVE_OWNER"); owner != "" {
		signer, err = auth.NewSessionKeySigner(key, common.HexToAddress(owner))
	} else {
		signer, err = auth.NewLocalSigner(key)
	}
	if err != nil {
		log.Fatalf("signer: %v", err)
	}

	restNetwork := rest.WithTestnet()
	if os.Getenv("DERIVE_NETWORK") == "mainnet" {
		restNetwork = rest.WithMainnet()
	}
	c, err := rest.New(restNetwork, rest.WithSigner(signer), rest.WithSubaccount(subaccount))
	if err != nil {
		log.Fatalf("rest.New: %v", err)
	}
	defer c.Close()
	if os.Getenv("DERIVE_RUN_LIVE_ORDERS") != "1" {
		log.Fatal("set DERIVE_RUN_LIVE_ORDERS=1 to actually run replace_quote")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := c.ReplaceQuote(ctx, types.ReplaceQuoteInput{
		SendQuoteInput: types.SendQuoteInput{
			RFQID:              "<rfq-id>",
			Direction:          enums.DirectionBuy,
			Legs:               nil,
			MaxFee:             types.MustDecimal("10"),
			Nonce:              0,
			Signature:          "",
			Signer:             signer.Owner().Hex(),
			SignatureExpirySec: 0,
		},
		QuoteIDToCancel: "<quote-id>",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "cancelled_quote:", res.CancelledQuote.QuoteID)
	if res.Quote != nil {
		fmt.Printf("%-30s %v\n", "new_quote:", res.Quote.QuoteID)
	}
	if res.CreateQuoteError != nil {
		fmt.Printf("%-30s %v\n", "create_quote_error code:", res.CreateQuoteError.Code)
		fmt.Printf("%-30s %v\n", "create_quote_error message:", res.CreateQuoteError.Message)
	}
}

// Submits a quote in response to an existing RFQ — the maker side
// of the RFQ flow. private/send_quote is a low-level pass-through:
// the caller supplies a fully-signed quote payload (signature,
// signer, nonce, signature_expiry_sec) per the docs at
// https://docs.derive.xyz/reference/private-send_quote.
//
// Requires DERIVE_RFQ_ID and DERIVE_RUN_LIVE_ORDERS=1 (since the
// quote, once accepted, may fill).
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
	rfqID := os.Getenv("DERIVE_RFQ_ID")
	if rfqID == "" {
		log.Fatal("DERIVE_RFQ_ID required")
	}
	if os.Getenv("DERIVE_RUN_LIVE_ORDERS") != "1" {
		log.Fatal("set DERIVE_RUN_LIVE_ORDERS=1 to actually send a quote")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Fill in legs / signature / signer / nonce / signature_expiry_sec
	// from your signing pipeline; this example only demonstrates the
	// call shape.
	q, err := c.SendQuote(ctx, types.SendQuoteInput{
		RFQID:              rfqID,
		Direction:          enums.DirectionBuy,
		Legs:               nil,
		MaxFee:             types.MustDecimal("10"),
		Nonce:              0,
		Signature:          "",
		Signer:             signer.Owner().Hex(),
		SignatureExpirySec: 0,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "quote id:", q.QuoteID)
	fmt.Printf("%-30s %v\n", "status:", q.Status)
}

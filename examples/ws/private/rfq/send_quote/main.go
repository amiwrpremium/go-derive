// Submits a quote in response to an existing RFQ over WebSocket —
// the maker side of the RFQ flow. private/send_quote is a low-level
// pass-through: the caller supplies a fully-signed quote payload.
//
// Requires DERIVE_RFQ_ID and DERIVE_RUN_LIVE_ORDERS=1.
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
	"github.com/amiwrpremium/go-derive/pkg/types"
	"github.com/amiwrpremium/go-derive/pkg/ws"
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	wsNetwork := ws.WithTestnet()
	if os.Getenv("DERIVE_NETWORK") == "mainnet" {
		wsNetwork = ws.WithMainnet()
	}
	c, err := ws.New(wsNetwork, ws.WithSigner(signer), ws.WithSubaccount(subaccount))
	if err != nil {
		log.Fatalf("ws.New: %v", err)
	}
	defer c.Close()
	if err := c.Connect(ctx); err != nil {
		log.Fatalf("ws.Connect: %v", err)
	}
	if err := c.Login(ctx); err != nil {
		log.Fatalf("ws.Login: %v", err)
	}
	rfqID := os.Getenv("DERIVE_RFQ_ID")
	if rfqID == "" {
		log.Fatal("DERIVE_RFQ_ID required")
	}
	if os.Getenv("DERIVE_RUN_LIVE_ORDERS") != "1" {
		log.Fatal("set DERIVE_RUN_LIVE_ORDERS=1 to actually send a quote")
	}
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
}

// Streams the platform-wide auctions.watch feed: every liquidation
// auction Derive is currently running emits per-state-transition events
// (one when an auction begins/updates, one when it ends).
//
// The channel takes no parameters — every subscriber receives the same
// engine-wide stream.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/ws"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	wsNetwork := ws.WithTestnet()
	if os.Getenv("DERIVE_NETWORK") == "mainnet" {
		wsNetwork = ws.WithMainnet()
	}
	c, err := ws.New(wsNetwork)
	if err != nil {
		log.Fatalf("ws.New: %v", err)
	}
	defer c.Close()
	if err := c.Connect(ctx); err != nil {
		log.Fatalf("ws.Connect: %v", err)
	}
	sub, err := c.SubscribeAuctionsWatch(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case ev, ok := <-sub.Updates():
			if !ok {
				return
			}
			fmt.Printf("%-30s %v\n", "subaccount:", ev.SubaccountID)
			fmt.Printf("%-30s %v\n", "  state:", string(ev.State))
			if ev.State == enums.AuctionStateOngoing && ev.Details != nil {
				fmt.Printf("%-30s %v\n", "  currency:", ev.Details.Currency)
				fmt.Printf("%-30s %v\n", "  est_bid_price:", ev.Details.EstimatedBidPrice.String())
				fmt.Printf("%-30s %v\n", "  est_percent_bid:", ev.Details.EstimatedPercentBid.String())
			}
		}
	}
}

// Streams the platform-wide auctions.watch feed: every liquidation
// auction Derive is currently running emits per-state-transition events
// (one when an auction begins/updates, one when it ends).
//
// The channel takes no parameters — every subscriber receives the same
// engine-wide stream. Demonstrates the typed convenience wrapper
// Client.SubscribeAuctionsWatch alongside the generic ws.Subscribe[T]
// (commented out below).
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/enums"
)

func main() {
	ctx, cancel := example.LongTimeout()
	defer cancel()
	c := example.MustWSPublic(ctx)
	defer c.Close()

	// Typed convenience method (preferred):
	sub, err := c.SubscribeAuctionsWatch(ctx)
	example.Fatal(err)
	defer sub.Close()

	// Equivalent generic form — uncomment to swap in:
	//
	//	sub, err := ws.Subscribe[types.AuctionWatchEvent](ctx, c, public.AuctionsWatch{})

	for {
		select {
		case <-ctx.Done():
			return
		case ev, ok := <-sub.Updates():
			if !ok {
				return
			}
			example.Print("subaccount", ev.SubaccountID)
			example.Print("  state", string(ev.State))
			if ev.State == enums.AuctionStateOngoing && ev.Details != nil {
				example.Print("  currency", ev.Details.Currency)
				example.Print("  est_bid_price", ev.Details.EstimatedBidPrice.String())
				example.Print("  est_percent_bid", ev.Details.EstimatedPercentBid.String())
			}
		}
	}
}

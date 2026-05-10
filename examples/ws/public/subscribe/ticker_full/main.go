// Streams the full ticker channel for one instrument — instrument
// metadata bundled with live market data. For the bandwidth-friendly
// compact wire variant, see ws/public/subscribe/ticker.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	ctx, cancel := example.LongTimeout()
	defer cancel()
	c := example.MustWSPublic(ctx)
	defer c.Close()

	sub, err := c.SubscribeTicker(ctx, example.Instrument(), "")
	example.Fatal(err)
	defer sub.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case feed, ok := <-sub.Updates():
			if !ok {
				return
			}
			example.Print("instrument", feed.Ticker.InstrumentName)
			example.Print("mark", feed.Ticker.MarkPrice.String())
			example.Print("best_bid", feed.Ticker.BestBidPrice.String())
			example.Print("best_ask", feed.Ticker.BestAskPrice.String())
		}
	}
}

// Streams the full ticker channel for one instrument — instrument
// metadata bundled with live market data. For the bandwidth-friendly
// compact wire variant, see ws/public/subscribe/ticker.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/channels/public"
	"github.com/amiwrpremium/go-derive/pkg/types"
	"github.com/amiwrpremium/go-derive/pkg/ws"
)

func main() {
	ctx, cancel := example.LongTimeout()
	defer cancel()
	c := example.MustWSPublic(ctx)
	defer c.Close()

	sub, err := ws.Subscribe[types.InstrumentTickerFeed](ctx, c, public.Ticker{
		Instrument: example.Instrument(),
	})
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

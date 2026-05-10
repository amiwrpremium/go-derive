// Streams ticker_slim updates for one instrument — the compact wire
// variant. For the full payload (instrument metadata + live market
// data), see ws/public/subscribe/ticker_full.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	ctx, cancel := example.LongTimeout()
	defer cancel()
	c := example.MustWSPublic(ctx)
	defer c.Close()

	sub, err := c.SubscribeTickerSlim(ctx, example.Instrument(), "")
	example.Fatal(err)
	defer sub.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case t, ok := <-sub.Updates():
			if !ok {
				return
			}
			example.Print("mark", t.Ticker.MarkPrice)
		}
	}
}

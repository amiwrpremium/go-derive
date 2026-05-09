// Streams oracle spot-feed updates for one currency.
//
// Use the spot_feed channel for liquidation monitoring, basis
// calculations, or any risk surface that needs an oracle anchor
// independent of the order book. Each update reports the current price
// + 24-hour-prior reading.
package main

import (
	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/channels/public"
	"github.com/amiwrpremium/go-derive/pkg/ws"
)

func main() {
	ctx, cancel := example.LongTimeout()
	defer cancel()
	c := example.MustWSPublic(ctx)
	defer c.Close()

	sub, err := ws.Subscribe[derive.SpotFeed](ctx, c, public.SpotFeed{Currency: "BTC"})
	example.Fatal(err)
	defer sub.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case sf, ok := <-sub.Updates():
			if !ok {
				return
			}
			if entry, has := sf.Feeds["BTC"]; has {
				example.Print("BTC oracle", entry.Price)
			}
		}
	}
}

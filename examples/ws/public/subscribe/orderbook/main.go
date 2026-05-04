// Streams orderbook updates for one instrument.
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

	// Derive accepts depth ∈ {1, 10}. Group ∈ {1, 10, 100} (price-bucket size).
	sub, err := ws.Subscribe[types.OrderBook](ctx, c, public.OrderBook{
		Instrument: example.Instrument(), Depth: 10,
	})
	example.Fatal(err)
	defer sub.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case ob, ok := <-sub.Updates():
			if !ok {
				return
			}
			example.Print(ob.InstrumentName, len(ob.Bids))
		}
	}
}

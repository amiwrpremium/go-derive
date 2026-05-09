// Streams ticker_slim updates for one instrument.
//
// `ticker_slim` is the only ticker channel Derive currently supports;
// the legacy `ticker` channel was deprecated.
package main

import (
	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/ws"
)

func main() {
	ctx, cancel := example.LongTimeout()
	defer cancel()
	c := example.MustWSPublic(ctx)
	defer c.Close()

	sub, err := ws.Subscribe[derive.TickerSlim](ctx, c, derive.PublicTickerSlim{
		Instrument: example.Instrument(),
	})
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

// Streams public trade prints.
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

	sub, err := ws.Subscribe[[]types.Trade](ctx, c, public.Trades{Instrument: example.Instrument()})
	example.Fatal(err)
	defer sub.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case batch, ok := <-sub.Updates():
			if !ok {
				return
			}
			example.Print("prints", len(batch))
		}
	}
}

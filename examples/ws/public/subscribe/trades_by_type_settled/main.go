// Streams public trades for one (instrument_type, currency) tuple
// filtered to settled fills only — same payload as
// ws/public/subscribe/trades_by_type but filtered server-side.
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

	sub, err := c.SubscribeTradesByTypeWithStatus(ctx, enums.InstrumentTypePerp, "BTC", enums.TxStatusSettled)
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
			example.Print("settled trades in batch", len(batch))
		}
	}
}

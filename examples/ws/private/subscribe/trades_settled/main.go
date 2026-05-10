// Streams the configured subaccount's fills filtered to settled-only —
// same payload as ws/private/subscribe/trades but filtered
// server-side. Useful for makers who only care about settled fills.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/enums"
)

func main() {
	ctx, cancel := example.LongTimeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	sub, err := c.SubscribeSubaccountTradesByStatus(ctx, example.Subaccount(), enums.TxStatusSettled)
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
			example.Print("settled fills in batch", len(batch))
		}
	}
}

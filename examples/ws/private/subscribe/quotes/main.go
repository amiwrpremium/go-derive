// Streams quote updates received against this subaccount's RFQs.
package main

import (
	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	ctx, cancel := example.LongTimeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	sub, err := derive.Subscribe[[]derive.Quote](ctx, c, derive.PrivateQuotes{SubaccountID: example.Subaccount()})
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
			example.Print("quotes", len(batch))
		}
	}
}

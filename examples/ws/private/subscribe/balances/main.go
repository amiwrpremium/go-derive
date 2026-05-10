// Streams subaccount balance updates.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	ctx, cancel := example.LongTimeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	sub, err := c.SubscribeBalances(ctx, example.Subaccount())
	example.Fatal(err)
	defer sub.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case bal, ok := <-sub.Updates():
			if !ok {
				return
			}
			example.Print("equity", bal.SubaccountValue)
		}
	}
}

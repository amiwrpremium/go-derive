// Streams wallet RFQ status updates. RFQs on Derive are wallet-scoped:
// one stream per signer address covers RFQs across every subaccount the
// wallet operates.
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

	wallet := example.MustSigner().Owner().Hex()
	sub, err := derive.Subscribe[[]derive.RFQ](ctx, c, derive.PrivateRFQs{Wallet: wallet})
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
			example.Print("rfqs", len(batch))
		}
	}
}

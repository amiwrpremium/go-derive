// Streams the running best-quote state for every open RFQ on the
// configured subaccount. Note the unusual channel-name format
// ({subaccount_id}.best.quotes per the docs) — see the package
// docs on private.BestQuotes for details.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/channels/private"
	"github.com/amiwrpremium/go-derive/pkg/types"
	"github.com/amiwrpremium/go-derive/pkg/ws"
)

func main() {
	ctx, cancel := example.LongTimeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	sub, err := ws.Subscribe[[]types.BestQuoteFeedEvent](ctx, c, private.BestQuotes{
		SubaccountID: example.Subaccount(),
	})
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
			example.Print("events in batch", len(batch))
			for i, ev := range batch {
				if i >= 3 {
					break
				}
				example.Print("rfq_id", ev.RFQID)
				if ev.Error != nil {
					example.Print("  error", ev.Error.Message)
				} else if ev.Result != nil {
					example.Print("  is_valid", ev.Result.IsValid)
				}
			}
		}
	}
}

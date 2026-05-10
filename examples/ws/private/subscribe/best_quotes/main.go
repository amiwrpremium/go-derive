// Streams the running best-quote state for every open RFQ on the
// configured subaccount. Wire channel: `{subaccount_id}.best.quotes`.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	ctx, cancel := example.LongTimeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	sub, err := c.SubscribeBestQuotes(ctx, example.Subaccount())
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

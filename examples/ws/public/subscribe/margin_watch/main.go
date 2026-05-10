// Streams the platform-wide margin_watch feed: every subaccount whose
// maintenance margin has crossed the watch threshold gets emitted as
// part of a per-timestamp batch.
//
// The channel takes no parameters — every subscriber receives the same
// engine-wide stream. Filter client-side on margin_type / subaccount_id
// if you only care about a subset.
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

	sub, err := ws.Subscribe[[]types.MarginWatch](ctx, c, public.MarginWatch{})
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
			example.Print("at-risk subaccounts in batch", len(batch))
			for i, ev := range batch {
				if i >= 3 {
					break
				}
				example.Print("subaccount", ev.SubaccountID)
				example.Print("  maintenance_margin", ev.MaintenanceMargin.String())
			}
		}
	}
}

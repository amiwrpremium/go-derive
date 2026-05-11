// Subscribes to three private channels for one subaccount and demuxes
// them in a single select loop. The canonical pattern for trading
// processes that need to react to fill, balance, and order-state
// changes from one place without per-channel goroutines.
//
//   - orders   → Client.SubscribeOrders
//   - balances → Client.SubscribeBalances
//   - trades   → Client.SubscribeSubaccountTrades
//
// Position state is not a Derive subscription channel — poll
// `private/get_positions` (REST or WS RPC) when you need it, or
// derive it from the trades feed.
//
// # Drop trade-off
//
// Each sub has its own buffer (default 256, see WithBufferSize) and
// its own drop policy (default DropNewest). If you do heavy work in
// one select arm, the OTHER subs' buffers keep filling while you're
// blocked — and once full, those other subs start dropping events.
// The drop happens on the slow path's NEIGHBOURS, not the slow path
// itself.
//
// For heavy handlers, either spawn one goroutine per Subscription
// (so each handler runs independently) or fan-in via SubscribeInto
// with one shared chan. Register WithErrorHandler to observe drops
// when they happen — wraps as ws.ErrBufferFull.
//
// # Cancellation
//
// Each c.SubscribeX takes a ctx. Cancelling ctx tears each sub down
// (sends unsubscribe, closes Updates). Below we share the same ctx
// across all three subs and the select loop, so one cancel cleans
// everything up. Per-sub lifetimes are also possible — give each a
// child ctx and cancel only that one.
//
// Requires `DERIVE_SESSION_KEY` (or `DERIVE_OWNER_KEY`) plus
// `DERIVE_SUBACCOUNT` in the environment; see the `examples/example`
// helper for the full env-var contract.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	ctx, cancel := example.LongTimeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	subID := example.Subaccount()

	orders, err := c.SubscribeOrders(ctx, subID)
	example.Fatal(err)
	defer orders.Close()

	balances, err := c.SubscribeBalances(ctx, subID)
	example.Fatal(err)
	defer balances.Close()

	trades, err := c.SubscribeSubaccountTrades(ctx, subID)
	example.Fatal(err)
	defer trades.Close()

	example.Print("multiplexing subaccount", subID)
	for {
		select {
		case <-ctx.Done():
			return
		case o, ok := <-orders.Updates():
			if !ok {
				return
			}
			example.Print("orders", len(o))
		case b, ok := <-balances.Updates():
			if !ok {
				return
			}
			example.Print("balance", b)
		case ts, ok := <-trades.Updates():
			if !ok {
				return
			}
			example.Print("trades", len(ts))
		}
	}
}

// Subscribes to three private channels for one subaccount and demuxes
// them in a single select loop. The canonical pattern for trading
// processes that need to react to fill, balance, and order-state
// changes from one place without per-channel goroutines.
//
//   - orders   → pkg/channels/private.Orders
//   - balances → pkg/channels/private.Balances
//   - trades   → pkg/channels/private.Trades
//
// Position state is not a Derive subscription channel — poll
// `private/get_positions` (REST or WS RPC) when you need it, or
// derive it from the trades feed.
//
// Requires `DERIVE_SESSION_KEY` (or `DERIVE_OWNER_KEY`) plus
// `DERIVE_SUBACCOUNT` in the environment; see the `examples/example`
// helper for the full env-var contract.
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

	subID := example.Subaccount()

	orders, err := ws.Subscribe[[]types.Order](ctx, c, private.Orders{SubaccountID: subID})
	example.Fatal(err)
	defer orders.Close()

	balances, err := ws.Subscribe[types.Balance](ctx, c, private.Balances{SubaccountID: subID})
	example.Fatal(err)
	defer balances.Close()

	trades, err := ws.Subscribe[[]types.Trade](ctx, c, private.Trades{SubaccountID: subID})
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

// Subscribes to three private channels for one subaccount and demuxes
// them in a single select loop. The canonical pattern for trading
// processes that need to react to fill, balance, and order-state
// changes from one place without per-channel goroutines.
//
//   - orders   → derive.PrivateOrders
//   - balances → derive.PrivateBalances
//   - trades   → derive.PrivateTrades
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
	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	ctx, cancel := example.LongTimeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	subID := example.Subaccount()

	orders, err := derive.Subscribe[[]derive.Order](ctx, c, derive.PrivateOrders{SubaccountID: subID})
	example.Fatal(err)
	defer orders.Close()

	balances, err := derive.Subscribe[derive.Balance](ctx, c, derive.PrivateBalances{SubaccountID: subID})
	example.Fatal(err)
	defer balances.Close()

	trades, err := derive.Subscribe[[]derive.Trade](ctx, c, derive.PrivateTrades{SubaccountID: subID})
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

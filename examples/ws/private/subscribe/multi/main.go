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
// Requires `DERIVE_SESSION_KEY` (or `DERIVE_OWNER` plus the session
// key) and `DERIVE_SUBACCOUNT` in the environment.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/amiwrpremium/go-derive/pkg/auth"
	"github.com/amiwrpremium/go-derive/pkg/ws"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	subStr := os.Getenv("DERIVE_SUBACCOUNT")
	if subStr == "" {
		log.Fatal("DERIVE_SUBACCOUNT required")
	}
	subaccount, err := strconv.ParseInt(subStr, 10, 64)
	if err != nil {
		log.Fatalf("DERIVE_SUBACCOUNT=%q: %v", subStr, err)
	}
	key := os.Getenv("DERIVE_SESSION_KEY")
	if key == "" {
		log.Fatal("DERIVE_SESSION_KEY required")
	}
	var signer auth.Signer
	if owner := os.Getenv("DERIVE_OWNER"); owner != "" {
		signer, err = auth.NewSessionKeySigner(key, common.HexToAddress(owner))
	} else {
		signer, err = auth.NewLocalSigner(key)
	}
	if err != nil {
		log.Fatalf("signer: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	wsNetwork := ws.WithTestnet()
	if os.Getenv("DERIVE_NETWORK") == "mainnet" {
		wsNetwork = ws.WithMainnet()
	}
	c, err := ws.New(wsNetwork, ws.WithSigner(signer), ws.WithSubaccount(subaccount))
	if err != nil {
		log.Fatalf("ws.New: %v", err)
	}
	defer c.Close()
	if err := c.Connect(ctx); err != nil {
		log.Fatalf("ws.Connect: %v", err)
	}
	if err := c.Login(ctx); err != nil {
		log.Fatalf("ws.Login: %v", err)
	}
	subID := subaccount

	orders, err := c.SubscribeOrders(ctx, subID)
	if err != nil {
		log.Fatal(err)
	}
	defer orders.Close()

	balances, err := c.SubscribeBalances(ctx, subID)
	if err != nil {
		log.Fatal(err)
	}
	defer balances.Close()

	trades, err := c.SubscribeSubaccountTrades(ctx, subID)
	if err != nil {
		log.Fatal(err)
	}
	defer trades.Close()

	fmt.Printf("%-30s %v\n", "multiplexing subaccount:", subID)
	for {
		select {
		case <-ctx.Done():
			return
		case o, ok := <-orders.Updates():
			if !ok {
				return
			}
			fmt.Printf("%-30s %v\n", "orders:", len(o))
		case b, ok := <-balances.Updates():
			if !ok {
				return
			}
			fmt.Printf("%-30s %v\n", "balance:", b)
		case ts, ok := <-trades.Updates():
			if !ok {
				return
			}
			fmt.Printf("%-30s %v\n", "trades:", len(ts))
		}
	}
}

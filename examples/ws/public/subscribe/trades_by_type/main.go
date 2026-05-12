// Streams every public trade for one (instrument_type, currency) pair.
//
// `trades.perp.BTC` covers every BTC perpetual; `trades.option.ETH`
// covers every ETH option. Useful for index-level analytics, where
// subscribing per-instrument would be both noisier and more expensive.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/ws"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	wsNetwork := ws.WithTestnet()
	if os.Getenv("DERIVE_NETWORK") == "mainnet" {
		wsNetwork = ws.WithMainnet()
	}
	c, err := ws.New(wsNetwork)
	if err != nil {
		log.Fatalf("ws.New: %v", err)
	}
	defer c.Close()
	if err := c.Connect(ctx); err != nil {
		log.Fatalf("ws.Connect: %v", err)
	}
	sub, err := c.SubscribeTradesByType(ctx, enums.InstrumentTypePerp, "BTC")
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case batch, ok := <-sub.Updates():
			if !ok {
				return
			}
			fmt.Printf("%-30s %v\n", "trades.perp.BTC:", len(batch))
		}
	}
}

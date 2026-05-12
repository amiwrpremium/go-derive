// Subscribes to three public channels for one instrument and demultiplexes
// them in a single select loop. This is the canonical pattern for trading
// processes that need order-book pressure, top-of-book, and trade prints
// in one place without per-channel goroutines.
//
//   - orderbook (depth 10)        →  Client.SubscribeOrderBook
//   - ticker_slim (1000 ms)       →  Client.SubscribeTickerSlim
//   - trades                      →  Client.SubscribeTrades
//
// Each select arm prints a one-line summary so it's obvious which channel
// fired. Cancel with Ctrl-C; the example exits cleanly via ctx-cancel.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amiwrpremium/go-derive/pkg/ws"
)

func main() {
	instrument := os.Getenv("DERIVE_INSTRUMENT")
	if instrument == "" {
		instrument = "BTC-PERP"
	}

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
	inst := instrument

	ob, err := c.SubscribeOrderBook(ctx, inst, "", 10)
	if err != nil {
		log.Fatal(err)
	}
	defer ob.Close()

	tk, err := c.SubscribeTickerSlim(ctx, inst, "")
	if err != nil {
		log.Fatal(err)
	}
	defer tk.Close()

	tr, err := c.SubscribeTrades(ctx, inst)
	if err != nil {
		log.Fatal(err)
	}
	defer tr.Close()

	fmt.Printf("%-30s %v\n", "multiplexing:", inst)
	for {
		select {
		case <-ctx.Done():
			return
		case b, ok := <-ob.Updates():
			if !ok {
				return
			}
			fmt.Printf("%-30s %v\n", "orderbook:", len(b.Bids))
		case s, ok := <-tk.Updates():
			if !ok {
				return
			}
			fmt.Printf("%-30s %v\n", "ticker:", s.Ticker.MarkPrice)
		case ts, ok := <-tr.Updates():
			if !ok {
				return
			}
			fmt.Printf("%-30s %v\n", "trades:", len(ts))
		}
	}
}

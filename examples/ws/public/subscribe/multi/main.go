// Subscribes to three public channels for one instrument and demultiplexes
// them across three goroutines — one per subscription — so a slow handler
// on one channel cannot starve the others.
//
//   - orderbook (depth 10)        →  Client.SubscribeOrderBook
//   - ticker_slim (1000 ms)       →  Client.SubscribeTickerSlim
//   - trades                      →  Client.SubscribeTrades
//
// # Why one goroutine per subscription
//
// Each Subscription owns its own buffer with its own drop policy. When
// several subscriptions are read from one select loop, a slow arm
// blocks the others; their buffers fill while the goroutine is stuck;
// and they start dropping events per their configured DropPolicy. The
// drop hits the slow path's NEIGHBOURS, not the slow path itself —
// usually the opposite of what the operator wants.
//
// Splitting the consumer into one goroutine per subscription removes
// the cross-channel coupling: a slow handler only ever drops its own
// channel. Pair with WithErrorHandler to observe ws.ErrBufferFull in
// real time when it does happen.
//
// For the inverse pattern — shared back-pressure across multiple
// subscriptions, one consumer — see examples/ws/public/subscribe/multi_into.
//
// Cancel with Ctrl-C; the example exits cleanly via ctx-cancel.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/amiwrpremium/go-derive/pkg/types"
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

	// Surface buffer-full and decode failures so a slow handler is
	// visible instead of silently lossy.
	onErr := func(err error) {
		switch {
		case errors.Is(err, ws.ErrBufferFull):
			log.Printf("warn: buffer-full drop — consumer can't keep up")
		default:
			log.Printf("warn: %v", err)
		}
	}

	ob, err := c.SubscribeOrderBook(ctx, instrument, "", 10, ws.WithErrorHandler(onErr))
	if err != nil {
		log.Fatal(err)
	}
	defer ob.Close()

	tk, err := c.SubscribeTickerSlim(ctx, instrument, "", ws.WithErrorHandler(onErr))
	if err != nil {
		log.Fatal(err)
	}
	defer tk.Close()

	tr, err := c.SubscribeTrades(ctx, instrument, ws.WithErrorHandler(onErr))
	if err != nil {
		log.Fatal(err)
	}
	defer tr.Close()

	fmt.Printf("%-30s %v\n", "multiplexing:", instrument)

	var wg sync.WaitGroup
	wg.Add(3)
	go pump(&wg, ob.Updates(), func(b types.OrderBook) {
		fmt.Printf("%-30s %v\n", "orderbook:", len(b.Bids))
	})
	go pump(&wg, tk.Updates(), func(s types.TickerSlim) {
		fmt.Printf("%-30s %v\n", "ticker:", s.Ticker.MarkPrice)
	})
	go pump(&wg, tr.Updates(), func(ts []types.Trade) {
		fmt.Printf("%-30s %v\n", "trades:", len(ts))
	})
	wg.Wait()
}

// pump drains one subscription's typed channel until it closes. Doing
// this in a dedicated goroutine is what isolates the slow-handler
// failure mode — see the package doc comment above.
func pump[T any](wg *sync.WaitGroup, ch <-chan T, fn func(T)) {
	defer wg.Done()
	for v := range ch {
		fn(v)
	}
}

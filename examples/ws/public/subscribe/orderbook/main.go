// Streams orderbook updates for one instrument and demos every
// SubscribeOption: a larger buffer, DropOldest so the freshest book
// always wins under back-pressure, and an error handler that
// surfaces drops to a counter.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync/atomic"
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
	var dropped atomic.Int64

	// Derive accepts depth in {1, 10}. Group in {1, 10, 100} (price-bucket size).
	//
	// Defaults if no opts are passed: buffer=256, policy=DropNewest,
	// no error handler. The opts here demonstrate every knob.
	sub, err := c.SubscribeOrderBook(ctx, instrument, "", 10,
		ws.WithBufferSize(1024),
		ws.WithDropPolicy(ws.DropOldest),
		ws.WithErrorHandler(func(error) { dropped.Add(1) }),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Close()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%-30s %v\n", "dropped events:", dropped.Load())
			return
		case ob, ok := <-sub.Updates():
			if !ok {
				fmt.Printf("%-30s %v\n", "dropped events:", dropped.Load())
				return
			}
			fmt.Printf("%-30s %v\n", ob.InstrumentName+":", len(ob.Bids))
		}
	}
}

// Streams the full ticker channel for one instrument — instrument
// metadata bundled with live market data. For the bandwidth-friendly
// compact wire variant, see ws/public/subscribe/ticker.
//
// The `ticker` channel was deprecated upstream on 2025-12-01 in favour
// of `ticker_slim`; this example exists to document the full payload
// shape and exercises [ws.Client.SubscribeTicker] deliberately. The
// `//lint:ignore SA1019` comment below silences staticcheck on the
// intentional deprecated call.
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
	//lint:ignore SA1019 example deliberately demonstrates the deprecated full ticker payload; the slim variant is in the sibling ticker/ example
	sub, err := c.SubscribeTicker(ctx, instrument, "")
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case feed, ok := <-sub.Updates():
			if !ok {
				return
			}
			fmt.Printf("%-30s %v\n", "instrument:", feed.Ticker.InstrumentName)
			fmt.Printf("%-30s %v\n", "mark:", feed.Ticker.MarkPrice.String())
			fmt.Printf("%-30s %v\n", "best_bid:", feed.Ticker.BestBidPrice.String())
			fmt.Printf("%-30s %v\n", "best_ask:", feed.Ticker.BestAskPrice.String())
		}
	}
}

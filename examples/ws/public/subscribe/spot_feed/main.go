// Streams oracle spot-feed updates for one currency.
//
// Use the spot_feed channel for liquidation monitoring, basis
// calculations, or any risk surface that needs an oracle anchor
// independent of the order book. Each update reports the current price
// + 24-hour-prior reading.
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
	sub, err := c.SubscribeSpotFeed(ctx, "BTC")
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case sf, ok := <-sub.Updates():
			if !ok {
				return
			}
			if entry, has := sf.Feeds["BTC"]; has {
				fmt.Printf("%-30s %v\n", "BTC oracle:", entry.Price)
			}
		}
	}
}

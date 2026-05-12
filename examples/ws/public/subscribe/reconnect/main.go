// Demonstrates the SDK's auto-reconnect + auto-resubscribe behaviour by
// subscribing to a ticker_slim feed and printing the update count seen in
// each 10-second window for 60 seconds total.
//
// Reconnect is enabled by default (ws.WithReconnect(true)). This example
// keeps that default and runs as a steady-state observer — when the
// connection drops naturally the SDK redials with exponential backoff
// (internal/retry.Backoff), then re-issues `subscribe` for every active
// channel. Updates resume on the same Updates() channel without the caller
// noticing.
//
// To verify the reconnect path manually:
//
//  1. Start the example (it prints `count=N` every 10 s).
//  2. Block testnet from your network for ~5 s, e.g.
//     `echo '127.0.0.1 api-demo.lyra.finance' | sudo tee -a /etc/hosts`
//     then revert.
//  3. Watch the count drop to 0 for one window, then resume.
//
// Disable reconnect via `ws.WithReconnect(false)` to confirm that without
// it the same disruption produces a permanent disconnect.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amiwrpremium/go-derive/pkg/ws"
)

const (
	totalDuration  = 60 * time.Second
	reportInterval = 10 * time.Second
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
	sub, err := c.SubscribeTickerSlim(ctx, instrument, "")
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Close()

	fmt.Printf("%-30s %v\n", "reconnect-demo:", "running 60s")
	deadline := time.Now().Add(totalDuration)
	tick := time.NewTicker(reportInterval)
	defer tick.Stop()

	var count, total int
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%-30s %v\n", "total updates:", total)
			return
		case <-tick.C:
			fmt.Printf("%-30s %v\n", "count (last 10s):", count)
			total += count
			count = 0
			if !time.Now().Before(deadline) {
				fmt.Printf("%-30s %v\n", "total updates:", total)
				return
			}
		case _, ok := <-sub.Updates():
			if !ok {
				fmt.Printf("%-30s %v\n", "subscription closed:", sub.Err())
				return
			}
			count++
		}
	}
}

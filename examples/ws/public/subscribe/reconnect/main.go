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
	"time"

	"github.com/amiwrpremium/go-derive/examples/example"
)

const (
	totalDuration  = 60 * time.Second
	reportInterval = 10 * time.Second
)

func main() {
	ctx, cancel := example.LongTimeout()
	defer cancel()
	c := example.MustWSPublic(ctx)
	defer c.Close()

	sub, err := c.SubscribeTickerSlim(ctx, example.Instrument(), "")
	example.Fatal(err)
	defer sub.Close()

	example.Print("reconnect-demo", "running 60s")
	deadline := time.Now().Add(totalDuration)
	tick := time.NewTicker(reportInterval)
	defer tick.Stop()

	var count, total int
	for {
		select {
		case <-ctx.Done():
			example.Print("total updates", total)
			return
		case <-tick.C:
			example.Print("count (last 10s)", count)
			total += count
			count = 0
			if !time.Now().Before(deadline) {
				example.Print("total updates", total)
				return
			}
		case _, ok := <-sub.Updates():
			if !ok {
				example.Print("subscription closed", sub.Err())
				return
			}
			count++
		}
	}
}

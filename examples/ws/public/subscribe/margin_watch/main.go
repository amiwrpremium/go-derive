// Streams the platform-wide margin_watch feed: every subaccount whose
// maintenance margin has crossed the watch threshold gets emitted as
// part of a per-timestamp batch.
//
// The channel takes no parameters — every subscriber receives the same
// engine-wide stream. Filter client-side on margin_type / subaccount_id
// if you only care about a subset.
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
	sub, err := c.SubscribeMarginWatch(ctx)
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
			fmt.Printf("%-30s %v\n", "at-risk subaccounts in batch:", len(batch))
			for i, ev := range batch {
				if i >= 3 {
					break
				}
				fmt.Printf("%-30s %v\n", "subaccount:", ev.SubaccountID)
				fmt.Printf("%-30s %v\n", "  maintenance_margin:", ev.MaintenanceMargin.String())
			}
		}
	}
}

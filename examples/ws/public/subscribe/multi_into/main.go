// Multiplexes two trade-channel subscriptions into one shared chan
// using SubscribeInto. The caller owns the buffer entirely — useful
// for a fan-in pattern where one consumer goroutine handles every
// subscription on the client.
//
// Required env: DERIVE_INSTRUMENT_A and DERIVE_INSTRUMENT_B (two
// instrument names whose public trades you want to merge).
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amiwrpremium/go-derive/pkg/types"
	"github.com/amiwrpremium/go-derive/pkg/ws"
)

func decodeTrades(raw json.RawMessage) ([]types.Trade, error) {
	var ts []types.Trade
	return ts, json.Unmarshal(raw, &ts)
}

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
	a := os.Getenv("DERIVE_INSTRUMENT_A")
	b := os.Getenv("DERIVE_INSTRUMENT_B")
	if a == "" || b == "" {
		log.Fatal("DERIVE_INSTRUMENT_A and DERIVE_INSTRUMENT_B required")
	}
	// One shared buffer fed by two subscriptions.
	out := make(chan []types.Trade, 256)
	subA, err := ws.SubscribeInto(ctx, c, "trades."+a, decodeTrades, out)
	if err != nil {
		log.Fatal(err)
	}
	defer subA.Close()
	subB, err := ws.SubscribeInto(ctx, c, "trades."+b, decodeTrades, out)
	if err != nil {
		log.Fatal(err)
	}
	defer subB.Close()

	fmt.Printf("%-30s %v\n", "merging:", a+" + "+b)
	for {
		select {
		case <-ctx.Done():
			return
		case batch := <-out:
			if len(batch) == 0 {
				continue
			}
			fmt.Printf("%-30s %v\n", batch[0].InstrumentName+":", len(batch))
		}
	}
}

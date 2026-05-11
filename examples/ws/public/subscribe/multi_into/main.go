// Multiplexes two trade-channel subscriptions into one shared chan
// using SubscribeInto. The caller owns the buffer entirely — useful
// for a fan-in pattern where one consumer goroutine handles every
// subscription on the client.
//
// Required env: DERIVE_INSTRUMENT_A and DERIVE_INSTRUMENT_B (two
// instrument names whose public trades you want to merge).
package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/types"
	"github.com/amiwrpremium/go-derive/pkg/ws"
)

func decodeTrades(raw json.RawMessage) ([]types.Trade, error) {
	var ts []types.Trade
	return ts, json.Unmarshal(raw, &ts)
}

func main() {
	a := os.Getenv("DERIVE_INSTRUMENT_A")
	b := os.Getenv("DERIVE_INSTRUMENT_B")
	if a == "" || b == "" {
		log.Fatal("DERIVE_INSTRUMENT_A and DERIVE_INSTRUMENT_B required")
	}

	ctx, cancel := example.LongTimeout()
	defer cancel()
	c := example.MustWSPublic(ctx)
	defer c.Close()

	// One shared buffer fed by two subscriptions.
	out := make(chan []types.Trade, 256)
	subA, err := ws.SubscribeInto(ctx, c, "trades."+a, decodeTrades, out)
	example.Fatal(err)
	defer subA.Close()
	subB, err := ws.SubscribeInto(ctx, c, "trades."+b, decodeTrades, out)
	example.Fatal(err)
	defer subB.Close()

	example.Print("merging", a+" + "+b)
	for {
		select {
		case <-ctx.Done():
			return
		case batch := <-out:
			if len(batch) == 0 {
				continue
			}
			example.Print(batch[0].InstrumentName, len(batch))
		}
	}
}

// Streams every public trade for one (instrument_type, currency) pair.
//
// `trades.perp.BTC` covers every BTC perpetual; `trades.option.ETH`
// covers every ETH option. Useful for index-level analytics, where
// subscribing per-instrument would be both noisier and more expensive.
package main

import "github.com/amiwrpremium/go-derive"

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/ws"
)

func main() {
	ctx, cancel := example.LongTimeout()
	defer cancel()
	c := example.MustWSPublic(ctx)
	defer c.Close()

	sub, err := ws.Subscribe[[]derive.Trade](ctx, c, derive.PublicTradesByType{
		InstrumentType: derive.InstrumentTypePerp,
		Currency:       "BTC",
	})
	example.Fatal(err)
	defer sub.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case batch, ok := <-sub.Updates():
			if !ok {
				return
			}
			example.Print("trades.perp.BTC", len(batch))
		}
	}
}

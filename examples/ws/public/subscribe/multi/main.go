// Subscribes to three public channels for one instrument and demultiplexes
// them in a single select loop. This is the canonical pattern for trading
// processes that need order-book pressure, top-of-book, and trade prints
// in one place without per-channel goroutines.
//
//   - orderbook (depth 10)        →  Client.SubscribeOrderBook
//   - ticker_slim (1000 ms)       →  Client.SubscribeTickerSlim
//   - trades                      →  Client.SubscribeTrades
//
// Each select arm prints a one-line summary so it's obvious which channel
// fired. Cancel with Ctrl-C; the example exits cleanly via ctx-cancel.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	ctx, cancel := example.LongTimeout()
	defer cancel()
	c := example.MustWSPublic(ctx)
	defer c.Close()

	inst := example.Instrument()

	ob, err := c.SubscribeOrderBook(ctx, inst, "", 10)
	example.Fatal(err)
	defer ob.Close()

	tk, err := c.SubscribeTickerSlim(ctx, inst, "")
	example.Fatal(err)
	defer tk.Close()

	tr, err := c.SubscribeTrades(ctx, inst)
	example.Fatal(err)
	defer tr.Close()

	example.Print("multiplexing", inst)
	for {
		select {
		case <-ctx.Done():
			return
		case b, ok := <-ob.Updates():
			if !ok {
				return
			}
			example.Print("orderbook", len(b.Bids))
		case s, ok := <-tk.Updates():
			if !ok {
				return
			}
			example.Print("ticker", s.Ticker.MarkPrice)
		case ts, ok := <-tr.Updates():
			if !ok {
				return
			}
			example.Print("trades", len(ts))
		}
	}
}

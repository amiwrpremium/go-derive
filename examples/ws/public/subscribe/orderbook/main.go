// Streams orderbook updates for one instrument and demos every
// SubscribeOption: a larger buffer, DropOldest so the freshest book
// always wins under back-pressure, and an error handler that
// surfaces drops to a counter.
package main

import (
	"sync/atomic"

	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/ws"
)

func main() {
	ctx, cancel := example.LongTimeout()
	defer cancel()
	c := example.MustWSPublic(ctx)
	defer c.Close()

	var dropped atomic.Int64

	// Derive accepts depth in {1, 10}. Group in {1, 10, 100} (price-bucket size).
	//
	// Defaults if no opts are passed: buffer=256, policy=DropNewest,
	// no error handler. The opts here demonstrate every knob.
	sub, err := c.SubscribeOrderBook(ctx, example.Instrument(), "", 10,
		ws.WithBufferSize(1024),
		ws.WithDropPolicy(ws.DropOldest),
		ws.WithErrorHandler(func(error) { dropped.Add(1) }),
	)
	example.Fatal(err)
	defer sub.Close()

	for {
		select {
		case <-ctx.Done():
			example.Print("dropped events", dropped.Load())
			return
		case ob, ok := <-sub.Updates():
			if !ok {
				example.Print("dropped events", dropped.Load())
				return
			}
			example.Print(ob.InstrumentName, len(ob.Bids))
		}
	}
}

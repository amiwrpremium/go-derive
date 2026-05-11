// Lists interest payments for the configured subaccount over WebSocket.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	events, err := c.GetInterestHistory(ctx, types.InterestHistoryQuery{})
	example.Fatal(err)
	example.Print("count", len(events))
}

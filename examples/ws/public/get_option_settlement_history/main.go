// Lists platform-wide option settlements over WebSocket.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPublic(ctx)
	defer c.Close()

	settlements, _, err := c.GetPublicOptionSettlementHistory(ctx, types.OptionSettlementHistoryQuery{}, types.PageRequest{})
	example.Fatal(err)
	example.Print("count", len(settlements))
}

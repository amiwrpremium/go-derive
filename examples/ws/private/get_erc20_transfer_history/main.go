// Lists ERC20 transfers for the configured subaccount over WebSocket.
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

	transfers, err := c.GetERC20TransferHistory(ctx, types.ERC20TransferHistoryQuery{})
	example.Fatal(err)
	example.Print("count", len(transfers))
}

// Lists ERC20 transfers for the configured subaccount over WebSocket.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	transfers, err := c.GetERC20TransferHistory(ctx, nil)
	example.Fatal(err)
	example.Print("count", len(transfers))
}

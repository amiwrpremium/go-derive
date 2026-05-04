// Lists every subaccount the wallet owns over WebSocket.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	ids, err := c.GetSubaccounts(ctx)
	example.Fatal(err)
	example.Print("ids", ids)
}

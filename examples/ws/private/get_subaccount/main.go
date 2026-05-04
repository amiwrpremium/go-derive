// Fetches the configured subaccount snapshot over WebSocket.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	sa, err := c.GetSubaccount(ctx)
	example.Fatal(err)
	example.Print("equity", sa.SubaccountValue)
}

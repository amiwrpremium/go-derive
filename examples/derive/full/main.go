// Builds the facade with private credentials and uses both REST and WS.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	c := example.MustDerivePrivate()
	defer c.Close()

	ctx, cancel := example.Timeout()
	defer cancel()
	example.Fatal(c.WS.Connect(ctx))
	example.Fatal(c.WS.Login(ctx))

	sa, err := c.REST.GetSubaccount(ctx)
	example.Fatal(err)
	example.Print("subaccount", sa.SubaccountID)
	example.Print("ws connected", c.WS.IsConnected())
}

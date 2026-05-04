// Fetches wallet-level account information for the configured signer.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	raw, err := c.GetAccount(ctx)
	example.Fatal(err)
	example.Print("account bytes", len(raw))
}

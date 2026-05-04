// Fetches funding payments received / paid by the configured subaccount.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	raw, err := c.GetFundingHistory(ctx, nil)
	example.Fatal(err)
	example.Print("funding history bytes", len(raw))
}

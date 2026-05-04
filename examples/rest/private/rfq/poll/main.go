// Polls outstanding RFQs for the subaccount.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	rfqs, err := c.PollRFQs(ctx)
	example.Fatal(err)
	example.Print("count", len(rfqs))
}

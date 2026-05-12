// Lists interest payments received / paid by the configured subaccount.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	events, err := c.GetInterestHistory(ctx, types.InterestHistoryQuery{})
	example.Fatal(err)
	example.Print("count", len(events))
}

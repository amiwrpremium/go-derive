// Lists option settlements for the configured subaccount.
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

	settlements, err := c.GetOptionSettlementHistory(ctx, types.OptionSettlementHistoryQuery{})
	example.Fatal(err)
	example.Print("count", len(settlements))
}

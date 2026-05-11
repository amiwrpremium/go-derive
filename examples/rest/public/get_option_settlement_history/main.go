// Lists platform-wide option settlements. Public — no auth needed.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	settlements, page, err := c.GetPublicOptionSettlementHistory(ctx, types.OptionSettlementHistoryQuery{}, types.PageRequest{})
	example.Fatal(err)
	example.Print("count", len(settlements))
	example.Print("page count", page.Count)
}

// Paginates the configured subaccount's outstanding and historical
// RFQs.
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

	rfqs, page, err := c.GetRFQs(ctx, types.RFQsQuery{}, types.PageRequest{PageSize: 10})
	example.Fatal(err)
	example.Print("count", len(rfqs))
	example.Print("total pages", page.NumPages)
}

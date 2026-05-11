// Lists ERC20 transfers for the configured subaccount.
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

	transfers, err := c.GetERC20TransferHistory(ctx, types.ERC20TransferHistoryQuery{})
	example.Fatal(err)
	example.Print("count", len(transfers))
	if len(transfers) > 0 {
		t := transfers[0]
		example.Print("first asset", t.Asset)
		example.Print("first amount", t.Amount.String())
	}
}

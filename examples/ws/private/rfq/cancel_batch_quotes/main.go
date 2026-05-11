// Bulk-cancels every quote on the configured subaccount over WebSocket.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	res, err := c.CancelBatchQuotes(ctx, types.CancelBatchInput{})
	example.Fatal(err)
	example.Print("cancelled", len(res.CancelledIDs))
}

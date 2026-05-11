// Bulk-cancels every outstanding RFQ on the configured subaccount
// over WebSocket.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	res, err := c.CancelBatchRFQs(ctx, nil)
	example.Fatal(err)
	example.Print("cancelled", len(res.CancelledIDs))
}

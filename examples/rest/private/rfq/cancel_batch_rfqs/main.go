// Bulk-cancels every outstanding RFQ on the configured subaccount.
// Returns the list of cancelled rfq ids.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	res, err := c.CancelBatchRFQs(ctx, nil)
	example.Fatal(err)
	example.Print("cancelled", len(res.CancelledIDs))
	for i, id := range res.CancelledIDs {
		if i >= 5 {
			break
		}
		example.Print("  rfq_id", id)
	}
}

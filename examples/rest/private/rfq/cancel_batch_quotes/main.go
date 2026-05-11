// Bulk-cancels every quote on the configured subaccount, optionally
// filtered by label / nonce. Returns the list of cancelled quote ids.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	res, err := c.CancelBatchQuotes(ctx, nil)
	example.Fatal(err)
	example.Print("cancelled", len(res.CancelledIDs))
	for i, id := range res.CancelledIDs {
		if i >= 5 {
			break
		}
		example.Print("  quote_id", id)
	}
}

// Bulk-updates notification statuses (e.g. mark every unseen as
// seen). The default no-filter call marks every notification on
// the configured subaccount as seen.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	res, err := c.UpdateNotifications(ctx, map[string]any{
		"status": "seen",
	})
	example.Fatal(err)
	example.Print("updated", res.UpdatedCount)
}

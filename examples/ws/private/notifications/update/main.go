// Bulk-updates notification statuses over WebSocket.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	res, err := c.UpdateNotifications(ctx, map[string]any{
		"status": "seen",
	})
	example.Fatal(err)
	example.Print("updated", res.UpdatedCount)
}

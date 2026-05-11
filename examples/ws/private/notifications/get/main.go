// Paginates notifications over WebSocket.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	events, _, err := c.GetNotifications(ctx, map[string]any{"page_size": 10})
	example.Fatal(err)
	example.Print("count", len(events))
}

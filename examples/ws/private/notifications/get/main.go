// Paginates notifications over WebSocket.
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

	events, _, err := c.GetNotifications(ctx, types.NotificationsQuery{}, types.PageRequest{PageSize: 10})
	example.Fatal(err)
	example.Print("count", len(events))
}

// Marks specific notifications as seen or hidden over WebSocket.
// Required env: DERIVE_NOTIFICATION_ID.
package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	raw := os.Getenv("DERIVE_NOTIFICATION_ID")
	if raw == "" {
		log.Fatal("DERIVE_NOTIFICATION_ID required (comma-separated ids)")
	}
	ids := []int64{}
	for _, tok := range strings.Split(raw, ",") {
		id, err := strconv.ParseInt(strings.TrimSpace(tok), 10, 64)
		example.Fatal(err)
		ids = append(ids, id)
	}
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	res, err := c.UpdateNotifications(ctx, types.UpdateNotificationsInput{
		NotificationIDs: ids,
		Status:          "seen",
	})
	example.Fatal(err)
	example.Print("updated", res.UpdatedCount)
}

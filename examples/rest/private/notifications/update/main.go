// Marks specific notifications as seen or hidden. Required env:
// DERIVE_NOTIFICATION_ID (comma-separated list of notification ids).
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
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	res, err := c.UpdateNotifications(ctx, types.UpdateNotificationsInput{
		NotificationIDs: ids,
		Status:          "seen",
	})
	example.Fatal(err)
	example.Print("updated", res.UpdatedCount)
}

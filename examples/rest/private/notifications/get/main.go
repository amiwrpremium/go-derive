// Paginates notifications for the configured subaccount.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	events, page, err := c.GetNotifications(ctx, map[string]any{"page_size": 10})
	example.Fatal(err)
	example.Print("count", len(events))
	example.Print("total pages", page.NumPages)
	for i, ev := range events {
		if i >= 3 {
			break
		}
		example.Print(ev.Event, ev.Status)
	}
}

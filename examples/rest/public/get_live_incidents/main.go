// Lists ongoing platform incidents. Empty list means no active incidents.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	incidents, err := c.GetLiveIncidents(ctx)
	example.Fatal(err)
	example.Print("active incidents", len(incidents))
	for _, in := range incidents {
		example.Print("incident", in.Label)
		example.Print("  severity", in.Severity)
		example.Print("  monitor_type", in.MonitorType)
		example.Print("  message", in.Message)
	}
}

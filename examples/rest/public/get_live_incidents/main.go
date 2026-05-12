// Lists ongoing platform incidents. Empty list means no active incidents.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amiwrpremium/go-derive/pkg/rest"
)

func main() {
	restNetwork := rest.WithTestnet()
	if os.Getenv("DERIVE_NETWORK") == "mainnet" {
		restNetwork = rest.WithMainnet()
	}
	c, err := rest.New(restNetwork)
	if err != nil {
		log.Fatalf("rest.New: %v", err)
	}
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	incidents, err := c.GetLiveIncidents(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "active incidents:", len(incidents))
	for _, in := range incidents {
		fmt.Printf("%-30s %v\n", "incident:", in.Label)
		fmt.Printf("%-30s %v\n", "  severity:", in.Severity)
		fmt.Printf("%-30s %v\n", "  monitor_type:", in.MonitorType)
		fmt.Printf("%-30s %v\n", "  message:", in.Message)
	}
}

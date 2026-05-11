// Reads MMP config over WebSocket. Optional DERIVE_CURRENCY filter.
package main

import (
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	currency := os.Getenv("DERIVE_CURRENCY")
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	configs, err := c.GetMMPConfig(ctx, currency)
	example.Fatal(err)
	example.Print("rules", len(configs))
}

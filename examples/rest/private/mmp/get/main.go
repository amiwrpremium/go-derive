// Reads market-maker-protection config for the configured subaccount,
// optionally filtered to one currency (empty string returns every
// configured currency).
package main

import (
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	currency := os.Getenv("DERIVE_CURRENCY") // empty → all currencies
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	configs, err := c.GetMMPConfig(ctx, currency)
	example.Fatal(err)
	example.Print("rules", len(configs))
	for _, cfg := range configs {
		example.Print(cfg.Currency, cfg.IsFrozen)
	}
}

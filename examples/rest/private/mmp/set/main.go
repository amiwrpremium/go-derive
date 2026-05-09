// Configures market-maker protection for one currency.
package main

import (
	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	example.Fatal(c.SetMMPConfig(ctx, derive.MMPConfig{
		Currency:        "BTC",
		MMPFrozenTimeMs: 5000,
		MMPIntervalMs:   1000,
	}))
	example.Print("set", "ok")
}

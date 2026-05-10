// Configures market-maker protection for one currency.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	example.Fatal(c.SetMMPConfig(ctx, types.MMPConfig{
		Currency:        "BTC",
		MMPFrozenTimeMs: 5000,
		MMPIntervalMs:   1000,
	}))
	example.Print("set", "ok")
}

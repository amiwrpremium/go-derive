// Sets MMP config over WebSocket.
package main

import (
	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	example.Fatal(c.SetMMPConfig(ctx, derive.MMPConfig{
		Currency:        "BTC",
		MMPFrozenTimeMs: 5000,
		MMPIntervalMs:   1000,
	}))
	example.Print("set", "ok")
}

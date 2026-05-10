// Sets MMP config over WebSocket.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	example.Fatal(c.SetMMPConfig(ctx, types.MMPConfig{
		Currency:        "BTC",
		MMPFrozenTimeMs: 5000,
		MMPIntervalMs:   1000,
	}))
	example.Print("set", "ok")
}

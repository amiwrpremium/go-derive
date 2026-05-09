// Arms (or disarms) the cancel-on-disconnect kill-switch for the wallet.
//
// When armed, every open order on this wallet is cancelled if the
// authenticated WebSocket session drops. Critical for makers who need to
// avoid stale quotes during a transient network blip.
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	example.Fatal(c.SetCancelOnDisconnect(ctx, true))
	example.Print("set_cancel_on_disconnect", "armed")
}

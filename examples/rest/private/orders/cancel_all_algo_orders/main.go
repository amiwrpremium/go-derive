// Cancels every in-flight algo order on the configured subaccount.
// Set DERIVE_RUN_LIVE_ORDERS=1 to actually run; otherwise this exits.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	if os.Getenv("DERIVE_RUN_LIVE_ORDERS") != "1" {
		log.Fatal("set DERIVE_RUN_LIVE_ORDERS=1 to actually cancel every algo order")
	}
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	example.Fatal(c.CancelAllAlgoOrders(ctx))
	example.Print("status", "ok")
}

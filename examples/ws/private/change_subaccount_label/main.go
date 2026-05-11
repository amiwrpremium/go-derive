// Sets a human-readable label on the configured subaccount over
// WebSocket. Required env: DERIVE_LABEL.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	label := os.Getenv("DERIVE_LABEL")
	if label == "" {
		log.Fatal("DERIVE_LABEL required")
	}
	ctx, cancel := example.Timeout()
	defer cancel()
	c := example.MustWSPrivate(ctx)
	defer c.Close()

	example.Fatal(c.ChangeSubaccountLabel(ctx, label))
	example.Print("label set", label)
}

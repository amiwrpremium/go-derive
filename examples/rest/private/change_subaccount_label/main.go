// Sets a human-readable label on the configured subaccount.
// Required env: DERIVE_LABEL.
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
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	example.Fatal(c.ChangeSubaccountLabel(ctx, label))
	example.Print("label set", label)
}

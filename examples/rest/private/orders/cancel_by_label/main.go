// Cancels every order carrying a given label. Set DERIVE_LABEL.
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

	n, err := c.CancelByLabel(ctx, label)
	example.Fatal(err)
	example.Print("cancelled", n)
}

// One wallet's points record for one program. Required env:
// DERIVE_PROGRAM_NAME and DERIVE_WALLET.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	name := os.Getenv("DERIVE_PROGRAM_NAME")
	wallet := os.Getenv("DERIVE_WALLET")
	if name == "" || wallet == "" {
		log.Fatal("DERIVE_PROGRAM_NAME and DERIVE_WALLET required")
	}
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	res, err := c.GetPoints(ctx, name, wallet)
	example.Fatal(err)
	example.Print("total_notional_volume", res.TotalNotionalVolume.String())
	example.Print("points (raw bytes)", len(res.Points))
}

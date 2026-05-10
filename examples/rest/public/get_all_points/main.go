// Program-wide points snapshot. Required env: DERIVE_PROGRAM_NAME.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	name := os.Getenv("DERIVE_PROGRAM_NAME")
	if name == "" {
		log.Fatal("DERIVE_PROGRAM_NAME required")
	}
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	res, err := c.GetAllPoints(ctx, name)
	example.Fatal(err)
	example.Print("total_users", res.TotalUsers)
	example.Print("total_notional_volume", res.TotalNotionalVolume.String())
	example.Print("points (raw bytes)", len(res.Points))
}

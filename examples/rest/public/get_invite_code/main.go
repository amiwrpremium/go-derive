// Returns the invite code allocated to one wallet plus its
// remaining-uses counter. Optional env: DERIVE_WALLET.
package main

import (
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	res, err := c.GetInviteCode(ctx, os.Getenv("DERIVE_WALLET"))
	example.Fatal(err)
	example.Print("code", res.Code)
	example.Print("remaining_uses (-1=unlimited)", res.RemainingUses)
}

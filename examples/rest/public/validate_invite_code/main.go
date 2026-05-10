// Validates one invite code. Required env: DERIVE_INVITE_CODE.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	code := os.Getenv("DERIVE_INVITE_CODE")
	if code == "" {
		log.Fatal("DERIVE_INVITE_CODE required")
	}
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	status, err := c.ValidateInviteCode(ctx, code)
	example.Fatal(err)
	example.Print("status", status)
}

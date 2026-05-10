// Returns the referral code currently associated with one wallet.
// Optional env: DERIVE_WALLET (defaults to signer's wallet if any).
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

	code, err := c.GetReferralCode(ctx, os.Getenv("DERIVE_WALLET"))
	example.Fatal(err)
	example.Print("referral_code", code)
}

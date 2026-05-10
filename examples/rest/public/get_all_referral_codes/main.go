// Lists every valid referral code for the configured signer's
// wallet (or omits the wallet param when no signer is configured).
package main

import "github.com/amiwrpremium/go-derive/examples/example"

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	codes, err := c.GetAllReferralCodes(ctx)
	example.Fatal(err)
	example.Print("codes", len(codes))
	for i, r := range codes {
		if i >= 3 {
			break
		}
		example.Print("code", r.ReferralCode)
		example.Print("  wallet", r.Wallet)
	}
}

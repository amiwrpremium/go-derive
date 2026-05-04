// Signs an EIP-712 ActionData and prints the resulting signature.
package main

import (
	"github.com/amiwrpremium/go-derive/examples/example"
	"github.com/amiwrpremium/go-derive/pkg/auth"
)

func main() {
	s := example.MustSigner()
	domain := example.Network().EIP712Domain()

	action := auth.ActionData{
		SubaccountID: example.Subaccount(),
		Nonce:        auth.NewNonceGen().Next(),
		Expiry:       1_700_000_000,
	}
	ctx, cancel := example.Timeout()
	defer cancel()
	sig, err := s.SignAction(ctx, domain, action)
	example.Fatal(err)
	example.Print("signature", sig.Hex())
}

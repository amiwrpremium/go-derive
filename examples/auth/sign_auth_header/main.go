// Signs an EIP-191 timestamp for use as the X-LyraSignature header.
package main

import (
	"time"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	s := example.MustSigner()
	ctx, cancel := example.Timeout()
	defer cancel()

	sig, err := s.SignAuthHeader(ctx, time.Now())
	example.Fatal(err)
	example.Print("signature", sig.Hex())
}

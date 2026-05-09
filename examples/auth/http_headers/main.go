// Builds the X-LyraWallet/Timestamp/Signature header bundle.
package main

import (
	"time"

	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	s := example.MustSigner()
	ctx, cancel := example.Timeout()
	defer cancel()

	h, err := derive.HTTPHeaders(ctx, s, time.Now())
	example.Fatal(err)
	example.Print("X-LyraWallet", h.Get("X-LyraWallet"))
	example.Print("X-LyraTimestamp", h.Get("X-LyraTimestamp"))
	example.Print("X-LyraSignature", h.Get("X-LyraSignature"))
}

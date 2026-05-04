// Package goderive is the root of the go-derive SDK for the Derive exchange
// (formerly Lyra). The user-facing API lives under pkg/. Most users want
// pkg/derive for the top-level facade, or pkg/rest and pkg/ws directly for
// fine-grained control.
//
// Quick start:
//
//	import (
//	    "context"
//	    "github.com/amiwrpremium/go-derive/pkg/auth"
//	    "github.com/amiwrpremium/go-derive/pkg/derive"
//	)
//
//	signer, _ := auth.NewLocalSigner("0xPRIVATEKEY")
//	c, _ := derive.NewClient(
//	    derive.WithMainnet(),
//	    derive.WithSigner(signer),
//	    derive.WithSubaccount(123),
//	)
//	defer c.Close()
//
//	instruments, err := c.REST.GetInstruments(context.Background(), "BTC", "perp")
package goderive

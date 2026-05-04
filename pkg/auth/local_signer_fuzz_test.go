package auth_test

import (
	"testing"

	"github.com/amiwrpremium/go-derive/pkg/auth"
)

// FuzzNewLocalSigner verifies that bad hex input never panics. Real keys
// produce a signer; everything else returns an error.
func FuzzNewLocalSigner(f *testing.F) {
	f.Add("0x4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318")
	f.Add("4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318")
	f.Add("0xZZZZ")
	f.Add("")
	f.Add("0x")
	f.Add(string(make([]byte, 256)))

	f.Fuzz(func(t *testing.T, s string) {
		// Either it parses and gives a non-nil signer with a valid address,
		// or it errors out — anything else is a defect.
		signer, err := auth.NewLocalSigner(s)
		if err != nil {
			if signer != nil {
				t.Fatalf("error path returned non-nil signer for %q", s)
			}
			return
		}
		if signer == nil {
			t.Fatalf("nil signer with no error for %q", s)
		}
	})
}

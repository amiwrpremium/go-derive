package types_test

import (
	"testing"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

// FuzzNewAddress verifies the parser is panic-free for arbitrary input.
func FuzzNewAddress(f *testing.F) {
	f.Add("0x1111111111111111111111111111111111111111")
	f.Add("0xZZZZ")
	f.Add("not-an-address")
	f.Add("")
	f.Add("0x")
	f.Add(string(make([]byte, 1024)))

	f.Fuzz(func(t *testing.T, s string) {
		a, err := types.NewAddress(s)
		if err != nil && !a.IsZero() {
			t.Fatalf("error path leaked non-zero address for %q", s)
		}
	})
}

// FuzzNewTxHash verifies the parser is panic-free for arbitrary input.
func FuzzNewTxHash(f *testing.F) {
	f.Add("0x1111111111111111111111111111111111111111111111111111111111111111")
	f.Add("0xabc")
	f.Add("0x")
	f.Add("")
	f.Add("not-a-hash")

	f.Fuzz(func(t *testing.T, s string) {
		h, err := types.NewTxHash(s)
		if err != nil && !h.IsZero() {
			t.Fatalf("error path leaked non-zero hash for %q", s)
		}
	})
}

package types_test

import (
	"encoding/json"
	"testing"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

// FuzzNewDecimal verifies the parser never panics on arbitrary input. Bad
// input must produce a (zero, error); good input must round-trip.
func FuzzNewDecimal(f *testing.F) {
	f.Add("0")
	f.Add("1.5")
	f.Add("-2.5")
	f.Add("0.000000000000000001")
	f.Add("100000000000000000000")
	f.Add("1e10")
	f.Add("not-a-number")
	f.Add("")
	f.Add("0x10")
	f.Add(string([]byte{0xff, 0xfe}))

	f.Fuzz(func(t *testing.T, s string) {
		d, err := types.NewDecimal(s)
		if err != nil {
			// Error path: receiver must be the zero value.
			if !d.IsZero() {
				t.Fatalf("error path returned non-zero decimal: %q -> %s", s, d)
			}
			return
		}
		// Success path: round-trip through JSON and back.
		b, err := json.Marshal(d)
		if err != nil {
			t.Fatalf("marshal succeeded-parse decimal: %v", err)
		}
		var back types.Decimal
		if err := json.Unmarshal(b, &back); err != nil {
			t.Fatalf("round-trip unmarshal: %v", err)
		}
	})
}

// FuzzDecimal_UnmarshalJSON checks the JSON unmarshaler doesn't panic.
func FuzzDecimal_UnmarshalJSON(f *testing.F) {
	f.Add([]byte(`"1.5"`))
	f.Add([]byte(`1.5`))
	f.Add([]byte(`null`))
	f.Add([]byte(`""`))
	f.Add([]byte(`"-0.001"`))
	f.Add([]byte(`{`))
	f.Add([]byte(``))
	f.Add([]byte{0xff, 0xfe})

	f.Fuzz(func(t *testing.T, raw []byte) {
		var d types.Decimal
		_ = d.UnmarshalJSON(raw) // panic = failure; error is fine
	})
}

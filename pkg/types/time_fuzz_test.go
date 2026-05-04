package types_test

import (
	"testing"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

// FuzzMillisTime_UnmarshalJSON ensures the time unmarshaler is panic-free
// on arbitrary input, including non-string-non-number JSON.
func FuzzMillisTime_UnmarshalJSON(f *testing.F) {
	f.Add([]byte(`1700000000000`))
	f.Add([]byte(`"1700000000000"`))
	f.Add([]byte(`null`))
	f.Add([]byte(`""`))
	f.Add([]byte(`"abc"`))
	f.Add([]byte(`-1`))
	f.Add([]byte(`{}`))

	f.Fuzz(func(t *testing.T, raw []byte) {
		var mt types.MillisTime
		_ = mt.UnmarshalJSON(raw)
	})
}

package types_test

import (
	"testing"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

// FuzzOrderBookLevel_UnmarshalJSON guards the array-shaped level decoder
// against panics on adversarial input.
func FuzzOrderBookLevel_UnmarshalJSON(f *testing.F) {
	f.Add([]byte(`["100","1"]`))
	f.Add([]byte(`["",""]`))
	f.Add([]byte(`[]`))
	f.Add([]byte(`[1,2,3]`))
	f.Add([]byte(`{"price":"1","amount":"2"}`))
	f.Add([]byte(`null`))
	f.Add([]byte(``))

	f.Fuzz(func(t *testing.T, raw []byte) {
		var l types.OrderBookLevel
		_ = l.UnmarshalJSON(raw)
	})
}

package errors_test

import (
	"encoding/json"
	"testing"

	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
)

// FuzzAPIError_UnmarshalJSON guards the on-the-wire error decoder against
// panics on adversarial input.
func FuzzAPIError_UnmarshalJSON(f *testing.F) {
	f.Add([]byte(`{"code":-32000,"message":"throttled"}`))
	f.Add([]byte(`{"code":11015,"message":"frozen","data":{"why":"none"}}`))
	f.Add([]byte(`{"code":"not-a-number"}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`null`))
	f.Add([]byte(``))

	f.Fuzz(func(t *testing.T, raw []byte) {
		var e derrors.APIError
		_ = json.Unmarshal(raw, &e)
		_ = e.Error()            // formatting must not panic on any state
		_ = e.CanonicalMessage() // lookup must not panic on any state
	})
}

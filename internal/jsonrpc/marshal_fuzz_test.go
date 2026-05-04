package jsonrpc_test

import (
	"encoding/json"
	"testing"

	"github.com/amiwrpremium/go-derive/internal/jsonrpc"
)

// FuzzDecodeResult drives a varied [Response] payload through the result
// decoder and checks for panics. Decoding failures are expected and fine.
func FuzzDecodeResult(f *testing.F) {
	f.Add([]byte(`{}`))
	f.Add([]byte(`{"id":1,"result":1}`))
	f.Add([]byte(`{"id":1,"error":{"code":-32600,"message":"x"}}`))
	f.Add([]byte(`{"id":1,"result":null}`))
	f.Add([]byte(``))
	f.Add([]byte(`{`))

	f.Fuzz(func(t *testing.T, raw []byte) {
		var resp jsonrpc.Response
		if err := json.Unmarshal(raw, &resp); err != nil {
			return // malformed envelope; skip
		}
		var out any
		_ = jsonrpc.DecodeResult(&resp, &out)
	})
}

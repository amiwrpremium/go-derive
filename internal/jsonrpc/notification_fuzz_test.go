package jsonrpc_test

import (
	"testing"

	"github.com/amiwrpremium/go-derive/internal/jsonrpc"
)

// FuzzIsNotification ensures the frame classifier never panics. It also
// runs on hot-path bytes received from the network — robustness here is
// load-bearing.
func FuzzIsNotification(f *testing.F) {
	f.Add([]byte(`{"jsonrpc":"2.0","method":"subscription","params":{}}`))
	f.Add([]byte(`{"jsonrpc":"2.0","id":1,"result":{}}`))
	f.Add([]byte(`{"id":null,"method":"x"}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(``))
	f.Add([]byte(`not json`))
	f.Add([]byte{0x00, 0xff})

	f.Fuzz(func(t *testing.T, raw []byte) {
		_ = jsonrpc.IsNotification(raw)
	})
}

package jsonrpc_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/internal/jsonrpc"
)

func TestNewRequest_NilParams_OmitsField(t *testing.T) {
	r, err := jsonrpc.NewRequest(7, "public/get_instruments", nil)
	require.NoError(t, err)
	b, err := json.Marshal(r)
	require.NoError(t, err)
	assert.JSONEq(t, `{"jsonrpc":"2.0","id":7,"method":"public/get_instruments"}`, string(b))
}

func TestNewRequest_TypedParams_Marshalled(t *testing.T) {
	r, err := jsonrpc.NewRequest(1, "private/order", map[string]string{"side": "buy"})
	require.NoError(t, err)
	b, err := json.Marshal(r)
	require.NoError(t, err)
	assert.JSONEq(t, `{"jsonrpc":"2.0","id":1,"method":"private/order","params":{"side":"buy"}}`, string(b))
}

func TestNewRequest_MarshalFailure(t *testing.T) {
	_, err := jsonrpc.NewRequest(1, "x", make(chan int))
	assert.Error(t, err, "channels can't be marshalled to JSON")
}

func TestError_Error_WithoutData(t *testing.T) {
	e := &jsonrpc.Error{Code: -32600, Message: "Invalid Request"}
	got := e.Error()
	assert.Contains(t, got, "-32600")
	assert.Contains(t, got, "Invalid Request")
	assert.NotContains(t, got, "(")
}

func TestError_Error_WithData(t *testing.T) {
	e := &jsonrpc.Error{
		Code: -32602, Message: "bad",
		Data: json.RawMessage(`{"why":"none"}`),
	}
	got := e.Error()
	assert.Contains(t, got, `{"why":"none"}`)
	assert.Contains(t, got, "(")
}

func TestVersionConstant(t *testing.T) {
	assert.Equal(t, "2.0", jsonrpc.Version)
}

func TestResponse_IDUint64(t *testing.T) {
	cases := []struct {
		name   string
		raw    string
		wantN  uint64
		wantOK bool
	}{
		{"numeric", `{"jsonrpc":"2.0","id":42,"result":null}`, 42, true},
		{"missing", `{"jsonrpc":"2.0","result":null}`, 0, false},
		{"string", `{"jsonrpc":"2.0","id":"abc","result":null}`, 0, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var r jsonrpc.Response
			require.NoError(t, json.Unmarshal([]byte(c.raw), &r))
			n, ok := r.IDUint64()
			assert.Equal(t, c.wantN, n)
			assert.Equal(t, c.wantOK, ok)
		})
	}
}

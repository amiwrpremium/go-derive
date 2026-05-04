package jsonrpc_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/internal/jsonrpc"
)

func TestDecodeResult_NilResponse(t *testing.T) {
	err := jsonrpc.DecodeResult(nil, nil)
	assert.Error(t, err)
}

func TestDecodeResult_ResponseWithError(t *testing.T) {
	resp := &jsonrpc.Response{
		Error: &jsonrpc.Error{Code: -32600, Message: "Invalid Request"},
	}
	var out struct{}
	err := jsonrpc.DecodeResult(resp, &out)
	assert.ErrorContains(t, err, "Invalid Request")
}

func TestDecodeResult_NilOut_SkipsDecode(t *testing.T) {
	resp := &jsonrpc.Response{Result: json.RawMessage(`{}`)}
	assert.NoError(t, jsonrpc.DecodeResult(resp, nil))
}

func TestDecodeResult_EmptyResult_SkipsDecode(t *testing.T) {
	resp := &jsonrpc.Response{}
	var out struct{ X int }
	assert.NoError(t, jsonrpc.DecodeResult(resp, &out))
}

func TestDecodeResult_DecodesIntoOut(t *testing.T) {
	resp := &jsonrpc.Response{
		Result: json.RawMessage(`{"name":"BTC-PERP"}`),
	}
	var out struct {
		Name string `json:"name"`
	}
	require.NoError(t, jsonrpc.DecodeResult(resp, &out))
	assert.Equal(t, "BTC-PERP", out.Name)
}

func TestDecodeResult_DecodeFailure(t *testing.T) {
	resp := &jsonrpc.Response{Result: json.RawMessage(`not-json`)}
	var out struct{ X int }
	assert.Error(t, jsonrpc.DecodeResult(resp, &out))
}

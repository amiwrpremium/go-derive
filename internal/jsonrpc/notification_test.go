package jsonrpc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/internal/jsonrpc"
)

func TestIsNotification_TrueForServerInitiated(t *testing.T) {
	frame := []byte(`{"jsonrpc":"2.0","method":"subscription","params":{"channel":"trades.BTC-PERP"}}`)
	assert.True(t, jsonrpc.IsNotification(frame))
}

func TestIsNotification_FalseForResponse(t *testing.T) {
	frame := []byte(`{"jsonrpc":"2.0","id":1,"result":{}}`)
	assert.False(t, jsonrpc.IsNotification(frame))
}

func TestIsNotification_FalseForRequest(t *testing.T) {
	frame := []byte(`{"jsonrpc":"2.0","id":1,"method":"public/get_time"}`)
	assert.False(t, jsonrpc.IsNotification(frame))
}

func TestIsNotification_FalseForGarbage(t *testing.T) {
	assert.False(t, jsonrpc.IsNotification([]byte(`not even json`)))
}

func TestIsNotification_FalseForEmptyObject(t *testing.T) {
	assert.False(t, jsonrpc.IsNotification([]byte(`{}`)))
}

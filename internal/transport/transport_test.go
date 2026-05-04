package transport_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/internal/transport"
)

// Compile-time interface assertions are the primary contract here.
// HTTPTransport implements Transport; WSTransport implements both
// Transport and Subscriber.
var (
	_ transport.Transport  = (*transport.HTTPTransport)(nil)
	_ transport.Transport  = (*transport.WSTransport)(nil)
	_ transport.Subscriber = (*transport.WSTransport)(nil)
)

func TestDecoder_TypeAlias(t *testing.T) {
	// Decoder is a function alias; verify it can be assigned and called.
	var d transport.Decoder = func(raw json.RawMessage) (any, error) {
		return string(raw), nil
	}
	v, err := d(json.RawMessage(`"hi"`))
	assert.NoError(t, err)
	assert.Equal(t, `"hi"`, v.(string))
}

// stubTransport is the smallest possible Transport — it lets us prove the
// interface is consumable from outside this package.
type stubTransport struct{ called bool }

func (s *stubTransport) Call(context.Context, string, any, any) error {
	s.called = true
	return nil
}

func (s *stubTransport) Close() error { return nil }

func TestTransport_StubImplements(t *testing.T) {
	var tr transport.Transport = &stubTransport{}
	require := assert.New(t)
	require.NoError(tr.Call(context.Background(), "x", nil, nil))
	require.True(tr.(*stubTransport).called)
	require.NoError(tr.Close())
}

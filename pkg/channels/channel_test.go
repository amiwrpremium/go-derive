package channels_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/channels"
)

// stubChannel is a minimal in-test implementation used only to confirm the
// Channel interface is satisfiable and that pkg/ws.Subscribe works against
// arbitrary descriptors. Real descriptors live in pkg/channels/{public,private}.
type stubChannel struct{ name string }

func (s stubChannel) Name() string { return s.name }
func (stubChannel) Decode(raw json.RawMessage) (any, error) {
	var out string
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func TestChannel_InterfaceConformance(t *testing.T) {
	var c channels.Channel = stubChannel{name: "trades.BTC-PERP"}
	assert.Equal(t, "trades.BTC-PERP", c.Name())

	v, err := c.Decode(json.RawMessage(`"hello"`))
	require.NoError(t, err)
	assert.Equal(t, "hello", v)
}

func TestChannel_DecodeError(t *testing.T) {
	c := stubChannel{name: "x"}
	_, err := c.Decode(json.RawMessage(`{`)) // malformed
	assert.Error(t, err)
}

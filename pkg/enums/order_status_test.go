package enums_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/enums"
)

func TestOrderStatus_Valid_AllArms(t *testing.T) {
	cases := []enums.OrderStatus{
		enums.OrderStatusOpen,
		enums.OrderStatusFilled,
		enums.OrderStatusCancelled,
		enums.OrderStatusExpired,
		enums.OrderStatusRejected,
	}
	for _, c := range cases {
		t.Run(string(c), func(t *testing.T) {
			assert.True(t, c.Valid())
		})
	}
}

func TestOrderStatus_Valid_RejectsUnknown(t *testing.T) {
	assert.False(t, enums.OrderStatus("").Valid())
	assert.False(t, enums.OrderStatus("haunted").Valid())
	// Previously hallucinated values that the canonical schema does not include:
	assert.False(t, enums.OrderStatus("untriggered").Valid())
	assert.False(t, enums.OrderStatus("insufficient_margin").Valid())
}

func TestOrderStatus_Terminal_TerminalArms(t *testing.T) {
	terminal := []enums.OrderStatus{
		enums.OrderStatusFilled,
		enums.OrderStatusCancelled,
		enums.OrderStatusExpired,
		enums.OrderStatusRejected,
	}
	for _, s := range terminal {
		t.Run(string(s), func(t *testing.T) {
			assert.True(t, s.Terminal())
		})
	}
}

func TestOrderStatus_Terminal_OpenArm(t *testing.T) {
	assert.False(t, enums.OrderStatusOpen.Terminal())
}

func TestOrderStatus_Terminal_DefaultArm(t *testing.T) {
	// Unknown values fall through the switch's default and report non-terminal.
	assert.False(t, enums.OrderStatus("???").Terminal())
}

func TestOrderStatus_JSONRoundTrip(t *testing.T) {
	type wrap struct {
		S enums.OrderStatus `json:"s"`
	}
	in := wrap{S: enums.OrderStatusFilled}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	assert.JSONEq(t, `{"s":"filled"}`, string(b))

	var out wrap
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in, out)
}

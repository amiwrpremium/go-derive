package enums_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/enums"
)

func TestOrderType_Valid_Limit(t *testing.T) {
	assert.True(t, enums.OrderTypeLimit.Valid())
}

func TestOrderType_Valid_Market(t *testing.T) {
	assert.True(t, enums.OrderTypeMarket.Valid())
}

func TestOrderType_Valid_RejectsUnknown(t *testing.T) {
	for _, v := range []string{"", "stop", "LIMIT"} {
		assert.False(t, enums.OrderType(v).Valid(), "value %q", v)
	}
}

func TestOrderType_JSONRoundTrip(t *testing.T) {
	type wrap struct {
		T enums.OrderType `json:"t"`
	}
	in := wrap{T: enums.OrderTypeMarket}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	assert.JSONEq(t, `{"t":"market"}`, string(b))

	var out wrap
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in, out)
}

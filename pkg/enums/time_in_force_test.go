package enums_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/enums"
)

func TestTimeInForce_Valid_AllArms(t *testing.T) {
	cases := []enums.TimeInForce{
		enums.TimeInForceGTC,
		enums.TimeInForcePostOnly,
		enums.TimeInForceFOK,
		enums.TimeInForceIOC,
	}
	for _, c := range cases {
		t.Run(string(c), func(t *testing.T) {
			assert.True(t, c.Valid())
		})
	}
}

func TestTimeInForce_Valid_RejectsUnknown(t *testing.T) {
	for _, v := range []string{"", "FOK", "forever", "day"} {
		assert.False(t, enums.TimeInForce(v).Valid(), "value %q", v)
	}
}

func TestTimeInForce_JSONRoundTrip(t *testing.T) {
	type wrap struct {
		T enums.TimeInForce `json:"t"`
	}
	in := wrap{T: enums.TimeInForcePostOnly}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	assert.JSONEq(t, `{"t":"post_only"}`, string(b))

	var out wrap
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in, out)
}

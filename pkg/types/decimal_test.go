package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestNewDecimal_Valid(t *testing.T) {
	d, err := types.NewDecimal("1.234")
	require.NoError(t, err)
	assert.Equal(t, "1.234", d.String())
}

func TestNewDecimal_Invalid(t *testing.T) {
	_, err := types.NewDecimal("not-a-number")
	assert.Error(t, err)
}

func TestMustDecimal_Panics(t *testing.T) {
	assert.Panics(t, func() { types.MustDecimal("not-a-number") })
	assert.NotPanics(t, func() { types.MustDecimal("0") })
}

func TestDecimalFromInt(t *testing.T) {
	d := types.DecimalFromInt(42)
	assert.Equal(t, "42", d.String())
	assert.False(t, d.IsZero())
	assert.Equal(t, 1, d.Sign())

	z := types.DecimalFromInt(0)
	assert.True(t, z.IsZero())
	assert.Equal(t, 0, z.Sign())

	n := types.DecimalFromInt(-5)
	assert.Equal(t, -1, n.Sign())
}

func TestDecimal_JSONRoundTrip(t *testing.T) {
	type wrap struct {
		Price types.Decimal `json:"price"`
	}
	in := wrap{Price: types.MustDecimal("65000.5")}

	b, err := json.Marshal(in)
	require.NoError(t, err)
	assert.JSONEq(t, `{"price":"65000.5"}`, string(b))

	var out wrap
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.Price.String(), out.Price.String())
}

func TestDecimal_UnmarshalAcceptsNumber(t *testing.T) {
	var d types.Decimal
	require.NoError(t, json.Unmarshal([]byte(`123.45`), &d))
	assert.Equal(t, "123.45", d.String())
}

func TestDecimal_UnmarshalNullAndEmpty(t *testing.T) {
	var d types.Decimal
	require.NoError(t, json.Unmarshal([]byte(`null`), &d))
	assert.True(t, d.IsZero())

	d2 := types.MustDecimal("1")
	require.NoError(t, json.Unmarshal([]byte(`""`), &d2))
	// Empty string left previous value alone (we treat as no-op).
	assert.Equal(t, "1", d2.String())
}

func TestDecimal_UnmarshalMalformedString(t *testing.T) {
	var d types.Decimal
	err := json.Unmarshal([]byte(`"abc"`), &d)
	assert.Error(t, err)
}

func TestDecimal_UnmarshalMalformedNumber(t *testing.T) {
	var d types.Decimal
	err := json.Unmarshal([]byte(`{`), &d)
	assert.Error(t, err)
}

func TestDecimal_Inner(t *testing.T) {
	d := types.MustDecimal("5.5")
	inner := d.Inner()
	assert.Equal(t, "5.5", inner.String())
}

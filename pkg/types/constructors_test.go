package types_test

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

// Tests cover the constructor symmetry added in B3 — every named
// identifier type now has a New/Must/*From* triplet.

func TestMustTxHash_PanicsOnInvalidInput(t *testing.T) {
	good := "0x" + string(make([]byte, 64))
	// Replace the zero bytes with '0' runes to make it a valid 64-hex string.
	b := []byte(good)
	for i := 2; i < len(b); i++ {
		b[i] = '0'
	}
	good = string(b)
	require.NotPanics(t, func() { types.MustTxHash(good) })
	require.Panics(t, func() { types.MustTxHash("not-a-hash") })
}

func TestTxHashFromCommon_NoStringRoundTrip(t *testing.T) {
	hash := common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	got := types.TxHashFromCommon(hash)
	assert.Equal(t, hash.Hex(), got.String())
}

func TestAddressFromCommon_NoStringRoundTrip(t *testing.T) {
	a := common.HexToAddress("0x1111111111111111111111111111111111111111")
	got := types.AddressFromCommon(a)
	assert.Equal(t, a.Hex(), got.String())
	assert.Equal(t, a, got.Common())
}

func TestDecimalFromShopspring_NoStringRoundTrip(t *testing.T) {
	d := decimal.RequireFromString("1.5")
	got := types.DecimalFromShopspring(d)
	assert.Equal(t, "1.5", got.String())
	assert.True(t, d.Equal(got.Inner()))
}

func TestNewMillisTime_ParsesEpochMillisString(t *testing.T) {
	got, err := types.NewMillisTime("1700000000000")
	require.NoError(t, err)
	assert.Equal(t, int64(1700000000000), got.Millis())
}

func TestNewMillisTime_ParsesRFC3339(t *testing.T) {
	got, err := types.NewMillisTime("2023-11-14T22:13:20Z")
	require.NoError(t, err)
	assert.Equal(t, int64(1700000000000), got.Millis())
}

func TestNewMillisTime_ParsesRFC3339Nano(t *testing.T) {
	got, err := types.NewMillisTime("2023-11-14T22:13:20.123456789Z")
	require.NoError(t, err)
	// Verify nanoseconds aren't lost on parse — the wire emits millis on
	// MarshalJSON, but the Go value can carry more precision.
	assert.Equal(t, 123456789, got.Time().Nanosecond())
}

func TestNewMillisTime_EmptyStringReturnsZero(t *testing.T) {
	got, err := types.NewMillisTime("")
	require.NoError(t, err)
	assert.True(t, got.Time().IsZero())
}

func TestNewMillisTime_RejectsGarbage(t *testing.T) {
	_, err := types.NewMillisTime("not a date")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "MillisTime")
}

func TestMustMillisTime_PanicsOnGarbage(t *testing.T) {
	require.NotPanics(t, func() { types.MustMillisTime("1700000000000") })
	require.Panics(t, func() { types.MustMillisTime("not a date") })
}

func TestMillisTimeFromTime_NoStringRoundTrip(t *testing.T) {
	now := time.Unix(1700000000, 0).UTC()
	got := types.MillisTimeFromTime(now)
	assert.Equal(t, now, got.Time())
}

func TestMillisTimeFromMillis_DirectEpochMillis(t *testing.T) {
	got := types.MillisTimeFromMillis(1700000000000)
	assert.Equal(t, int64(1700000000000), got.Millis())
}

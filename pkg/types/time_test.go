package types_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestMillisTime_RoundTripFromNumber(t *testing.T) {
	now := time.Date(2026, 5, 2, 10, 0, 0, 0, time.UTC)
	mt := types.NewMillisTime(now)
	b, err := json.Marshal(mt)
	require.NoError(t, err)

	var got types.MillisTime
	require.NoError(t, json.Unmarshal(b, &got))
	assert.Equal(t, now.UnixMilli(), got.Millis())
	assert.Equal(t, now.UnixMilli(), got.Time().UnixMilli())
}

func TestMillisTime_UnmarshalString(t *testing.T) {
	var mt types.MillisTime
	require.NoError(t, json.Unmarshal([]byte(`"1700000000000"`), &mt))
	assert.Equal(t, int64(1700000000000), mt.Millis())
}

func TestMillisTime_UnmarshalNullEmpty(t *testing.T) {
	var mt types.MillisTime
	require.NoError(t, json.Unmarshal([]byte(`null`), &mt))
	assert.True(t, mt.Time().IsZero())

	require.NoError(t, json.Unmarshal([]byte(`""`), &mt))
	assert.True(t, mt.Time().IsZero())
}

func TestMillisTime_UnmarshalInvalid(t *testing.T) {
	var mt types.MillisTime
	assert.Error(t, json.Unmarshal([]byte(`"abc"`), &mt))
	assert.Error(t, json.Unmarshal([]byte(`{`), &mt))
}

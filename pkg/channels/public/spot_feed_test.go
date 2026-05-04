package public_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/channels/public"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestSpotFeed_Name(t *testing.T) {
	assert.Equal(t, "spot_feed.BTC", public.SpotFeed{Currency: "BTC"}.Name())
	assert.Equal(t, "spot_feed.ETH", public.SpotFeed{Currency: "ETH"}.Name())
}

func TestSpotFeed_Decode(t *testing.T) {
	raw := json.RawMessage(`{"timestamp":1,"feeds":{"BTC":{"price":"100","confidence":"1","price_prev_daily":"99","confidence_prev_daily":"1","timestamp_prev_daily":0}}}`)
	v, err := public.SpotFeed{}.Decode(raw)
	require.NoError(t, err)
	sf, ok := v.(types.SpotFeed)
	require.True(t, ok)
	assert.Equal(t, "100", sf.Feeds["BTC"].Price.String())
}

func TestSpotFeed_Decode_Malformed(t *testing.T) {
	_, err := public.SpotFeed{}.Decode([]byte(`{`))
	assert.Error(t, err)
}

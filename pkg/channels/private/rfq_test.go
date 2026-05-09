package private_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/pkg/channels/private"
)

func TestRFQs_Name(t *testing.T) {
	assert.Equal(t,
		"wallet.0xC0FFee0000000000000000000000000000000000.rfqs",
		private.RFQs{Wallet: "0xC0FFee0000000000000000000000000000000000"}.Name())
}

func TestRFQs_Decode_Happy(t *testing.T) {
	raw := json.RawMessage(`[{"rfq_id":"R","subaccount_id":9,"status":"open","legs":[],"creation_timestamp":1,"last_update_timestamp":1}]`)
	v, err := private.RFQs{}.Decode(raw)
	require.NoError(t, err)
	rfqs, ok := v.([]derive.RFQ)
	require.True(t, ok)
	require.Len(t, rfqs, 1)
	assert.Equal(t, "R", rfqs[0].RFQID)
}

func TestRFQs_Decode_Malformed(t *testing.T) {
	_, err := private.RFQs{}.Decode([]byte(`{`))
	assert.Error(t, err)
}

func TestQuotes_Name(t *testing.T) {
	assert.Equal(t, "subaccount.9.quotes", private.Quotes{SubaccountID: 9}.Name())
}

func TestQuotes_Decode_Happy(t *testing.T) {
	raw := json.RawMessage(`[{"quote_id":"Q","rfq_id":"R","subaccount_id":9,"direction":"buy","legs":[],"price":"1","status":"open","creation_timestamp":1}]`)
	v, err := private.Quotes{}.Decode(raw)
	require.NoError(t, err)
	quotes, ok := v.([]derive.Quote)
	require.True(t, ok)
	require.Len(t, quotes, 1)
}

func TestQuotes_Decode_Malformed(t *testing.T) {
	_, err := private.Quotes{}.Decode([]byte(`{`))
	assert.Error(t, err)
}

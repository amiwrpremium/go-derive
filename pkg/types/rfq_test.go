package types_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func validLeg() types.RFQLeg {
	return types.RFQLeg{
		InstrumentName: "BTC-PERP",
		Direction:      enums.DirectionBuy,
		Amount:         types.MustDecimal("1"),
	}
}

func TestRFQLeg_Validate_Happy(t *testing.T) {
	require.NoError(t, validLeg().Validate())
}

func TestRFQLeg_Validate_Rejects(t *testing.T) {
	cases := []struct {
		name string
		mut  func(*types.RFQLeg)
		want string
	}{
		{"empty instrument", func(l *types.RFQLeg) { l.InstrumentName = "" }, "instrument_name"},
		{"bad direction", func(l *types.RFQLeg) { l.Direction = enums.Direction("sideways") }, "direction"},
		{"zero amount", func(l *types.RFQLeg) { l.Amount = types.MustDecimal("0") }, "amount"},
		{"negative amount", func(l *types.RFQLeg) { l.Amount = types.MustDecimal("-1") }, "amount"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := validLeg()
			c.mut(&l)
			err := l.Validate()
			require.Error(t, err)
			assert.True(t, errors.Is(err, types.ErrInvalidParams))
			assert.Contains(t, err.Error(), c.want)
		})
	}
}

func TestRFQLeg_RoundTrip(t *testing.T) {
	// RFQ legs do not carry per-leg prices — that's the QuoteLeg shape.
	in := types.RFQLeg{
		InstrumentName: "BTC-PERP",
		Direction:      enums.DirectionBuy,
		Amount:         types.MustDecimal("1"),
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	var out types.RFQLeg
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.InstrumentName, out.InstrumentName)
	assert.Equal(t, in.Direction, out.Direction)
}

func TestQuoteLeg_RoundTrip(t *testing.T) {
	in := types.QuoteLeg{
		InstrumentName: "BTC-PERP",
		Direction:      enums.DirectionBuy,
		Amount:         types.MustDecimal("1"),
		Price:          types.MustDecimal("65000"),
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	var out types.QuoteLeg
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, "65000", out.Price.String())
}

func TestRFQ_Decode(t *testing.T) {
	payload := `{
		"rfq_id": "R1",
		"subaccount_id": 1,
		"status": "open",
		"legs": [
			{"instrument_name":"BTC-PERP","direction":"buy","amount":"1"}
		],
		"max_total_fee": "10",
		"creation_timestamp": 1700000000000,
		"last_update_timestamp": 1700000000001
	}`
	var rfq types.RFQ
	require.NoError(t, json.Unmarshal([]byte(payload), &rfq))
	assert.Equal(t, "R1", rfq.RFQID)
	assert.Equal(t, enums.QuoteStatusOpen, rfq.Status)
	require.Len(t, rfq.Legs, 1)
}

func TestQuote_Decode(t *testing.T) {
	payload := `{
		"quote_id": "Q1",
		"rfq_id": "R1",
		"subaccount_id": 1,
		"direction": "sell",
		"legs": [{"instrument_name":"BTC-PERP","direction":"sell","amount":"1","price":"65000"}],
		"price": "65000",
		"status": "open",
		"creation_timestamp": 1700000000000
	}`
	var q types.Quote
	require.NoError(t, json.Unmarshal([]byte(payload), &q))
	assert.Equal(t, "Q1", q.QuoteID)
	assert.Equal(t, enums.DirectionSell, q.Direction)
	assert.Equal(t, enums.QuoteStatusOpen, q.Status)
	require.Len(t, q.Legs, 1)
	assert.Equal(t, "65000", q.Legs[0].Price.String())
}

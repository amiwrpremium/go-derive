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
		"wallet": "0xabc",
		"status": "open",
		"cancel_reason": "",
		"legs": [
			{"instrument_name":"BTC-PERP","direction":"buy","amount":"1"}
		],
		"counterparties": ["0xmaker1","0xmaker2"],
		"label": "spread-1",
		"preferred_direction": "buy",
		"reducing_direction": "",
		"filled_direction": "",
		"filled_pct": "0",
		"max_total_cost": "10",
		"min_total_cost": "0",
		"total_cost": "0",
		"ask_total_cost": "10",
		"bid_total_cost": "9",
		"mark_total_cost": "9.5",
		"partial_fill_step": "0.1",
		"valid_until": 1700000060000,
		"creation_timestamp": 1700000000000,
		"last_update_timestamp": 1700000000001
	}`
	var rfq types.RFQ
	require.NoError(t, json.Unmarshal([]byte(payload), &rfq))
	assert.Equal(t, "R1", rfq.RFQID)
	assert.Equal(t, "0xabc", rfq.Wallet)
	assert.Equal(t, enums.QuoteStatusOpen, rfq.Status)
	require.Len(t, rfq.Legs, 1)
	assert.Equal(t, []string{"0xmaker1", "0xmaker2"}, rfq.Counterparties)
	assert.Equal(t, "spread-1", rfq.Label)
	assert.Equal(t, "10", rfq.MaxTotalCost.String())
	assert.Equal(t, "9", rfq.BidTotalCost.String())
	assert.Equal(t, int64(1700000060000), rfq.ValidUntil.Millis())
}

func TestQuote_Decode(t *testing.T) {
	payload := `{
		"quote_id": "Q1",
		"rfq_id": "R1",
		"subaccount_id": 1,
		"direction": "sell",
		"legs": [{"instrument_name":"BTC-PERP","direction":"sell","amount":"1","price":"65000"}],
		"legs_hash": "0xleghash",
		"status": "open",
		"cancel_reason": "",
		"liquidity_role": "maker",
		"fee": "5",
		"max_fee": "10",
		"extra_fee": "1",
		"fill_pct": "0",
		"is_transfer": false,
		"label": "lbl",
		"mmp": false,
		"nonce": 42,
		"signer": "0x0000000000000000000000000000000000000001",
		"signature": "0xsig",
		"signature_expiry_sec": 1700000300,
		"tx_hash": "",
		"tx_status": "",
		"creation_timestamp": 1700000000000,
		"last_update_timestamp": 1700000000000
	}`
	var q types.Quote
	require.NoError(t, json.Unmarshal([]byte(payload), &q))
	assert.Equal(t, "Q1", q.QuoteID)
	assert.Equal(t, enums.DirectionSell, q.Direction)
	assert.Equal(t, enums.QuoteStatusOpen, q.Status)
	require.Len(t, q.Legs, 1)
	assert.Equal(t, "65000", q.Legs[0].Price.String())
	assert.Equal(t, "5", q.Fee.String())
	assert.Equal(t, "1", q.ExtraFee.String())
	assert.Equal(t, uint64(42), q.Nonce)
}

func TestQuotePublic_Decode(t *testing.T) {
	payload := `{
		"quote_id": "Q1",
		"rfq_id": "R1",
		"subaccount_id": 1,
		"wallet": "0xmaker",
		"direction": "sell",
		"legs": [{"instrument_name":"BTC-PERP","direction":"sell","amount":"1","price":"65000"}],
		"legs_hash": "0xleghash",
		"status": "open",
		"cancel_reason": "",
		"liquidity_role": "maker",
		"fill_pct": "0",
		"tx_hash": "",
		"tx_status": "",
		"creation_timestamp": 1700000000000,
		"last_update_timestamp": 1700000000000
	}`
	var q types.QuotePublic
	require.NoError(t, json.Unmarshal([]byte(payload), &q))
	assert.Equal(t, "Q1", q.QuoteID)
	assert.Equal(t, "0xmaker", q.Wallet)
	assert.Equal(t, enums.DirectionSell, q.Direction)
}

func TestExecuteQuoteResult_Decode(t *testing.T) {
	payload := `{
		"quote_id": "Q1","rfq_id": "R1","subaccount_id": 1,"direction":"sell",
		"legs":[{"instrument_name":"BTC-PERP","direction":"sell","amount":"1","price":"65000"}],
		"legs_hash":"","status":"open","cancel_reason":"","liquidity_role":"maker",
		"fee":"5","max_fee":"10","extra_fee":"0","fill_pct":"0.5","is_transfer":false,
		"label":"","mmp":false,"nonce":1,"signer":"0x0000000000000000000000000000000000000001",
		"signature":"","signature_expiry_sec":0,"tx_hash":"","tx_status":"",
		"creation_timestamp":1700000000000,"last_update_timestamp":1700000000000,
		"rfq_filled_pct": "0.25"
	}`
	var r types.ExecuteQuoteResult
	require.NoError(t, json.Unmarshal([]byte(payload), &r))
	assert.Equal(t, "0.25", r.RFQFilledPct.String())
	assert.Equal(t, "0.5", r.FillPct.String())
	assert.Equal(t, "Q1", r.QuoteID)
}

func TestCancelBatchResult_Decode(t *testing.T) {
	raw := []byte(`{"cancelled_ids":["a","b","c"]}`)
	var r types.CancelBatchResult
	require.NoError(t, json.Unmarshal(raw, &r))
	assert.Equal(t, []string{"a", "b", "c"}, r.CancelledIDs)
}

func TestBestQuoteResult_Decode(t *testing.T) {
	raw := []byte(`{
		"best_quote": null,
		"direction": "buy",
		"is_valid": true,
		"invalid_reason": "",
		"estimated_fee": "1",
		"estimated_realized_pnl": "0",
		"estimated_realized_pnl_excl_fees": "0",
		"estimated_total_cost": "1000",
		"filled_pct": "0",
		"orderbook_total_cost": null,
		"suggested_max_fee": "5",
		"pre_initial_margin": "100",
		"post_initial_margin": "120",
		"post_liquidation_price": null,
		"down_liquidation_price": null,
		"up_liquidation_price": null
	}`)
	var r types.BestQuoteResult
	require.NoError(t, json.Unmarshal(raw, &r))
	assert.Nil(t, r.BestQuote)
	assert.True(t, r.IsValid)
	assert.Equal(t, "1000", r.EstimatedTotalCost.String())
	assert.Equal(t, enums.DirectionBuy, r.Direction)
}

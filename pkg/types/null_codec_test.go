package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

// Several Derive responses carry fields the OAS marks as
// "required + nullable" (e.g. `LiquidationAuction.end_timestamp` —
// REQUIRED on the wire but null when the auction is open). The Go
// SDK uses bare types (`Decimal`, `MillisTime`, `enums.Direction`,
// `enums.TxStatus`) for these rather than `*T`, on the basis that
// each codec's UnmarshalJSON tolerates JSON `null` and decodes to
// the zero value.
//
// These tests pin that contract — if a future codec change starts
// rejecting `null`, the next response with a nullable field will
// break decode and these tests catch it.

func TestCodec_AcceptsNull_Decimal(t *testing.T) {
	var d types.Decimal
	require.NoError(t, json.Unmarshal([]byte("null"), &d))
	assert.Equal(t, "0", d.String(), "JSON null must decode to zero-value Decimal")
}

func TestCodec_AcceptsNull_MillisTime(t *testing.T) {
	var m types.MillisTime
	require.NoError(t, json.Unmarshal([]byte("null"), &m))
	assert.True(t, m.Time().IsZero(), "JSON null must decode to zero-value MillisTime")
}

func TestCodec_AcceptsNull_Address(t *testing.T) {
	var a types.Address
	require.NoError(t, json.Unmarshal([]byte("null"), &a))
}

func TestCodec_AcceptsNull_TxHash(t *testing.T) {
	var h types.TxHash
	require.NoError(t, json.Unmarshal([]byte("null"), &h))
}

func TestCodec_AcceptsNull_NestedStruct(t *testing.T) {
	// A representative case: every doc-nullable field on
	// `BestQuoteResult` arrives as `null`. The whole struct must
	// decode without error and every `Decimal` must end up zero.
	raw := []byte(`{
		"best_quote": null,
		"direction": "buy",
		"is_valid": true,
		"invalid_reason": null,
		"estimated_fee": "0",
		"estimated_realized_pnl": "0",
		"estimated_realized_pnl_excl_fees": "0",
		"estimated_total_cost": "0",
		"filled_pct": "0",
		"orderbook_total_cost": null,
		"suggested_max_fee": "0",
		"pre_initial_margin": "0",
		"post_initial_margin": "0",
		"post_liquidation_price": null,
		"down_liquidation_price": null,
		"up_liquidation_price": null
	}`)
	var got types.BestQuoteResult
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.Nil(t, got.BestQuote)
	assert.Equal(t, "0", got.OrderbookTotalCost.String())
	assert.Equal(t, "0", got.PostLiquidationPrice.String())
	assert.Equal(t, "0", got.DownLiquidationPrice.String())
	assert.Equal(t, "0", got.UpLiquidationPrice.String())
}

func TestCodec_AcceptsNull_RFQNullables(t *testing.T) {
	// `RFQResultSchema` declares ten fields as required-and-nullable
	// (cost columns + the three direction discriminators). A response
	// with all of them null must decode cleanly.
	raw := []byte(`{
		"rfq_id":"R1",
		"subaccount_id":1,
		"wallet":"0x1",
		"status":"open",
		"cancel_reason":"",
		"legs":[],
		"counterparties":null,
		"label":"",
		"preferred_direction":null,
		"reducing_direction":null,
		"filled_direction":null,
		"filled_pct":"0",
		"max_total_cost":null,
		"min_total_cost":null,
		"total_cost":null,
		"ask_total_cost":null,
		"bid_total_cost":null,
		"mark_total_cost":null,
		"partial_fill_step":"0",
		"valid_until":1700000060000,
		"creation_timestamp":1700000000000,
		"last_update_timestamp":1700000000001
	}`)
	var got types.RFQ
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.Equal(t, "R1", got.RFQID)
	assert.Equal(t, "0", got.AskTotalCost.String())
	assert.Equal(t, enums.Direction(""), got.FilledDirection)
	assert.Empty(t, got.Counterparties)
}

func TestCodec_AcceptsNull_LiquidationEndTimestamp(t *testing.T) {
	// `end_timestamp` is REQUIRED + nullable on AuctionHistoryResultSchema.
	// An open auction emits `null`; decode must succeed and the
	// MillisTime must be zero.
	raw := []byte(`{
		"auction_id":"a",
		"auction_type":"solvent",
		"bids":[],
		"end_timestamp":null,
		"fee":"0",
		"start_timestamp":1,
		"subaccount_id":1,
		"tx_hash":"0x0000000000000000000000000000000000000000000000000000000000000000"
	}`)
	var got types.LiquidationAuction
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.True(t, got.EndTimestamp.Time().IsZero())
}

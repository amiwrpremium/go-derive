package methods_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestSendRFQ_Happy(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/send_rfq", map[string]any{
		"rfq_id": "R1", "subaccount_id": 1, "status": "open",
		"legs": []any{}, "creation_timestamp": 1, "last_update_timestamp": 1,
	})
	rfq, err := api.SendRFQ(context.Background(), nil, types.MustDecimal("100"))
	require.NoError(t, err)
	assert.Equal(t, "R1", rfq.RFQID)
}

func TestSendRFQ_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	_, err := api.SendRFQ(context.Background(), nil, types.MustDecimal("0"))
	assert.ErrorIs(t, err, derrors.ErrSubaccountRequired)
}

func TestPollRFQs_Happy(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/poll_rfqs", map[string]any{"rfqs": []any{}})
	got, err := api.PollRFQs(context.Background())
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestPollRFQs_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	_, err := api.PollRFQs(context.Background())
	assert.ErrorIs(t, err, derrors.ErrSubaccountRequired)
}

func TestCancelRFQ_Happy(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/cancel_rfq", nil)
	require.NoError(t, api.CancelRFQ(context.Background(), "R1"))
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, "R1", params["rfq_id"])
}

func TestCancelRFQ_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	err := api.CancelRFQ(context.Background(), "R1")
	assert.ErrorIs(t, err, derrors.ErrSubaccountRequired)
}

// Decode coverage for the eight RFQ-flow methods that were retyped
// in the rfq_extras.go fold.

func TestGetRFQs_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 7)
	ft.HandleResult("private/get_rfqs", map[string]any{
		"rfqs": []any{
			map[string]any{
				"rfq_id": "R1", "subaccount_id": int64(7), "wallet": "0xa",
				"status": "open", "cancel_reason": "", "legs": []any{},
				"counterparties": []any{}, "label": "",
				"preferred_direction": "", "reducing_direction": "",
				"filled_direction": "", "filled_pct": "0",
				"max_total_cost": "10", "min_total_cost": "0",
				"total_cost": "0", "ask_total_cost": "0", "bid_total_cost": "0",
				"mark_total_cost": "0", "partial_fill_step": "0",
				"valid_until":        int64(1700000060000),
				"creation_timestamp": int64(1700000000000), "last_update_timestamp": int64(1700000000001),
			},
		},
		"pagination": map[string]any{"count": 1, "num_pages": 1},
	})
	rfqs, page, err := api.GetRFQs(context.Background(), nil)
	require.NoError(t, err)
	require.Len(t, rfqs, 1)
	assert.Equal(t, "R1", rfqs[0].RFQID)
	assert.Equal(t, 1, page.Count)
}

func TestGetRFQs_RequiresSigner(t *testing.T) {
	api, _ := newAPI(t, false, 0)
	_, _, err := api.GetRFQs(context.Background(), nil)
	assert.ErrorIs(t, err, derrors.ErrUnauthorized)
}

func TestGetQuotes_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 7)
	ft.HandleResult("private/get_quotes", map[string]any{
		"quotes":     []any{},
		"pagination": map[string]any{"count": 0, "num_pages": 0},
	})
	got, page, err := api.GetQuotes(context.Background(), nil)
	require.NoError(t, err)
	assert.Empty(t, got)
	assert.Equal(t, 0, page.NumPages)
}

func TestPollQuotes_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 7)
	ft.HandleResult("private/poll_quotes", map[string]any{
		"quotes":     []any{},
		"pagination": map[string]any{"count": 0, "num_pages": 0},
	})
	got, _, err := api.PollQuotes(context.Background(), nil)
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestSendQuote_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 7)
	ft.HandleResult("private/send_quote", map[string]any{
		"quote_id": "Q1", "rfq_id": "R1", "subaccount_id": int64(7),
		"direction": "sell", "legs": []any{}, "legs_hash": "",
		"status": "open", "cancel_reason": "", "liquidity_role": "maker",
		"fee": "5", "max_fee": "10", "extra_fee": "0", "fill_pct": "0",
		"is_transfer": false, "label": "", "mmp": false, "nonce": int64(1),
		"signer":    "0x0000000000000000000000000000000000000001",
		"signature": "0x", "signature_expiry_sec": int64(0),
		"tx_hash": "", "tx_status": "",
		"creation_timestamp":    int64(1700000000000),
		"last_update_timestamp": int64(1700000000000),
	})
	q, err := api.SendQuote(context.Background(), map[string]any{"rfq_id": "R1"})
	require.NoError(t, err)
	assert.Equal(t, "Q1", q.QuoteID)
}

func TestExecuteQuote_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 7)
	ft.HandleResult("private/execute_quote", map[string]any{
		"quote_id": "Q1", "rfq_id": "R1", "subaccount_id": int64(7),
		"direction": "buy", "legs": []any{}, "legs_hash": "",
		"status": "filled", "cancel_reason": "", "liquidity_role": "taker",
		"fee": "0", "max_fee": "0", "extra_fee": "0", "fill_pct": "1",
		"is_transfer": false, "label": "", "mmp": false, "nonce": int64(1),
		"signer":    "0x0000000000000000000000000000000000000001",
		"signature": "", "signature_expiry_sec": int64(0),
		"tx_hash": "", "tx_status": "",
		"creation_timestamp":    int64(1700000000000),
		"last_update_timestamp": int64(1700000000000),
		"rfq_filled_pct":        "0.5",
	})
	res, err := api.ExecuteQuote(context.Background(), map[string]any{"quote_id": "Q1"})
	require.NoError(t, err)
	assert.Equal(t, "0.5", res.RFQFilledPct.String())
}

func TestCancelQuote_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 7)
	ft.HandleResult("private/cancel_quote", map[string]any{
		"quote_id": "Q1", "rfq_id": "R1", "subaccount_id": int64(7),
		"direction": "sell", "legs": []any{}, "legs_hash": "",
		"status": "cancelled", "cancel_reason": "user_request", "liquidity_role": "maker",
		"fee": "0", "max_fee": "0", "extra_fee": "0", "fill_pct": "0",
		"is_transfer": false, "label": "", "mmp": false, "nonce": int64(1),
		"signer":    "0x0000000000000000000000000000000000000001",
		"signature": "", "signature_expiry_sec": int64(0),
		"tx_hash": "", "tx_status": "",
		"creation_timestamp":    int64(1700000000000),
		"last_update_timestamp": int64(1700000000000),
	})
	q, err := api.CancelQuote(context.Background(), map[string]any{"quote_id": "Q1"})
	require.NoError(t, err)
	assert.Equal(t, "cancelled", string(q.Status))
}

func TestReplaceQuote_Decode_HappyPath(t *testing.T) {
	quote := map[string]any{
		"quote_id": "Q2", "rfq_id": "R1", "subaccount_id": int64(7),
		"direction": "sell", "legs": []any{}, "legs_hash": "",
		"status": "open", "cancel_reason": "", "liquidity_role": "maker",
		"fee": "0", "max_fee": "10", "extra_fee": "0", "fill_pct": "0",
		"is_transfer": false, "label": "", "mmp": false, "nonce": int64(2),
		"signer":    "0x0000000000000000000000000000000000000001",
		"signature": "0x", "signature_expiry_sec": int64(0),
		"tx_hash": "", "tx_status": "",
		"creation_timestamp":    int64(1700000000000),
		"last_update_timestamp": int64(1700000000000),
	}
	api, ft := newAPI(t, true, 7)
	ft.HandleResult("private/replace_quote", map[string]any{
		"cancelled_quote": map[string]any{
			"quote_id": "Q1", "rfq_id": "R1", "subaccount_id": int64(7),
			"direction": "sell", "legs": []any{}, "legs_hash": "",
			"status": "cancelled", "cancel_reason": "user_request",
			"liquidity_role": "maker",
			"fee":            "0", "max_fee": "10", "extra_fee": "0", "fill_pct": "0",
			"is_transfer": false, "label": "", "mmp": false, "nonce": int64(1),
			"signer":    "0x0000000000000000000000000000000000000001",
			"signature": "0x", "signature_expiry_sec": int64(0),
			"tx_hash": "", "tx_status": "",
			"creation_timestamp":    int64(1700000000000),
			"last_update_timestamp": int64(1700000000000),
		},
		"quote":              quote,
		"create_quote_error": nil,
	})
	res, err := api.ReplaceQuote(context.Background(), map[string]any{
		"rfq_id": "R1", "quote_id_to_cancel": "Q1",
	})
	require.NoError(t, err)
	assert.Equal(t, "Q1", res.CancelledQuote.QuoteID)
	require.NotNil(t, res.Quote)
	assert.Equal(t, "Q2", res.Quote.QuoteID)
	assert.Nil(t, res.CreateQuoteError)
}

func TestReplaceQuote_Decode_RejectedReplacement(t *testing.T) {
	api, ft := newAPI(t, true, 7)
	ft.HandleResult("private/replace_quote", map[string]any{
		"cancelled_quote": map[string]any{
			"quote_id": "Q1", "rfq_id": "R1", "subaccount_id": int64(7),
			"direction": "sell", "legs": []any{}, "legs_hash": "",
			"status": "cancelled", "cancel_reason": "user_request",
			"liquidity_role": "maker",
			"fee":            "0", "max_fee": "10", "extra_fee": "0", "fill_pct": "0",
			"is_transfer": false, "label": "", "mmp": false, "nonce": int64(1),
			"signer":    "0x0000000000000000000000000000000000000001",
			"signature": "0x", "signature_expiry_sec": int64(0),
			"tx_hash": "", "tx_status": "",
			"creation_timestamp":    int64(1700000000000),
			"last_update_timestamp": int64(1700000000000),
		},
		"quote":              nil,
		"create_quote_error": map[string]any{"code": -32000, "message": "insufficient_margin"},
	})
	res, err := api.ReplaceQuote(context.Background(), map[string]any{
		"rfq_id": "R1", "quote_id_to_cancel": "Q1",
	})
	require.NoError(t, err)
	assert.Nil(t, res.Quote)
	require.NotNil(t, res.CreateQuoteError)
	assert.Equal(t, -32000, res.CreateQuoteError.Code)
	assert.Equal(t, "insufficient_margin", res.CreateQuoteError.Message)
}

func TestReplaceQuote_RequiresSigner(t *testing.T) {
	api, _ := newAPI(t, false, 0)
	_, err := api.ReplaceQuote(context.Background(), nil)
	assert.ErrorIs(t, err, derrors.ErrUnauthorized)
}

func TestCancelBatchQuotes_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 7)
	ft.HandleResult("private/cancel_batch_quotes", map[string]any{
		"cancelled_ids": []any{"a", "b", "c"},
	})
	got, err := api.CancelBatchQuotes(context.Background(), nil)
	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "c"}, got.CancelledIDs)
}

func TestCancelBatchRFQs_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 7)
	ft.HandleResult("private/cancel_batch_rfqs", map[string]any{
		"cancelled_ids": []any{"x"},
	})
	got, err := api.CancelBatchRFQs(context.Background(), nil)
	require.NoError(t, err)
	assert.Equal(t, []string{"x"}, got.CancelledIDs)
}

func TestRFQGetBestQuote_Decode_NoQuote(t *testing.T) {
	api, ft := newAPI(t, true, 7)
	ft.HandleResult("private/rfq_get_best_quote", map[string]any{
		"best_quote":                       nil,
		"direction":                        "buy",
		"is_valid":                         true,
		"invalid_reason":                   "",
		"estimated_fee":                    "1",
		"estimated_realized_pnl":           "0",
		"estimated_realized_pnl_excl_fees": "0",
		"estimated_total_cost":             "1000",
		"filled_pct":                       "0",
		"orderbook_total_cost":             nil,
		"suggested_max_fee":                "5",
		"pre_initial_margin":               "100",
		"post_initial_margin":              "120",
		"post_liquidation_price":           nil,
		"down_liquidation_price":           nil,
		"up_liquidation_price":             nil,
	})
	res, err := api.RFQGetBestQuote(context.Background(), map[string]any{"legs": []any{}, "direction": "buy"})
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Nil(t, res.BestQuote)
	assert.True(t, res.IsValid)
	assert.Equal(t, "1000", res.EstimatedTotalCost.String())
}

func TestOrderQuote_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 7)
	ft.HandleResult("private/order_quote", map[string]any{
		"is_valid":               true,
		"invalid_reason":         nil,
		"estimated_fill_amount":  "1",
		"estimated_fill_price":   "50000",
		"estimated_fee":          "5",
		"estimated_realized_pnl": "0",
		"estimated_order_status": "filled",
		"suggested_max_fee":      "10",
		"pre_initial_margin":     "100",
		"post_initial_margin":    "120",
		"post_liquidation_price": nil,
	})
	got, err := api.OrderQuote(context.Background(), map[string]any{})
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.True(t, got.IsValid)
	assert.Equal(t, "filled", got.EstimatedOrderStatus)
	assert.Equal(t, "5", got.EstimatedFee.String())
	assert.Equal(t, "0", got.PostLiquidationPrice.String(), "null decimal decodes to zero-value")
}

func TestOrderQuote_Invalid(t *testing.T) {
	api, ft := newAPI(t, true, 7)
	ft.HandleResult("private/order_quote", map[string]any{
		"is_valid":               false,
		"invalid_reason":         "Insufficient buying power.",
		"estimated_fill_amount":  "0",
		"estimated_fill_price":   "0",
		"estimated_fee":          "0",
		"estimated_realized_pnl": "0",
		"estimated_order_status": "rejected",
		"suggested_max_fee":      "0",
		"pre_initial_margin":     "100",
		"post_initial_margin":    "100",
		"post_liquidation_price": "45000",
	})
	got, err := api.OrderQuote(context.Background(), map[string]any{})
	require.NoError(t, err)
	assert.False(t, got.IsValid)
	assert.Equal(t, enums.RFQInvalidReasonInsufficientBuyingPower, got.InvalidReason)
	assert.Equal(t, "rejected", got.EstimatedOrderStatus)
	assert.Equal(t, "45000", got.PostLiquidationPrice.String())
}

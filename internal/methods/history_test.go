package methods_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
)

func TestGetFundingHistory_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 7)
	ft.HandleResult("private/get_funding_history", map[string]any{
		"events": []any{
			map[string]any{"instrument_name": "BTC-PERP", "subaccount_id": int64(7), "timestamp": int64(1700000000000), "funding": "0.001", "pnl": "-0.0005"},
		},
		"pagination": map[string]any{"count": 1, "num_pages": 1},
	})
	events, page, err := api.GetFundingHistory(context.Background(), nil)
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, "BTC-PERP", events[0].InstrumentName)
	assert.Equal(t, "0.001", events[0].Funding.String())
	assert.Equal(t, 1, page.Count)
	assert.Equal(t, 1, page.NumPages)
}

func TestGetFundingHistory_RequiresSigner(t *testing.T) {
	api, _ := newAPI(t, false, 0)
	_, _, err := api.GetFundingHistory(context.Background(), nil)
	assert.True(t, errors.Is(err, derrors.ErrUnauthorized))
}

func TestGetFundingHistory_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	_, _, err := api.GetFundingHistory(context.Background(), nil)
	assert.True(t, errors.Is(err, derrors.ErrSubaccountRequired))
}

func TestGetFundingHistory_PropagatesAPIError(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleError("private/get_funding_history", boom)
	_, _, err := api.GetFundingHistory(context.Background(), nil)
	assert.ErrorAs(t, err, new(*derrors.APIError))
}

func TestGetLiquidatorHistory_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 9)
	ft.HandleResult("private/get_liquidator_history", map[string]any{
		"bids": []any{
			map[string]any{
				"amounts_liquidated":               map[string]any{"BTC-PERP": "0.5"},
				"cash_received":                    "100",
				"discount_pnl":                     "5",
				"percent_liquidated":               "0.5",
				"positions_realized_pnl":           map[string]any{"BTC-PERP": "10"},
				"positions_realized_pnl_excl_fees": map[string]any{"BTC-PERP": "11"},
				"realized_pnl":                     "10",
				"realized_pnl_excl_fees":           "11",
				"timestamp":                        int64(1700000000000),
				"tx_hash":                          "0x" + strings.Repeat("a", 64),
			},
		},
		"pagination": map[string]any{"count": 1, "num_pages": 1},
	})
	bids, page, err := api.GetLiquidatorHistory(context.Background(), nil)
	require.NoError(t, err)
	require.Len(t, bids, 1)
	assert.Equal(t, "100", bids[0].CashReceived.String())
	assert.Equal(t, 1, page.Count)
}

func TestGetLiquidatorHistory_RequiresSigner(t *testing.T) {
	api, _ := newAPI(t, false, 0)
	_, _, err := api.GetLiquidatorHistory(context.Background(), nil)
	assert.True(t, errors.Is(err, derrors.ErrUnauthorized))
}

func TestGetLiquidatorHistory_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	_, _, err := api.GetLiquidatorHistory(context.Background(), nil)
	assert.True(t, errors.Is(err, derrors.ErrSubaccountRequired))
}

func TestGetLiquidationHistory_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 9)
	ft.HandleResult("private/get_liquidation_history", []any{
		map[string]any{
			"auction_id": "auc-1", "auction_type": "solvent", "subaccount_id": int64(9),
			"start_timestamp": int64(1700000000000), "end_timestamp": int64(1700000060000),
			"fee": "1", "tx_hash": "0x1111111111111111111111111111111111111111111111111111111111111111",
			"bids": []any{},
		},
	})
	got, err := api.GetLiquidationHistory(context.Background(), nil)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "auc-1", got[0].AuctionID)
	assert.Equal(t, "solvent", got[0].AuctionType)
}

func TestGetLiquidationHistory_PropagatesAPIError(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleError("private/get_liquidation_history", boom)
	_, err := api.GetLiquidationHistory(context.Background(), nil)
	assert.ErrorAs(t, err, new(*derrors.APIError))
}

func TestGetOptionSettlementHistory_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 9)
	ft.HandleResult("private/get_option_settlement_history", map[string]any{
		"settlements": []any{
			map[string]any{
				"subaccount_id": int64(9), "instrument_name": "BTC-20240101-50000-C",
				"expiry": int64(1704067200), "amount": "1", "settlement_price": "50000",
				"option_settlement_pnl": "100", "option_settlement_pnl_excl_fees": "105",
			},
		},
	})
	got, err := api.GetOptionSettlementHistory(context.Background(), nil)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, int64(1704067200), got[0].Expiry)
}

func TestGetPublicOptionSettlementHistory_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_option_settlement_history", map[string]any{
		"settlements": []any{
			map[string]any{
				"subaccount_id": int64(1), "instrument_name": "BTC-20240101-50000-P",
				"expiry": int64(1704067200), "amount": "1", "settlement_price": "50000",
				"option_settlement_pnl": "0", "option_settlement_pnl_excl_fees": "0",
			},
		},
		"pagination": map[string]any{"count": 1, "num_pages": 1},
	})
	got, page, err := api.GetPublicOptionSettlementHistory(context.Background(), nil)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, 1, page.Count)
}

func TestGetSubaccountValueHistory_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 9)
	ft.HandleResult("private/get_subaccount_value_history", map[string]any{
		"subaccount_id": int64(9),
		"subaccount_value_history": []any{
			map[string]any{
				"timestamp": int64(1700000000000), "currency": "USDC", "margin_type": "PM",
				"subaccount_value": "10000", "initial_margin": "100",
				"maintenance_margin": "50", "delayed_maintenance_margin": "55",
			},
		},
	})
	subaccountID, got, err := api.GetSubaccountValueHistory(context.Background(), map[string]any{
		"period": int64(86400), "start_timestamp": 0, "end_timestamp": 1,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(9), subaccountID)
	require.Len(t, got, 1)
	assert.Equal(t, "PM", got[0].MarginType)
	assert.Equal(t, "10000", got[0].SubaccountValue.String())
}

func TestGetERC20TransferHistory_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/get_erc20_transfer_history", map[string]any{
		"events": []any{
			map[string]any{
				"subaccount_id": int64(1), "counterparty_subaccount_id": int64(2),
				"asset": "USDC", "amount": "100", "is_outgoing": true,
				"timestamp": int64(1700000000000),
				"tx_hash":   "0x1111111111111111111111111111111111111111111111111111111111111111",
			},
		},
	})
	got, err := api.GetERC20TransferHistory(context.Background(), nil)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.True(t, got[0].IsOutgoing)
	assert.Equal(t, "USDC", got[0].Asset)
}

func TestGetInterestHistory_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/get_interest_history", map[string]any{
		"events": []any{
			map[string]any{"subaccount_id": int64(1), "timestamp": int64(1700000000000), "interest": "0.5"},
		},
	})
	got, err := api.GetInterestHistory(context.Background(), nil)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "0.5", got[0].Interest.String())
}

func TestExpiredAndCancelledHistory_Decode(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/expired_and_cancelled_history", map[string]any{
		"presigned_urls": []any{"https://s3.example/a", "https://s3.example/b"},
	})
	got, err := api.ExpiredAndCancelledHistory(context.Background(), map[string]any{
		"start_timestamp": 0, "end_timestamp": 1, "expiry": int64(1704067200),
	})
	require.NoError(t, err)
	require.Equal(t, []string{"https://s3.example/a", "https://s3.example/b"}, got.PresignedURLs)
}

func TestExpiredAndCancelledHistory_PropagatesAPIError(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleError("private/expired_and_cancelled_history", boom)
	_, err := api.ExpiredAndCancelledHistory(context.Background(), nil)
	assert.ErrorAs(t, err, new(*derrors.APIError))
}

package methods_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestGetInstruments(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_instruments", []map[string]any{
		{"instrument_name": "BTC-PERP", "instrument_type": "perp", "is_active": true,
			"base_currency": "BTC", "quote_currency": "USDC",
			"tick_size": "0.5", "minimum_amount": "0.001", "maximum_amount": "100", "amount_step": "0.001"},
	})

	insts, err := api.GetInstruments(context.Background(), types.InstrumentsQuery{Currency: "BTC", Kind: enums.InstrumentTypePerp})
	require.NoError(t, err)
	require.Len(t, insts, 1)
	assert.Equal(t, "BTC-PERP", insts[0].Name)

	last := ft.LastCall()
	assert.Equal(t, "public/get_instruments", last.Method)
	params := paramsAsMap(t, last.Params)
	assert.Equal(t, "BTC", params["currency"])
	assert.Equal(t, "perp", params["instrument_type"])
	assert.Equal(t, false, params["expired"])
}

func TestGetInstruments_NoFilters(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_instruments", []map[string]any{})

	_, err := api.GetInstruments(context.Background(), types.InstrumentsQuery{})
	require.NoError(t, err)
	params := paramsAsMap(t, ft.LastCall().Params)
	_, hasCurrency := params["currency"]
	_, hasKind := params["instrument_type"]
	assert.False(t, hasCurrency)
	assert.False(t, hasKind)
}

func TestGetInstrument(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_instrument", map[string]any{
		"instrument_name": "BTC-PERP",
		"instrument_type": "perp",
		"is_active":       true,
		"base_currency":   "BTC",
		"quote_currency":  "USDC",
		"tick_size":       "0.5",
		"minimum_amount":  "0.001",
		"maximum_amount":  "100",
		"amount_step":     "0.001",
	})

	got, err := api.GetInstrument(context.Background(), types.InstrumentQuery{Name: "BTC-PERP"})
	require.NoError(t, err)
	assert.Equal(t, "BTC-PERP", got.Name)
}

func TestGetTicker(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_ticker", map[string]any{
		"instrument_name": "BTC-PERP",
		"best_bid_price":  "100", "best_bid_amount": "1",
		"best_ask_price": "101", "best_ask_amount": "2",
		"mark_price":  "100.5",
		"index_price": "100.5",
		"timestamp":   1700000000000,
	})
	got, err := api.GetTicker(context.Background(), types.TickerQuery{Name: "BTC-PERP"})
	require.NoError(t, err)
	assert.Equal(t, "BTC-PERP", got.InstrumentName)
}

func TestGetPublicTradeHistory(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_trade_history", map[string]any{
		"trades":     []any{},
		"pagination": map[string]any{"num_pages": 1, "count": 0, "current_page": 1, "page_size": 50},
	})
	trades, page, err := api.GetPublicTradeHistory(context.Background(),
		types.PublicTradeHistoryQuery{
			InstrumentName: "BTC-PERP",
			Currency:       "BTC",
			TxStatus:       "settled",
			FromTimestamp:  types.NewMillisTime(time.UnixMilli(1700000000000)),
		},
		types.PageRequest{Page: 2, PageSize: 50})
	require.NoError(t, err)
	assert.Empty(t, trades)
	assert.Equal(t, 1, page.NumPages)
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, "BTC-PERP", params["instrument_name"])
	assert.Equal(t, "BTC", params["currency"])
	assert.Equal(t, "settled", params["tx_status"])
	assert.Equal(t, float64(1700000000000), params["from_timestamp"])
	assert.Equal(t, float64(2), params["page"])
	assert.Equal(t, float64(50), params["page_size"])
}

func TestGetTime(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_time", 1700000000000)
	got, err := api.GetTime(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(1700000000000), got)
}

func TestGetCurrencies(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_all_currencies", []map[string]any{
		{"currency": "USDC"}, {"currency": "WETH"},
	})
	got, err := api.GetCurrencies(context.Background())
	require.NoError(t, err)
	assert.Equal(t, []string{"USDC", "WETH"}, got)
}

func TestPublicMethods_PropagateUnhandledTransportError(t *testing.T) {
	api, _ := newAPI(t, false, 0)
	// No handler registered → FakeTransport returns "unhandled" error.
	_, err := api.GetTime(context.Background())
	require.Error(t, err)
}

func TestGetStatistics_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/statistics", map[string]any{
		"daily_fees":            "100",
		"daily_notional_volume": "1000000",
		"daily_premium_volume":  "50000",
		"daily_trades":          int64(250),
		"open_interest":         "500",
		"total_fees":            "10000",
		"total_notional_volume": "100000000",
		"total_premium_volume":  "500000",
		"total_trades":          int64(25000),
	})
	got, err := api.GetStatistics(context.Background(), types.StatisticsQuery{
		InstrumentName: "BTC-PERP",
		Currency:       "BTC",
		EndTime:        1700000000,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(250), got.DailyTrades)
	assert.Equal(t, "500", got.OpenInterest.String())
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, "BTC", params["currency"])
	assert.Equal(t, float64(1700000000), params["end_time"])
}

func TestGetCurrency_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_currency", map[string]any{
		"currency": "ETH", "spot_price": "2500", "spot_price_24h": "2450",
		"instrument_types": []any{"option", "perp"}, "market_type": "ALL",
		"managers": []any{}, "pm2_collateral_discounts": []any{},
		"protocol_asset_addresses":         map[string]any{},
		"asset_cap_and_supply_per_manager": map[string]any{},
		"srm_im_discount":                  "0", "srm_mm_discount": "0",
		"borrow_apy": "0", "supply_apy": "0", "total_borrow": "0", "total_supply": "0",
	})
	got, err := api.GetCurrency(context.Background(), types.CurrencyQuery{Currency: "ETH"})
	require.NoError(t, err)
	assert.Equal(t, "ETH", got.Currency)
	assert.Equal(t, "2500", got.SpotPrice.String())
	assert.Equal(t, "ETH", paramsAsMap(t, ft.LastCall().Params)["currency"])
}

func TestGetAllInstruments_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_all_instruments", map[string]any{
		"instruments": []any{
			map[string]any{
				"instrument_name": "BTC-PERP", "instrument_type": "perp", "is_active": true,
				"base_currency": "BTC", "quote_currency": "USDC",
				"tick_size": "0.5", "minimum_amount": "0.001", "maximum_amount": "100", "amount_step": "0.001",
			},
		},
		"pagination": map[string]any{"num_pages": 3, "count": 250},
	})
	insts, page, err := api.GetAllInstruments(context.Background(), types.AllInstrumentsQuery{Kind: enums.InstrumentTypePerp, IncludeExpired: true}, types.PageRequest{Page: 2, PageSize: 100})
	require.NoError(t, err)
	require.Len(t, insts, 1)
	assert.Equal(t, "BTC-PERP", insts[0].Name)
	assert.Equal(t, 3, page.NumPages)
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, "perp", params["instrument_type"])
	assert.Equal(t, true, params["expired"])
	assert.Equal(t, float64(2), params["page"])
	assert.Equal(t, float64(100), params["page_size"])
}

func TestGetTickers_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_tickers", map[string]any{
		"tickers": map[string]any{
			"BTC-PERP": map[string]any{
				"A": "1", "B": "2", "I": "65000", "M": "65000",
				"a": "65010", "b": "64990", "f": "0.0001", "maxp": "65500", "minp": "64500",
				"option_pricing": nil,
				"stats":          map[string]any{"c": "100", "h": "65500", "l": "64500", "n": int64(50), "oi": "10", "p": "0.01", "pr": "1000", "v": "5000000"},
				"t":              int64(1700000000000),
			},
		},
	})
	tickers, err := api.GetTickers(context.Background(), types.TickersQuery{InstrumentType: enums.InstrumentTypePerp})
	require.NoError(t, err)
	require.Contains(t, tickers, "BTC-PERP")
	assert.Equal(t, "65000", tickers["BTC-PERP"].MarkPrice.String())
}

func TestGetOptionSettlementPrices_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_option_settlement_prices", map[string]any{
		"expiries": []any{
			map[string]any{"expiry_date": "20260327", "utc_expiry_sec": int64(1774828800), "price": "65000"},
			map[string]any{"expiry_date": "20260626", "utc_expiry_sec": int64(1782604800), "price": nil},
		},
	})
	prices, err := api.GetOptionSettlementPrices(context.Background(), types.OptionSettlementPricesQuery{Currency: "BTC"})
	require.NoError(t, err)
	require.Len(t, prices, 2)
	assert.Equal(t, "65000", prices[0].Price.String())
	assert.Equal(t, "0", prices[1].Price.String())
}

func TestGetLiveIncidents_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_live_incidents", map[string]any{
		"incidents": []any{
			map[string]any{
				"creation_timestamp_sec": int64(1700000000),
				"label":                  "matching", "message": "elevated latency",
				"monitor_type": "auto", "severity": "medium",
			},
		},
	})
	got, err := api.GetLiveIncidents(context.Background())
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "matching", got[0].Label)
}

func TestGetAllStatistics_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/all_statistics", []any{
		map[string]any{
			"currency": "BTC", "instrument_type": "perp",
			"daily_fees": "100", "daily_notional_volume": "1000000",
			"daily_premium_volume": "0", "daily_trades": int64(250),
			"open_interest": "10", "total_fees": "100000",
			"total_notional_volume": "1000000000", "total_premium_volume": "0",
			"total_trades": int64(25000),
		},
	})
	got, err := api.GetAllStatistics(context.Background(), types.AllStatisticsQuery{})
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "BTC", got[0].Currency)
	assert.Equal(t, "perp", got[0].InstrumentType)
}

func TestGetAllUserStatistics_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/all_user_statistics", []any{
		map[string]any{
			"wallet":         "0x1111111111111111111111111111111111111111",
			"total_base_fee": "10", "total_contract_fee": "5", "total_fees": "15",
			"total_notional_volume": "100000", "total_premium_volume": "0",
			"total_regular_base_fee": "10", "total_regular_contract_fee": "5",
			"total_trades":          int64(7),
			"first_trade_timestamp": int64(1700000000000),
			"last_trade_timestamp":  int64(1700100000000),
		},
	})
	got, err := api.GetAllUserStatistics(context.Background(), types.AllUserStatisticsQuery{})
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "0x1111111111111111111111111111111111111111", got[0].Wallet)
	assert.Equal(t, "15", got[0].TotalFees.String())
}

func TestGetUserStatistics_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/user_statistics", map[string]any{
		"total_base_fee": "10", "total_contract_fee": "5", "total_fees": "15",
		"total_notional_volume": "100000", "total_premium_volume": "0",
		"total_regular_base_fee": "10", "total_regular_contract_fee": "5",
		"total_trades":          int64(7),
		"first_trade_timestamp": int64(1700000000000),
		"last_trade_timestamp":  int64(1700100000000),
	})
	got, err := api.GetUserStatistics(context.Background(), types.UserStatisticsQuery{Wallet: "0x1111111111111111111111111111111111111111"})
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, int64(7), got.TotalTrades)
}

func TestGetAsset_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_asset", map[string]any{
		"address":        "0x1111111111111111111111111111111111111111",
		"asset_id":       "1",
		"asset_name":     "USDC",
		"asset_type":     "erc20",
		"currency":       "USDC",
		"is_collateral":  true,
		"is_position":    false,
		"erc20_details":  map[string]any{"decimals": 6},
		"option_details": nil,
		"perp_details":   nil,
	})
	got, err := api.GetAsset(context.Background(), types.AssetQuery{Name: "USDC"})
	require.NoError(t, err)
	assert.Equal(t, "USDC", got.AssetName)
	assert.True(t, got.IsCollateral)
}

func TestGetAssets_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_assets", []any{
		map[string]any{
			"address":  "0x1111111111111111111111111111111111111111",
			"asset_id": "1", "asset_name": "USDC", "asset_type": "erc20",
			"currency": "USDC", "is_collateral": true, "is_position": false,
		},
	})
	got, err := api.GetAssets(context.Background(), types.AssetsQuery{AssetType: enums.AssetTypeERC20, Currency: "USDC"})
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "USDC", got[0].AssetName)
}

func TestGetBridgeBalances_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_bridge_balances", []any{
		map[string]any{
			"name": "across", "integrator": "across-protocol",
			"chain_id": int64(1), "balance": "1000000", "balance_hours": "72",
		},
	})
	got, err := api.GetBridgeBalances(context.Background())
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "across", got[0].Name)
}

func TestGetStDRVSnapshots_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_stdrv_snapshots", map[string]any{
		"wallet": "0x1111111111111111111111111111111111111111",
		"snapshots": []any{
			map[string]any{"amount": "1000", "timestamp_sec": int64(1700000000)},
			map[string]any{"amount": "1010", "timestamp_sec": int64(1700003600)},
		},
	})
	got, err := api.GetStDRVSnapshots(context.Background(), types.STDRVSnapshotsQuery{
		Wallet: "0x1111111111111111111111111111111111111111",
	})
	require.NoError(t, err)
	require.Len(t, got.Snapshots, 2)
	assert.Equal(t, "1010", got.Snapshots[1].Amount.String())
}

func TestGetDescendantTree_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_descendant_tree", map[string]any{
		"parent":      "0x1111111111111111111111111111111111111111",
		"descendants": []any{},
	})
	got, err := api.GetDescendantTree(context.Background(), types.DescendantTreeQuery{WalletOrInviteCode: "0x1111111111111111111111111111111111111111"})
	require.NoError(t, err)
	assert.Equal(t, "0x1111111111111111111111111111111111111111", got.Parent)
	assert.NotEmpty(t, got.Descendants, "descendants is preserved as raw JSON")
}

func TestGetTreeRoots_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_tree_roots", map[string]any{
		"roots": []any{},
	})
	got, err := api.GetTreeRoots(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, got.Roots)
}

func TestMarginWatch_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/margin_watch", map[string]any{
		"subaccount_id": int64(7), "currency": "USDC", "margin_type": "PM",
		"subaccount_value": "10000", "initial_margin": "100",
		"maintenance_margin": "50", "valuation_timestamp": int64(1700000000),
		"collaterals": []any{}, "positions": []any{},
	})
	got, err := api.MarginWatch(context.Background(), types.MarginWatchQuery{SubaccountID: 7})
	require.NoError(t, err)
	assert.Equal(t, int64(7), got.SubaccountID)
	assert.Equal(t, "PM", got.MarginType)
}

// silence unused json import warning on platforms where compiler is strict.
var _ = json.Marshal

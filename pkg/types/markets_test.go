package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestPortfolio_Decode(t *testing.T) {
	// Golden payload covering every required field on
	// PrivateGetSubaccountResultSchema, with one open order and one
	// position so the nested arrays are exercised.
	raw := []byte(`{
		"subaccount_id": 42,
		"currency": "USDC",
		"label": "trader-1",
		"margin_type": "PM",
		"is_under_liquidation": false,
		"subaccount_value": "1000",
		"initial_margin": "100",
		"maintenance_margin": "50",
		"open_orders_margin": "10",
		"projected_margin_change": "-5",
		"collaterals_initial_margin": "1000",
		"collaterals_maintenance_margin": "1000",
		"collaterals_value": "1000",
		"positions_initial_margin": "0",
		"positions_maintenance_margin": "0",
		"positions_value": "0",
		"collaterals": [{
			"asset_name":"USDC","asset_type":"erc20","amount":"1000","mark_value":"1000"
		}],
		"open_orders": [],
		"positions": []
	}`)
	var p types.Portfolio
	require.NoError(t, json.Unmarshal(raw, &p))
	assert.Equal(t, int64(42), p.SubaccountID)
	assert.Equal(t, "USDC", p.Currency)
	assert.Equal(t, "trader-1", p.Label)
	assert.Equal(t, "PM", p.MarginType)
	assert.Equal(t, "1000", p.SubaccountValue.String())
	assert.Equal(t, "100", p.InitialMargin.String())
	assert.Equal(t, "10", p.OpenOrdersMargin.String())
	assert.Equal(t, "-5", p.ProjectedMarginChange.String())
	assert.Equal(t, "1000", p.CollateralsValue.String())
	require.Len(t, p.Collaterals, 1)
	assert.Equal(t, "USDC", p.Collaterals[0].AssetName)
}

func TestCurrency_Decode(t *testing.T) {
	raw := []byte(`{
		"currency": "ETH",
		"spot_price": "2500",
		"spot_price_24h": "2450",
		"instrument_types": ["option","perp"],
		"market_type": "ALL",
		"managers": [{"address":"0x1111111111111111111111111111111111111111","margin_type":"PM2"}],
		"pm2_collateral_discounts": [{"manager_currency":"USDC","im_discount":"0.9","mm_discount":"0.95"}],
		"protocol_asset_addresses": {
			"option":"0x2222222222222222222222222222222222222222",
			"perp":"0x3333333333333333333333333333333333333333",
			"spot":"0x4444444444444444444444444444444444444444",
			"underlying_erc20":"0x5555555555555555555555555555555555555555"
		},
		"asset_cap_and_supply_per_manager": {
			"0x1111111111111111111111111111111111111111": {
				"option": [{"current_open_interest":"100","interest_cap":"1000","manager_currency":"USDC"}]
			}
		},
		"srm_im_discount": "0.5",
		"srm_mm_discount": "0.6",
		"srm_perp_margin_requirements": {"im_perp_req":"0.05","mm_perp_req":"0.025","max_leverage":"20"},
		"borrow_apy": "0.08",
		"supply_apy": "0.04",
		"total_borrow": "5000",
		"total_supply": "10000"
	}`)
	var c types.Currency
	require.NoError(t, json.Unmarshal(raw, &c))
	assert.Equal(t, "ETH", c.Currency)
	assert.Equal(t, "2500", c.SpotPrice.String())
	assert.Equal(t, "2450", c.SpotPrice24h.String())
	assert.Equal(t, []string{"option", "perp"}, c.InstrumentTypes)
	assert.Equal(t, "ALL", c.MarketType)
	require.Len(t, c.Managers, 1)
	assert.Equal(t, "PM2", c.Managers[0].MarginType)
	require.Len(t, c.PM2CollateralDiscounts, 1)
	assert.Equal(t, "USDC", c.PM2CollateralDiscounts[0].ManagerCurrency)
	require.NotNil(t, c.SRMPerpMarginRequirements)
	assert.Equal(t, "20", c.SRMPerpMarginRequirements.MaxLeverage.String())
	require.Contains(t, c.AssetCapAndSupplyPerManager, "0x1111111111111111111111111111111111111111")
	stats := c.AssetCapAndSupplyPerManager["0x1111111111111111111111111111111111111111"]["option"]
	require.Len(t, stats, 1)
	assert.Equal(t, "100", stats[0].CurrentOpenInterest.String())
}

func TestCurrency_NullableFieldsAbsent(t *testing.T) {
	// Smaller currency record — `spot_price_24h`,
	// `srm_perp_margin_requirements`, `erc20_details` are all absent
	// (or `null`) on the wire for currencies without 24h data, no SRM
	// perp market, or non-ERC-20 underlyings. All must decode to
	// Go zero-values without error.
	raw := []byte(`{
		"currency": "BTC",
		"spot_price": "65000",
		"spot_price_24h": null,
		"instrument_types": ["perp"],
		"market_type": "ALL",
		"managers": [],
		"pm2_collateral_discounts": [],
		"protocol_asset_addresses": {"option":null,"perp":null,"spot":null,"underlying_erc20":null},
		"asset_cap_and_supply_per_manager": {},
		"srm_im_discount": "0",
		"srm_mm_discount": "0",
		"srm_perp_margin_requirements": null,
		"erc20_details": null,
		"borrow_apy": "0",
		"supply_apy": "0",
		"total_borrow": "0",
		"total_supply": "0"
	}`)
	var c types.Currency
	require.NoError(t, json.Unmarshal(raw, &c))
	assert.Equal(t, "BTC", c.Currency)
	assert.Equal(t, "0", c.SpotPrice24h.String())
	assert.Nil(t, c.SRMPerpMarginRequirements)
	assert.Nil(t, c.ERC20Details)
	assert.True(t, c.ProtocolAssetAddresses.Option.IsZero())
}

func TestOptionSettlementPrice_Decode(t *testing.T) {
	raw := []byte(`{
		"expiry_date": "20260327",
		"utc_expiry_sec": 1774828800,
		"price": "65000.5"
	}`)
	var o types.OptionSettlementPrice
	require.NoError(t, json.Unmarshal(raw, &o))
	assert.Equal(t, "20260327", o.ExpiryDate)
	assert.Equal(t, int64(1774828800), o.UTCExpirySec)
	assert.Equal(t, "65000.5", o.Price.String())
}

func TestOptionSettlementPrice_NullPrice(t *testing.T) {
	// Pre-settlement: the OAS marks `price` nullable; null must
	// decode to the zero-value Decimal rather than fail.
	raw := []byte(`{
		"expiry_date": "20260327",
		"utc_expiry_sec": 1774828800,
		"price": null
	}`)
	var o types.OptionSettlementPrice
	require.NoError(t, json.Unmarshal(raw, &o))
	assert.Equal(t, "0", o.Price.String())
}

// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the response shapes returned by the public
// markets / metadata endpoints: `public/get_currency` and
// `public/get_option_settlement_prices`.
package types

// Currency is the response of `public/get_currency`. It carries the
// per-asset margin parameters, manager addresses, and protocol-asset
// addresses Derive's risk engine uses for one underlying currency.
//
// The shape mirrors `PublicGetCurrencyResultSchema` in Derive's v2.2
// OpenAPI spec.
type Currency struct {
	// Currency is the underlying currency symbol (e.g. "ETH", "BTC").
	Currency string `json:"currency"`
	// SpotPrice is the latest oracle spot price.
	SpotPrice Decimal `json:"spot_price"`
	// SpotPrice24h is the spot price 24 hours ago. Nullable on the
	// wire (zero-value when absent).
	SpotPrice24h Decimal `json:"spot_price_24h,omitempty"`

	// InstrumentTypes lists the instrument kinds Derive supports for
	// this currency. Wire enum: "erc20", "option", "perp".
	InstrumentTypes []string `json:"instrument_types"`
	// MarketType is the market category. Wire enum: "ALL",
	// "SRM_BASE_ONLY", "SRM_OPTION_ONLY", "SRM_PERP_ONLY", "CASH".
	MarketType string `json:"market_type"`

	// Managers is the list of margin-manager contracts that support
	// the currency.
	Managers []ManagerContract `json:"managers"`

	// PM2CollateralDiscounts is the IM/MM discount table for the
	// currency under Portfolio Margin 2.
	PM2CollateralDiscounts []PM2CollateralDiscount `json:"pm2_collateral_discounts"`

	// SRMIMDiscount is the Standard-Manager initial-margin discount
	// (a.k.a. LTV) — only the standard manager supports non-USDC
	// collateral.
	SRMIMDiscount Decimal `json:"srm_im_discount"`
	// SRMMMDiscount is the Standard-Manager maintenance-margin
	// discount (liquidation threshold).
	SRMMMDiscount Decimal `json:"srm_mm_discount"`
	// SRMPerpMarginRequirements is the standard-manager perp margin
	// table. Nullable on the wire — nil when the currency has no perp
	// market under SRM.
	SRMPerpMarginRequirements *SRMPerpMarginRequirements `json:"srm_perp_margin_requirements,omitempty"`

	// ProtocolAssetAddresses carries the Derive-protocol contract
	// addresses for the currency's option / perp / spot / underlying
	// ERC-20 contracts.
	ProtocolAssetAddresses ProtocolAssetAddresses `json:"protocol_asset_addresses"`

	// AssetCapAndSupplyPerManager carries the open-interest stats for
	// the currency keyed by manager address → asset type → list of
	// stats. Nested map shape mirrors the OAS verbatim.
	AssetCapAndSupplyPerManager map[string]map[string][]OpenInterestStats `json:"asset_cap_and_supply_per_manager"`

	// ERC20Details is a free-form map of ERC-20 metadata attached to
	// the currency (token address, decimals, name, etc). Nullable on
	// the wire — nil when the currency is not an ERC-20.
	ERC20Details map[string]string `json:"erc20_details,omitempty"`

	// BorrowAPY is the borrow APY (only populated for USDC).
	BorrowAPY Decimal `json:"borrow_apy"`
	// SupplyAPY is the supply APY (only populated for USDC).
	SupplyAPY Decimal `json:"supply_apy"`
	// TotalBorrow is the total collateral borrowed in the protocol
	// (only USDC is borrowable).
	TotalBorrow Decimal `json:"total_borrow"`
	// TotalSupply is the total collateral supplied in the protocol.
	TotalSupply Decimal `json:"total_supply"`
}

// ManagerContract is one margin-manager contract entry on a [Currency].
//
// Mirrors `ManagerContractResponseSchema`.
type ManagerContract struct {
	// Address is the manager contract address.
	Address Address `json:"address"`
	// Currency is the manager's currency. Only populated for
	// portfolio managers; nullable on the wire.
	Currency string `json:"currency,omitempty"`
	// MarginType is "PM", "SM", or "PM2" — the margin model the
	// manager implements.
	MarginType string `json:"margin_type"`
}

// PM2CollateralDiscount is one entry in [Currency.PM2CollateralDiscounts]
// — the IM and MM discount for the currency under Portfolio Margin 2,
// keyed by the manager's quote currency.
//
// Mirrors `PM2CollateralDiscountsSchema`.
type PM2CollateralDiscount struct {
	// ManagerCurrency is the quote currency of the manager (e.g. "USDC").
	ManagerCurrency string `json:"manager_currency"`
	// IMDiscount is the initial-margin discount.
	IMDiscount Decimal `json:"im_discount"`
	// MMDiscount is the maintenance-margin discount.
	MMDiscount Decimal `json:"mm_discount"`
}

// ProtocolAssetAddresses carries the Derive-protocol contract addresses
// for one currency's option / perp / spot / underlying ERC-20
// instruments.
//
// All four fields are nullable on the wire — empty Address means the
// currency does not support that asset class.
//
// Mirrors `ProtocolAssetAddressesSchema`.
type ProtocolAssetAddresses struct {
	// Option is the on-chain option contract address.
	Option Address `json:"option,omitempty"`
	// Perp is the on-chain perp contract address.
	Perp Address `json:"perp,omitempty"`
	// Spot is the on-chain spot contract address.
	Spot Address `json:"spot,omitempty"`
	// UnderlyingERC20 is the underlying ERC-20 token address on Derive
	// chain.
	UnderlyingERC20 Address `json:"underlying_erc20,omitempty"`
}

// SRMPerpMarginRequirements is the standard-manager perp margin
// schedule for one currency.
//
// Mirrors `SRMPerpMarginRequirementsPublicSchema`.
type SRMPerpMarginRequirements struct {
	// IMPerpReq is the initial margin requirement for perp positions
	// (fraction of notional).
	IMPerpReq Decimal `json:"im_perp_req"`
	// MMPerpReq is the maintenance margin requirement for perp
	// positions (fraction of notional).
	MMPerpReq Decimal `json:"mm_perp_req"`
	// MaxLeverage is `1 / im_perp_req` — the cap the standard manager
	// will enforce on opening leverage.
	MaxLeverage Decimal `json:"max_leverage"`
}

// OpenInterestStats is one entry in
// [Currency.AssetCapAndSupplyPerManager]. The wire shape reports the
// current open interest and the cap on it for one (manager, asset
// type) combination.
//
// Mirrors `OpenInterestStatsSchema`.
type OpenInterestStats struct {
	// CurrentOpenInterest is the engine's current open interest.
	CurrentOpenInterest Decimal `json:"current_open_interest"`
	// InterestCap is the configured cap.
	InterestCap Decimal `json:"interest_cap"`
	// ManagerCurrency is the manager's currency. Only populated for
	// portfolio managers; nullable on the wire.
	ManagerCurrency string `json:"manager_currency,omitempty"`
}

// OptionSettlementPrice is one entry in
// `public/get_option_settlement_prices.expiries`. Each entry is one
// expiry and its (eventual) settlement price.
//
// Mirrors `ExpiryResponseSchema`. `Price` is nullable on the wire — it
// stays at the zero value until the expiry settles on chain.
type OptionSettlementPrice struct {
	// ExpiryDate is the expiry in `YYYYMMDD` form.
	ExpiryDate string `json:"expiry_date"`
	// UTCExpirySec is the expiry as a Unix timestamp in seconds.
	UTCExpirySec int64 `json:"utc_expiry_sec"`
	// Price is the on-chain settlement price. Zero-value (Decimal "0")
	// until the expiry settles.
	Price Decimal `json:"price,omitempty"`
}

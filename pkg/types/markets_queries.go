// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the query DTOs for the public market-data endpoints
// (instruments, tickers, assets, currency, descendant tree, margin
// watch, settlement prices, statistics).
package types

import "github.com/amiwrpremium/go-derive/pkg/enums"

// InstrumentQuery parameterises public/get_instrument.
type InstrumentQuery struct {
	// Name is the canonical instrument name (e.g. "BTC-PERP").
	Name string
}

// InstrumentsQuery parameterises public/get_instruments — the active,
// currency-filtered listing UI endpoint.
type InstrumentsQuery struct {
	// Currency narrows the listing to one underlying currency.
	Currency string
	// Kind narrows the listing to one instrument type (perp, option,
	// erc20). Zero value lists every kind.
	Kind enums.InstrumentType
}

// AllInstrumentsQuery parameterises the filter portion of
// public/get_all_instruments. Pagination travels separately on
// methods that accept it.
type AllInstrumentsQuery struct {
	// Kind narrows the listing to one instrument type.
	Kind enums.InstrumentType
	// IncludeExpired opts into expired instruments. Set to true for
	// historical lookups; leave false for trading-side cache warming.
	IncludeExpired bool
}

// TickerQuery parameterises public/get_ticker.
type TickerQuery struct {
	// Name identifies the instrument.
	Name string
}

// TickersQuery parameterises public/get_tickers — the bulk
// ticker-by-instrument-type listing.
type TickersQuery struct {
	// InstrumentType is required; one of perp / option / erc20.
	InstrumentType enums.InstrumentType
	// Currency is required for option queries, ignored for perp /
	// erc20.
	Currency string
	// ExpiryDate is the YYYYMMDD numeric form. Required for option
	// queries, ignored for perp / erc20.
	ExpiryDate int64
}

// CurrencyQuery parameterises public/get_currency.
type CurrencyQuery struct {
	// Currency is the underlying symbol (e.g. "USDC", "BTC").
	Currency string
}

// AssetQuery parameterises public/get_asset.
type AssetQuery struct {
	// Name identifies the asset.
	Name string
}

// AssetsQuery parameterises public/get_assets.
type AssetsQuery struct {
	// AssetType filters by asset class. Zero value returns all.
	AssetType enums.AssetType
	// Currency narrows the listing to one currency.
	Currency string
	// Expired opts into expired assets.
	Expired bool
}

// OptionSettlementPricesQuery parameterises
// public/get_option_settlement_prices.
type OptionSettlementPricesQuery struct {
	// Currency selects the underlying currency.
	Currency string
}

// DescendantTreeQuery parameterises public/get_descendant_tree.
type DescendantTreeQuery struct {
	// WalletOrInviteCode is either an Ethereum wallet address or an
	// invite code — the endpoint accepts either form.
	WalletOrInviteCode string
}

// TransactionQuery parameterises public/get_transaction.
type TransactionQuery struct {
	// TransactionID is the server-side transaction id.
	TransactionID string
}

// MarginWatchQuery parameterises public/margin_watch.
type MarginWatchQuery struct {
	// SubaccountID is the subaccount whose maintenance margin is
	// being checked.
	SubaccountID int64
	// ForceOnchain forces an on-chain re-computation of the margin
	// snapshot rather than reading the engine's cached view.
	ForceOnchain bool
	// IsDelayedLiquidation flags the request as part of a
	// delayed-liquidation flow.
	IsDelayedLiquidation bool
}

// LatestSignedFeedsQuery parameterises public/get_latest_signed_feeds.
type LatestSignedFeedsQuery struct {
	// Currency selects the underlying currency.
	Currency string
	// Expiry is an optional expiry filter (zero means "all").
	Expiry int64
}

// PerpImpactTWAPQuery parameterises public/get_perp_impact_twap.
type PerpImpactTWAPQuery struct {
	// Currency selects the perpetual underlying.
	Currency string
	// StartTime is the inclusive window start (Unix seconds).
	StartTime int64
	// EndTime is the inclusive window end (Unix seconds).
	EndTime int64
}

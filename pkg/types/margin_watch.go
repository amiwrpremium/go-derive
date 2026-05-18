// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the per-event payload of Derive's `margin_watch`
// WebSocket channel: a stream of subaccounts whose maintenance margin
// has crossed the watch threshold and may be at imminent liquidation
// risk.
package types

import "github.com/amiwrpremium/go-derive/pkg/enums"

// MarginWatch is one event from the `margin.watch` channel. The
// channel emits a slice of these — one per subaccount that's near
// or below its maintenance-margin floor at the snapshot's
// `valuation_timestamp`.
//
// The shape mirrors the per-event payload documented at
// docs.derive.xyz/reference/margin-watch. The docs publish the full
// nested object (top-level margins + per-collateral and per-position
// breakdowns) — this type carries all of them. The Go
// [enums.MarginType] also models "PM2".
type MarginWatch struct {
	// SubaccountID identifies the at-risk subaccount.
	SubaccountID int64 `json:"subaccount_id"`
	// Currency is the subaccount's quote currency (e.g. "USDC").
	Currency string `json:"currency"`
	// MarginType is the margin model in use ("PM", "PM2", or "SM").
	MarginType enums.MarginType `json:"margin_type"`
	// SubaccountValue is the total mark-to-market value of all
	// positions and collaterals.
	SubaccountValue Decimal `json:"subaccount_value"`
	// InitialMargin is the total initial-margin requirement.
	InitialMargin Decimal `json:"initial_margin"`
	// MaintenanceMargin is the total maintenance-margin requirement.
	// If this falls below zero the subaccount is flagged for
	// liquidation.
	MaintenanceMargin Decimal `json:"maintenance_margin"`
	// ValuationTimestamp is when the engine computed the margin /
	// MtM (Unix seconds).
	ValuationTimestamp int64 `json:"valuation_timestamp"`
	// Collaterals is the per-collateral breakdown that feeds
	// [SubaccountValue], [InitialMargin], and [MaintenanceMargin].
	// Reuses the channel-side [MarginWatchCollateral] shape (a
	// subset of [Collateral] — only the fields the channel
	// publishes).
	Collaterals []MarginWatchCollateral `json:"collaterals"`
	// Positions is the per-position breakdown. Reuses the channel-
	// side [MarginWatchPosition] shape — Greeks plus margin
	// contribution and per-position liquidation price.
	Positions []MarginWatchPosition `json:"positions"`
}

// MarginWatchCollateral is one collateral row inside a [MarginWatch]
// event. The shape is a strict subset of [Collateral] — only the
// fields the `margin.watch` channel publishes.
type MarginWatchCollateral struct {
	// AssetName is the collateral's symbol (e.g. "USDC", "weETH").
	AssetName string `json:"asset_name"`
	// AssetType classifies the collateral ("erc20", "option", or
	// "perp").
	AssetType enums.AssetType `json:"asset_type"`
	// Amount is the held quantity in native units.
	Amount Decimal `json:"amount"`
	// MarkPrice is the current mark in quote units.
	MarkPrice Decimal `json:"mark_price"`
	// MarkValue is Amount × MarkPrice — the dollar-equivalent
	// contribution to subaccount value.
	MarkValue Decimal `json:"mark_value"`
	// Delta is this collateral's contribution to the subaccount's
	// total delta.
	Delta Decimal `json:"delta"`
	// DeltaCurrency is the currency Delta is denominated in.
	DeltaCurrency string `json:"delta_currency"`
	// InitialMargin is this collateral's contribution to the
	// subaccount's initial-margin requirement.
	InitialMargin Decimal `json:"initial_margin"`
	// MaintenanceMargin is this collateral's contribution to the
	// subaccount's maintenance-margin requirement.
	MaintenanceMargin Decimal `json:"maintenance_margin"`
}

// MarginWatchPosition is one position row inside a [MarginWatch]
// event. The shape is a strict subset of [Position] — only the
// fields the `margin.watch` channel publishes.
type MarginWatchPosition struct {
	// InstrumentName identifies the position's market.
	InstrumentName string `json:"instrument_name"`
	// InstrumentType is "perp", "option", or "erc20".
	InstrumentType enums.InstrumentType `json:"instrument_type"`
	// Amount is the signed position size (positive = long).
	Amount Decimal `json:"amount"`
	// Delta is the position's contribution to subaccount delta.
	// Always present (zero for non-option positions).
	Delta Decimal `json:"delta"`
	// Gamma is the option Greek; zero for perps and ERC-20.
	Gamma Decimal `json:"gamma"`
	// Theta is the option Greek; zero for perps and ERC-20.
	Theta Decimal `json:"theta"`
	// Vega is the option Greek; zero for perps and ERC-20.
	Vega Decimal `json:"vega"`
	// IndexPrice is the underlying index price (options only;
	// engine may emit zero for non-options).
	IndexPrice Decimal `json:"index_price"`
	// MarkPrice is the position's current mark.
	MarkPrice Decimal `json:"mark_price"`
	// MarkValue is Amount × MarkPrice.
	MarkValue Decimal `json:"mark_value"`
	// InitialMargin is this position's contribution to the
	// subaccount's initial-margin requirement.
	InitialMargin Decimal `json:"initial_margin"`
	// MaintenanceMargin is this position's contribution to the
	// subaccount's maintenance-margin requirement.
	MaintenanceMargin Decimal `json:"maintenance_margin"`
	// LiquidationPrice is the price at which the position would be
	// liquidated, or null when the engine cannot compute one (zero
	// here after JSON-null decode).
	LiquidationPrice Decimal `json:"liquidation_price,omitempty"`
}

// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// All numeric fields use [Decimal], a thin wrapper around shopspring/decimal,
// so price/size/fee values never lose precision through float64 round-trips.
// On the wire, [Decimal] reads and writes JSON strings (Derive's preferred
// representation); a fallback path also accepts JSON numbers for resilience.
//
// Identifier types ([Address], [TxHash], [MillisTime]) carry the same
// round-trip guarantees: each one preserves the canonical wire format
// regardless of how Go marshals the surrounding struct.
//
// # Why named types
//
// Plain string and int64 fields would parse just fine, but named types let
// the SDK enforce invariants at construction time (NewAddress checksum
// check, NewDecimal precision check) and let callers tell at a glance which
// values are amounts vs prices vs subaccount ids.
package types

import "github.com/amiwrpremium/go-derive/pkg/enums"

// Position is a held position in one instrument on a subaccount.
//
// Amount is signed: positive for long, negative for short, zero for flat.
// Most numeric fields are denominated in the quote currency (USDC for the
// vast majority of Derive markets); Amount itself is in base-currency units.
//
// Greeks (Delta/Gamma/Theta/Vega) are populated for option positions and
// zero for perp/erc20 positions.
type Position struct {
	// InstrumentName identifies the market.
	InstrumentName string `json:"instrument_name"`
	// InstrumentType identifies whether this is a perp, option or ERC20.
	InstrumentType enums.InstrumentType `json:"instrument_type"`
	// CreationTimestamp is when the position first appeared on the engine.
	CreationTimestamp MillisTime `json:"creation_timestamp,omitempty"`
	// Amount is the signed position size (positive=long, negative=short).
	Amount Decimal `json:"amount"`
	// AveragePrice is the volume-weighted entry price.
	AveragePrice Decimal `json:"average_price"`
	// MarkPrice is the engine's current mark.
	MarkPrice Decimal `json:"mark_price"`
	// MarkValue is the position's mark-to-market value in quote currency.
	MarkValue Decimal `json:"mark_value"`
	// IndexPrice is the underlying index price (zero if not yet computed).
	IndexPrice Decimal `json:"index_price,omitempty"`
	// Leverage is the position's effective leverage.
	Leverage Decimal `json:"leverage,omitempty"`
	// LiquidationPrice is the price at which the engine would liquidate
	// (zero if no liquidation risk).
	LiquidationPrice Decimal `json:"liquidation_price,omitempty"`

	// InitialMargin is the engine's initial-margin requirement for this
	// position alone.
	InitialMargin Decimal `json:"initial_margin,omitempty"`
	// MaintenanceMargin is the maintenance-margin requirement.
	MaintenanceMargin Decimal `json:"maintenance_margin,omitempty"`
	// OpenOrdersMargin is the margin reserved for open orders against this position.
	OpenOrdersMargin Decimal `json:"open_orders_margin,omitempty"`

	// UnrealizedPNL is the mark-to-market PnL.
	UnrealizedPNL Decimal `json:"unrealized_pnl"`
	// RealizedPNL is the cumulative realized PnL across closes.
	RealizedPNL Decimal `json:"realized_pnl"`

	// CumulativeFunding is the total funding paid/received over the
	// position's lifetime (perps only).
	CumulativeFunding Decimal `json:"cumulative_funding,omitempty"`
	// PendingFunding is funding accrued since the last settlement.
	PendingFunding Decimal `json:"pending_funding,omitempty"`
	// NetSettlements is the cumulative net of perp / option settlements
	// applied to the position.
	NetSettlements Decimal `json:"net_settlements,omitempty"`

	// Delta, Gamma, Theta, Vega are the option greeks (option positions
	// only; zero for perp / erc20).
	Delta Decimal `json:"delta,omitempty"`
	Gamma Decimal `json:"gamma,omitempty"`
	Theta Decimal `json:"theta,omitempty"`
	Vega  Decimal `json:"vega,omitempty"`
}

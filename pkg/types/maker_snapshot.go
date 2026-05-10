// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the per-quote maker-snapshot row returned by
// `public/get_detailed_maker_snapshot_history`.
package types

import "encoding/json"

// MakerSnapshot is one row of the maker-program quoting history —
// captured at a single timestamp on a single instrument, scoring how
// the maker's quote contributes to the program's coverage / quality
// metrics.
//
// Mirrors the per-row shape per
// docs.derive.xyz/reference/public-get_detailed_maker_snapshot_history.
type MakerSnapshot struct {
	// Wallet is the market maker's wallet address.
	Wallet string `json:"wallet"`
	// AssetType identifies the instrument's product class (e.g.
	// "perp", "option").
	AssetType string `json:"asset_type"`
	// Currency is the underlying currency (e.g. "ETH").
	Currency string `json:"currency"`
	// InstrumentName identifies the market.
	InstrumentName string `json:"instrument_name"`
	// IsBid is true for bid-side quotes, false for ask-side.
	IsBid bool `json:"is_bid"`
	// Timestamp is when the snapshot was captured (millisecond
	// Unix epoch).
	Timestamp MillisTime `json:"timestamp"`

	// MidPrice is the engine's mid price at the snapshot.
	MidPrice Decimal `json:"mid_price"`
	// BestPrice is the maker's best resting price.
	BestPrice Decimal `json:"best_price"`
	// IndexPrice is the underlying index price.
	IndexPrice Decimal `json:"index_price"`

	// Notional is the maker's quoted notional.
	Notional Decimal `json:"notional"`
	// ScaledNotional is Notional after the program's instrument /
	// BBO factors are applied.
	ScaledNotional Decimal `json:"scaled_notional"`
	// TotalScaledNotional is the maker's running total of
	// ScaledNotional for the epoch.
	TotalScaledNotional Decimal `json:"total_scaled_notional"`
	// DeductedNotional is the notional component deducted by the
	// program's risk filters.
	DeductedNotional Decimal `json:"deducted_notional,omitempty"`

	// BBOFactor is the BBO-tightness factor applied to Notional.
	BBOFactor Decimal `json:"bbo_factor,omitempty"`
	// InstrumentFactor is the per-instrument weighting factor.
	InstrumentFactor Decimal `json:"instrument_factor,omitempty"`
	// DeductedFactor is the deduction factor applied for risk
	// filters.
	DeductedFactor Decimal `json:"deducted_factor,omitempty"`

	// CoverageScore is the snapshot's contribution to the maker's
	// coverage component.
	CoverageScore Decimal `json:"coverage_score,omitempty"`
	// QualityScore is the snapshot's contribution to the maker's
	// quality component.
	QualityScore Decimal `json:"quality_score,omitempty"`

	// Quotes is the per-quote backing data. Preserved as raw JSON
	// because the inner shape isn't documented as a fixed schema —
	// decode further at the call site.
	Quotes json.RawMessage `json:"quotes,omitempty"`
}

// DetailedMakerSnapshotHistory is the response of
// `public/get_detailed_maker_snapshot_history`. It pairs the program
// metadata with a paginated slice of per-quote snapshots.
type DetailedMakerSnapshotHistory struct {
	// Program is the program metadata (same shape as
	// `public/get_maker_programs` entries).
	Program MakerProgram `json:"program"`
	// Snapshots is the per-quote snapshot rows for the requested
	// page.
	Snapshots []MakerSnapshot `json:"snapshots"`
	// Pagination carries the totals for the query.
	Pagination Page `json:"pagination"`
}

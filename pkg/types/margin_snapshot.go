// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the per-subaccount margin snapshot returned by
// `public/margin_watch` (an RPC, distinct from the
// platform-wide `margin_watch` WebSocket channel).
package types

// MarginSnapshot is the response of `public/margin_watch` — a
// schema-faithful margin snapshot for one subaccount including its
// collaterals and positions.
//
// MarginType on this RPC matches the documented two-value enum on
// the wire ("PM" / "SM"); the Go `enums.MarginType` includes "PM2"
// too, which is fine since defined-string enums tolerate extra
// values.
type MarginSnapshot struct {
	// SubaccountID is the queried subaccount.
	SubaccountID int64 `json:"subaccount_id"`
	// Currency is the subaccount's quote currency (e.g. "USDC").
	Currency string `json:"currency"`
	// MarginType is the margin model in use ("PM" or "SM").
	MarginType string `json:"margin_type"`
	// SubaccountValue is the total mark-to-market equity.
	SubaccountValue Decimal `json:"subaccount_value"`
	// InitialMargin is the IM requirement.
	InitialMargin Decimal `json:"initial_margin"`
	// MaintenanceMargin is the MM requirement; falling below zero
	// flags the subaccount for liquidation.
	MaintenanceMargin Decimal `json:"maintenance_margin"`
	// ValuationTimestamp is when the engine computed the snapshot
	// (Unix seconds).
	ValuationTimestamp int64 `json:"valuation_timestamp"`
	// Collaterals is the per-asset collateral breakdown.
	Collaterals []Collateral `json:"collaterals"`
	// Positions is the held positions.
	Positions []Position `json:"positions"`
}

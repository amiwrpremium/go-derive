// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the per-event payload of Derive's `margin_watch`
// WebSocket channel: a stream of subaccounts whose maintenance margin
// has crossed the watch threshold and may be at imminent liquidation
// risk.
package types

import "github.com/amiwrpremium/go-derive/pkg/enums"

// MarginWatch is one event from the `margin_watch` channel. The
// channel emits a slice of these — one per subaccount that's near
// or below its maintenance-margin floor at the snapshot's
// `valuation_timestamp`.
//
// The shape mirrors the per-event payload documented at
// docs.derive.xyz/reference/margin-watch (the channel page —
// distinct from the public/margin_watch RPC, which returns a fuller
// per-subaccount shape covered by [MarginSnapshot]). MarginType on
// this channel is restricted to "PM" or "SM"; the Go
// [enums.MarginType] also includes "PM2" but margin_watch will not
// emit it.
type MarginWatch struct {
	// SubaccountID identifies the at-risk subaccount.
	SubaccountID int64 `json:"subaccount_id"`
	// Currency is the subaccount's quote currency (e.g. "USDC").
	Currency string `json:"currency"`
	// MarginType is the margin model in use ("PM" or "SM").
	MarginType enums.MarginType `json:"margin_type"`
	// SubaccountValue is the total mark-to-market value of all
	// positions and collaterals.
	SubaccountValue Decimal `json:"subaccount_value"`
	// MaintenanceMargin is the total maintenance-margin requirement.
	// If this falls below zero the subaccount is flagged for
	// liquidation.
	MaintenanceMargin Decimal `json:"maintenance_margin"`
	// ValuationTimestamp is when the engine computed the margin /
	// MtM (Unix seconds).
	ValuationTimestamp int64 `json:"valuation_timestamp"`
}

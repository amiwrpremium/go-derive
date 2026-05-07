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

// SubAccount is a snapshot of one subaccount as returned by GetSubaccount.
//
// A wallet has one or more subaccounts, each with its own positions,
// collateral and margin state. Subaccounts isolate risk: a liquidation in
// one does not cascade into another.
type SubAccount struct {
	// SubaccountID is the unique numeric id.
	SubaccountID int64 `json:"subaccount_id"`
	// OwnerAddress is the smart-account owner that controls this subaccount.
	OwnerAddress Address `json:"owner_address"`
	// MarginType is "PM" (portfolio margin), "SM" (standard margin), etc.
	MarginType string `json:"margin_type"`
	// IsUnderLiquidation is true when the engine is actively liquidating
	// the subaccount.
	IsUnderLiquidation bool `json:"is_under_liquidation"`
	// SubaccountValue is the total equity.
	SubaccountValue Decimal `json:"subaccount_value"`
	// InitialMargin is the margin required to open new orders.
	InitialMargin Decimal `json:"initial_margin"`
	// MaintenanceMargin is the liquidation floor.
	MaintenanceMargin Decimal `json:"maintenance_margin"`
	// OpenOrders is the list of currently-open orders.
	OpenOrders []Order `json:"open_orders,omitempty"`
	// Positions is the list of open positions.
	Positions []Position `json:"positions,omitempty"`
	// Collaterals is the per-asset collateral breakdown.
	Collaterals []Collateral `json:"collaterals,omitempty"`
}

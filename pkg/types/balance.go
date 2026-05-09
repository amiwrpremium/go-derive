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

import "github.com/amiwrpremium/go-derive"

// Collateral is one collateral asset balance for a subaccount.
//
// Each subaccount can hold multiple collaterals; PMRM (portfolio-margin
// risk-managed) subaccounts are restricted to USDC.
type Collateral struct {
	// AssetName is the human-readable symbol (e.g. "USDC", "weETH").
	AssetName string `json:"asset_name"`
	// AssetType identifies the asset class — see [derive.AssetType].
	AssetType derive.AssetType `json:"asset_type"`
	// Currency is the underlying currency (e.g. "USDC", "ETH").
	Currency string `json:"currency,omitempty"`
	// Amount is the balance in the asset's native units.
	Amount Decimal `json:"amount"`
	// MarkPrice is the asset's mark price in quote currency (USDC).
	MarkPrice Decimal `json:"mark_price,omitempty"`
	// MarkValue is the asset balance valued at the current mark.
	MarkValue Decimal `json:"mark_value"`
	// CumulativeInterest is the lifetime interest earned/paid on this asset.
	CumulativeInterest Decimal `json:"cumulative_interest,omitempty"`
	// PendingInterest is interest accrued but not yet settled.
	PendingInterest Decimal `json:"pending_interest,omitempty"`
	// InitialMargin is the asset's contribution to the subaccount's IM.
	InitialMargin Decimal `json:"initial_margin,omitempty"`
	// MaintenanceMargin is the asset's contribution to the subaccount's MM.
	MaintenanceMargin Decimal `json:"maintenance_margin,omitempty"`
}

// Balance summarises a subaccount's value and margin posture in one struct.
//
// SubaccountValue is the headline equity number; InitialMargin and
// MaintenanceMargin set the bands inside which open orders are accepted
// and outside which the engine liquidates.
type Balance struct {
	// SubaccountID identifies the subaccount this balance belongs to.
	SubaccountID int64 `json:"subaccount_id"`
	// SubaccountValue is the total equity (collateral + unrealized PnL +
	// pending funding).
	SubaccountValue Decimal `json:"subaccount_value"`
	// InitialMargin is the margin required to open new orders.
	InitialMargin Decimal `json:"initial_margin"`
	// MaintenanceMargin is the margin floor; breaching it triggers
	// liquidation.
	MaintenanceMargin Decimal `json:"maintenance_margin"`
	// Collaterals is the per-asset balance breakdown.
	Collaterals []Collateral `json:"collaterals"`
	// Positions is the open positions by instrument (omitted by some endpoints).
	Positions []Position `json:"positions,omitempty"`
}

// BalanceUpdate is one entry on the `subaccount.{id}.balances` subscription
// channel. Where Balance is a snapshot, BalanceUpdate is a delta event:
// it carries the [derive.BalanceUpdateType] explaining what caused the
// change (a fill, a deposit, an interest accrual, etc.).
type BalanceUpdate struct {
	// SubaccountID identifies the subaccount this update belongs to.
	SubaccountID int64 `json:"subaccount_id"`
	// AssetName is the affected asset.
	AssetName string `json:"asset_name,omitempty"`
	// AssetType identifies the asset class.
	AssetType derive.AssetType `json:"asset_type,omitempty"`
	// Amount is the new balance after the update.
	Amount Decimal `json:"amount,omitempty"`
	// PreviousAmount is the balance before the update.
	PreviousAmount Decimal `json:"previous_amount,omitempty"`
	// Delta is the signed change.
	Delta Decimal `json:"delta,omitempty"`
	// UpdateType classifies the cause of the update — see
	// [derive.BalanceUpdateType].
	UpdateType derive.BalanceUpdateType `json:"update_type,omitempty"`
	// TxHash is the on-chain transaction hash that generated the update,
	// for update types that involve on-chain settlement.
	TxHash TxHash `json:"tx_hash,omitempty"`
	// TxStatus is the on-chain settlement state.
	TxStatus derive.TxStatus `json:"tx_status,omitempty"`
	// Timestamp is when the update was recorded.
	Timestamp MillisTime `json:"timestamp,omitempty"`
}

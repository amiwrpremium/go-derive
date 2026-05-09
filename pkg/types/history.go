// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the per-event types returned by Derive's history
// endpoints: funding payments, interest accruals, ERC-20 transfers,
// liquidation auctions and bids, option settlements, subaccount value
// snapshots, and the expired-and-cancelled archive bundle.
package types

// FundingPayment is one entry in `private/get_funding_history`. Funding
// is debited or credited at every funding tick on each open perp; the
// sum across an interval is the perp PnL component on funding.
type FundingPayment struct {
	// InstrumentName is the perpetual market the payment applies to.
	InstrumentName string `json:"instrument_name"`
	// SubaccountID is the subaccount the payment was attributed to.
	SubaccountID int64 `json:"subaccount_id"`
	// Timestamp is the funding tick (millisecond Unix epoch).
	Timestamp MillisTime `json:"timestamp"`
	// Funding is the signed funding amount in quote currency. Positive
	// means the subaccount received funding; negative means paid.
	Funding Decimal `json:"funding"`
	// PnL is the realized PnL contribution from the funding tick.
	PnL Decimal `json:"pnl"`
}

// InterestPayment is one entry in `private/get_interest_history`.
// Interest accrues on borrowed/lent collateral at Derive's per-asset
// rates.
type InterestPayment struct {
	// SubaccountID is the subaccount the interest was applied to.
	SubaccountID int64 `json:"subaccount_id"`
	// Timestamp is when the interest was applied (millisecond Unix
	// epoch).
	Timestamp MillisTime `json:"timestamp"`
	// Interest is the signed interest amount. Positive is income,
	// negative is expense.
	Interest Decimal `json:"interest"`
}

// ERC20Transfer is one entry in `private/get_erc20_transfer_history`.
// Each ERC-20 deposit or withdrawal between two subaccounts on the
// wallet shows up as a record here.
type ERC20Transfer struct {
	// SubaccountID is the subaccount the transfer was attributed to.
	SubaccountID int64 `json:"subaccount_id"`
	// CounterpartySubaccountID is the other side of the transfer.
	CounterpartySubaccountID int64 `json:"counterparty_subaccount_id"`
	// Asset is the asset symbol (e.g. "USDC").
	Asset string `json:"asset"`
	// Amount is the transferred quantity.
	Amount Decimal `json:"amount"`
	// IsOutgoing reports the direction relative to SubaccountID.
	IsOutgoing bool `json:"is_outgoing"`
	// Timestamp is when the transfer was observed (millisecond Unix
	// epoch).
	Timestamp MillisTime `json:"timestamp"`
	// TxHash is the on-chain transfer transaction hash.
	TxHash TxHash `json:"tx_hash"`
}

// OptionSettlement is one entry in
// `private/get_option_settlement_history` and
// `public/get_option_settlement_history`. Each record describes one
// option contract settling at expiry against the marked settlement
// price.
type OptionSettlement struct {
	// SubaccountID is the subaccount whose position settled.
	SubaccountID int64 `json:"subaccount_id"`
	// InstrumentName is the option contract name.
	InstrumentName string `json:"instrument_name"`
	// Expiry is the option expiry (Unix seconds).
	Expiry int64 `json:"expiry"`
	// Amount is the settled position size (signed).
	Amount Decimal `json:"amount"`
	// SettlementPrice is the index price at expiry.
	SettlementPrice Decimal `json:"settlement_price"`
	// OptionSettlementPnL is the PnL credited at settlement.
	OptionSettlementPnL Decimal `json:"option_settlement_pnl"`
	// OptionSettlementPnLExclFees is OptionSettlementPnL with the
	// fee component of cost basis excluded.
	OptionSettlementPnLExclFees Decimal `json:"option_settlement_pnl_excl_fees"`
}

// SubaccountValueRecord is one snapshot in
// `private/get_subaccount_value_history`. The series sampled at the
// requested period gives a portable equity curve for the subaccount.
type SubaccountValueRecord struct {
	// Timestamp is the sample point (millisecond Unix epoch).
	Timestamp MillisTime `json:"timestamp"`
	// Currency is the quote currency the values are denominated in.
	Currency string `json:"currency"`
	// MarginType is "PM", "PM2", or "SM" depending on the
	// subaccount's margin regime.
	MarginType string `json:"margin_type"`
	// SubaccountValue is the total portfolio value at the sample
	// point.
	SubaccountValue Decimal `json:"subaccount_value"`
	// InitialMargin is the initial-margin requirement at the sample
	// point.
	InitialMargin Decimal `json:"initial_margin"`
	// MaintenanceMargin is the maintenance-margin requirement at the
	// sample point.
	MaintenanceMargin Decimal `json:"maintenance_margin"`
	// DelayedMaintenanceMargin is the delayed-maintenance-margin
	// requirement at the sample point.
	DelayedMaintenanceMargin Decimal `json:"delayed_maintenance_margin"`
}

// ExpiredAndCancelledExport is the response of
// `private/expired_and_cancelled_history`. The endpoint produces
// pre-signed S3 URLs (one per archive shard) the caller can download
// directly to retrieve the bulk export of expired and cancelled
// orders.
type ExpiredAndCancelledExport struct {
	// PresignedURLs is the list of S3 download URLs.
	PresignedURLs []string `json:"presigned_urls"`
}

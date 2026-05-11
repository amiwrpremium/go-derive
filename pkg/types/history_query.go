// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds query DTOs for the time-windowed history methods.
// Endpoints share two common axes: an optional `[StartTimestamp,
// EndTimestamp]` window in milliseconds since the Unix epoch
// (defaults: 0 / current time) and an optional Wallet override
// (when non-empty the query spans every subaccount under that
// wallet; when empty the client-configured subaccount is used).
// Pagination is carried separately by [PageRequest] where supported.
package types

// HistoryWindow narrows a history query to a `[StartTimestamp,
// EndTimestamp]` window in milliseconds since the Unix epoch.
// Either side can be the zero value to defer to the server-side
// default (0 / current time).
type HistoryWindow struct {
	StartTimestamp MillisTime
	EndTimestamp   MillisTime
}

// FundingHistoryQuery narrows a paginated
// `private/get_funding_history` request. InstrumentName, when
// non-empty, restricts the result to one perpetual. Wallet, when
// non-empty, overrides the configured subaccount and spans every
// subaccount under that wallet.
type FundingHistoryQuery struct {
	HistoryWindow
	InstrumentName string
	Wallet         string
}

// InterestHistoryQuery narrows a `private/get_interest_history`
// request. Wallet, when non-empty, overrides the configured
// subaccount.
type InterestHistoryQuery struct {
	HistoryWindow
	Wallet string
}

// LiquidationHistoryQuery narrows a `private/get_liquidation_history`
// request. Wallet, when non-empty, takes precedence over
// SubaccountID per the engine's contract.
type LiquidationHistoryQuery struct {
	HistoryWindow
	Wallet string
}

// LiquidatorHistoryQuery narrows a paginated
// `private/get_liquidator_history` request to a time window.
type LiquidatorHistoryQuery struct {
	HistoryWindow
}

// OptionSettlementHistoryQuery narrows the option-settlement
// listings (both `private/get_option_settlement_history` and
// `public/get_option_settlement_history`). Wallet, when non-empty,
// queries across every subaccount under that wallet; the engine
// ignores SubaccountID when Wallet is set.
//
// Neither variant accepts timestamps — the query is purely by
// account identity.
type OptionSettlementHistoryQuery struct {
	Wallet string
}

// ERC20TransferHistoryQuery narrows a
// `private/get_erc20_transfer_history` request. Wallet, when
// non-empty, overrides SubaccountID.
type ERC20TransferHistoryQuery struct {
	HistoryWindow
	Wallet string
}

// SubaccountValueHistoryQuery selects the period bucket and time
// window for `private/get_subaccount_value_history`. PeriodSec
// accepts one of 900 (15m), 3600 (1h), 86400 (1d), 604800 (1w).
type SubaccountValueHistoryQuery struct {
	HistoryWindow
	PeriodSec int64
}

// Validate performs schema-level checks on the receiver. Returns
// nil on success or a wrapped [ErrInvalidParams].
func (q SubaccountValueHistoryQuery) Validate() error {
	switch q.PeriodSec {
	case 900, 3600, 86400, 604800:
		return nil
	}
	return invalidParam("period", "must be one of 900, 3600, 86400, 604800")
}

// ExpiredAndCancelledHistoryInput parameterises
// `private/expired_and_cancelled_history`, which triggers an
// archive export of expired and cancelled orders.
//
// The engine requires both Wallet and a subaccount filter; when
// SubaccountID is zero the SDK substitutes the client-configured
// subaccount. ExpirySec caps how long the presigned download URLs
// remain accessible — at most 604800 seconds (one week).
type ExpiredAndCancelledHistoryInput struct {
	HistoryWindow
	// Wallet is the wallet address the export is scoped to.
	Wallet string
	// SubaccountID restricts the export to one subaccount; zero
	// defaults to the client-configured subaccount.
	SubaccountID int64
	// ExpirySec is the lifetime (in seconds) of the returned
	// presigned URLs. Max 604800.
	ExpirySec int64
}

// Validate performs schema-level checks on the receiver.
func (in ExpiredAndCancelledHistoryInput) Validate() error {
	if in.Wallet == "" {
		return invalidParam("wallet", "required")
	}
	if in.ExpirySec <= 0 || in.ExpirySec > 604800 {
		return invalidParam("expiry", "must be in (0, 604800]")
	}
	return nil
}

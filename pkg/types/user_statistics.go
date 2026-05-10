// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the per-wallet trading statistics shape returned by
// `public/user_statistics` (single-wallet) and
// `public/all_user_statistics` (slice).
package types

// UserStatistics is one wallet's rolling trading statistics. Returned
// by `public/user_statistics` (with Wallet supplied as the request
// param, omitted from the response) and as one element of the
// `public/all_user_statistics` slice (Wallet populated).
//
// The shape mirrors the response per
// docs.derive.xyz/reference/public-user_statistics and -all_user_statistics.
type UserStatistics struct {
	// Wallet identifies the trader. Empty for the single-wallet
	// variant (since the wallet is the request param); populated
	// when the record is an entry in the all-user list.
	Wallet string `json:"wallet,omitempty"`
	// TotalBaseFee is the total $ base-fee component paid.
	TotalBaseFee Decimal `json:"total_base_fee"`
	// TotalContractFee is the total $ contract-fee component paid.
	TotalContractFee Decimal `json:"total_contract_fee"`
	// TotalFees is the sum of all fee components paid.
	TotalFees Decimal `json:"total_fees"`
	// TotalNotionalVolume is the cumulative notional volume.
	TotalNotionalVolume Decimal `json:"total_notional_volume"`
	// TotalPremiumVolume is the cumulative premium volume.
	TotalPremiumVolume Decimal `json:"total_premium_volume"`
	// TotalRegularBaseFee is the regular-rate base-fee component
	// (i.e. before maker rebates / referral discounts).
	TotalRegularBaseFee Decimal `json:"total_regular_base_fee"`
	// TotalRegularContractFee is the regular-rate contract-fee
	// component.
	TotalRegularContractFee Decimal `json:"total_regular_contract_fee"`
	// TotalTrades is the cumulative trade count.
	TotalTrades int64 `json:"total_trades"`
	// FirstTradeTimestamp is the wallet's first-ever trade time
	// (millisecond Unix epoch).
	FirstTradeTimestamp MillisTime `json:"first_trade_timestamp"`
	// LastTradeTimestamp is the wallet's most-recent trade time.
	LastTradeTimestamp MillisTime `json:"last_trade_timestamp"`
}

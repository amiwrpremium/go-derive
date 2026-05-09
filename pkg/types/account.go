// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds account-level shapes: the response of
// `private/get_account`, the per-fee discount table, and the response
// of `private/get_margin` and `public/get_margin` (both share the same
// shape).
package types

import "encoding/json"

// AccountResult is the response of `private/get_account`. It carries
// wallet-level metadata (subaccount ids, RFQ-maker status, referral
// code) and the rate-limit and fee schedules attached to the wallet.
//
// The shape mirrors `AccountResultSchema` in Derive's v2.2 OpenAPI
// spec.
type AccountResult struct {
	// SubaccountIDs is every subaccount registered under the wallet.
	SubaccountIDs []int64 `json:"subaccount_ids"`
	// Wallet is the smart-account owner address (lower-cased hex).
	Wallet string `json:"wallet"`
	// CancelOnDisconnect reports whether the kill-switch is currently
	// armed for this wallet. See `private/set_cancel_on_disconnect`.
	CancelOnDisconnect bool `json:"cancel_on_disconnect"`
	// CreationTimestampSec is the wallet creation time in Unix seconds.
	CreationTimestampSec int64 `json:"creation_timestamp_sec"`
	// IsRFQMaker reports whether the wallet is whitelisted to respond
	// to RFQs.
	IsRFQMaker bool `json:"is_rfq_maker"`
	// ReferralCode is the wallet's referral code (empty if none).
	ReferralCode string `json:"referral_code"`
	// WebSocketMatchingTPS is the per-second matching-engine call
	// budget for this wallet over WS.
	WebSocketMatchingTPS int64 `json:"websocket_matching_tps"`
	// WebSocketNonMatchingTPS is the per-second non-matching call
	// budget over WS.
	WebSocketNonMatchingTPS int64 `json:"websocket_non_matching_tps"`
	// WebSocketOptionTPS is the option-specific TPS budget over WS.
	WebSocketOptionTPS int64 `json:"websocket_option_tps"`
	// WebSocketPerpTPS is the perp-specific TPS budget over WS.
	WebSocketPerpTPS int64 `json:"websocket_perp_tps"`
	// PerEndpointTPS is a free-form map of per-endpoint TPS overrides
	// keyed by JSON-RPC method name. The OAS declares it as an open
	// `object`, so the SDK keeps it as a raw payload — decode at the
	// call site if you need it.
	PerEndpointTPS json.RawMessage `json:"per_endpoint_tps"`
	// FeeInfo is the per-product fee schedule.
	FeeInfo FeeInfo `json:"fee_info"`
}

// FeeInfo is the wallet's fee schedule. All values are decimal
// fractions on the wire; e.g. `"0.0003"` for 3 bps.
type FeeInfo struct {
	// BaseFeeDiscount is the wallet's flat fee discount (decimal).
	BaseFeeDiscount Decimal `json:"base_fee_discount"`
	// OptionMakerFee is the option maker fee (decimal fraction).
	OptionMakerFee Decimal `json:"option_maker_fee"`
	// OptionTakerFee is the option taker fee (decimal fraction).
	OptionTakerFee Decimal `json:"option_taker_fee"`
	// PerpMakerFee is the perpetual maker fee (decimal fraction).
	PerpMakerFee Decimal `json:"perp_maker_fee"`
	// PerpTakerFee is the perpetual taker fee (decimal fraction).
	PerpTakerFee Decimal `json:"perp_taker_fee"`
	// RFQMakerDiscount is the RFQ maker discount (decimal fraction).
	RFQMakerDiscount Decimal `json:"rfq_maker_discount"`
	// RFQTakerDiscount is the RFQ taker discount (decimal fraction).
	RFQTakerDiscount Decimal `json:"rfq_taker_discount"`
	// SpotMakerFee is the spot maker fee (decimal fraction).
	SpotMakerFee Decimal `json:"spot_maker_fee"`
	// SpotTakerFee is the spot taker fee (decimal fraction).
	SpotTakerFee Decimal `json:"spot_taker_fee"`
}

// MarginResult is the response of `private/get_margin` and
// `public/get_margin`. Both endpoints simulate a margin calculation
// against a (possibly modified) subaccount and report the pre/post
// initial- and maintenance-margin values.
//
// The shape mirrors `PrivateGetMarginResultSchema` /
// `PublicGetMarginResultSchema` in Derive's v2.2 OpenAPI spec.
type MarginResult struct {
	// SubaccountID is the subaccount the calculation ran against.
	SubaccountID int64 `json:"subaccount_id"`
	// IsValidTrade is true when the simulated trade leaves the
	// subaccount above the margin floor.
	IsValidTrade bool `json:"is_valid_trade"`
	// PreInitialMargin is the initial-margin requirement before the
	// simulated changes are applied.
	PreInitialMargin Decimal `json:"pre_initial_margin"`
	// PostInitialMargin is the requirement after.
	PostInitialMargin Decimal `json:"post_initial_margin"`
	// PreMaintenanceMargin is the maintenance-margin requirement
	// before the simulated changes are applied.
	PreMaintenanceMargin Decimal `json:"pre_maintenance_margin"`
	// PostMaintenanceMargin is the requirement after.
	PostMaintenanceMargin Decimal `json:"post_maintenance_margin"`
}

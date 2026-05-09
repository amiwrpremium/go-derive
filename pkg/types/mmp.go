// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the response shape of `private/get_mmp_config`.
// The matching input type for `private/set_mmp_config` lives in the
// internal/methods package because it's a method-side concern.
package types

// MMPConfigResult is one entry in the response of
// `private/get_mmp_config`. The endpoint returns one record per
// (subaccount, currency) pair the wallet has configured.
//
// Mirrors `MMPConfigResultSchema` in Derive's v2.2 OpenAPI spec.
type MMPConfigResult struct {
	// SubaccountID is the subaccount the rule applies to.
	SubaccountID int64 `json:"subaccount_id"`
	// Currency is the rule's underlying currency (e.g. "BTC",
	// "ETH"). Each subaccount can have at most one rule per
	// currency.
	Currency string `json:"currency"`
	// MMPFrozenTime is how long (ms) the subaccount stays frozen
	// after an MMP trigger. Zero means freeze until manual reset.
	MMPFrozenTime int64 `json:"mmp_frozen_time"`
	// MMPInterval is the rolling window (ms) the limits are
	// evaluated over. Zero disables MMP for this rule.
	MMPInterval int64 `json:"mmp_interval"`
	// MMPAmountLimit is the maximum total order amount (gross,
	// no netting) that can trade within MMPInterval before the
	// rule trips. Zero means no amount limit.
	MMPAmountLimit Decimal `json:"mmp_amount_limit,omitempty"`
	// MMPDeltaLimit is the maximum total signed delta (netted)
	// that can trade within MMPInterval before the rule trips.
	// Zero means no delta limit.
	MMPDeltaLimit Decimal `json:"mmp_delta_limit,omitempty"`
	// MMPUnfreezeTime is when the subaccount will be unfrozen
	// (millisecond Unix epoch). Zero when not currently frozen.
	MMPUnfreezeTime int64 `json:"mmp_unfreeze_time"`
	// IsFrozen reports whether the subaccount is currently
	// frozen for this rule.
	IsFrozen bool `json:"is_frozen"`
}

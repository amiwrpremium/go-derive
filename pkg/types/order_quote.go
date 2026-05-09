// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the response shape of `private/order_quote`.
package types

// OrderQuoteResult is the response of `private/order_quote`. The
// endpoint runs a hypothetical order through the matching engine
// without submitting and reports the engine's estimates for fill
// price, fee, and post-trade margin balance — useful for
// pre-flighting orders against thin books.
//
// The shape mirrors `PrivateOrderQuoteResultSchema` in
// `derivexyz/cockpit/orderbook-types`.
type OrderQuoteResult struct {
	// IsValid reports whether the order is expected to clear
	// margin requirements.
	IsValid bool `json:"is_valid"`
	// InvalidReason carries a human-readable reason when IsValid
	// is false. Empty when valid; the wire field is nullable.
	InvalidReason string `json:"invalid_reason,omitempty"`
	// EstimatedFillAmount is the amount the engine projects will
	// be crossed instantly on submission.
	EstimatedFillAmount Decimal `json:"estimated_fill_amount"`
	// EstimatedFillPrice is the engine's projected average fill
	// price.
	EstimatedFillPrice Decimal `json:"estimated_fill_price"`
	// EstimatedFee is the projected fee for the trade ($, whole
	// order).
	EstimatedFee Decimal `json:"estimated_fee"`
	// EstimatedRealizedPnL is the projected realized PnL on the
	// (possibly partial) fill.
	EstimatedRealizedPnL Decimal `json:"estimated_realized_pnl"`
	// EstimatedOrderStatus is the engine's projected lifecycle
	// state for the order on submission. Fully filled = "filled";
	// limit/maker = "open"; partially filled IOC/market =
	// "cancelled"; "rejected" / "expired" if margin or expiry
	// rules trip.
	EstimatedOrderStatus string `json:"estimated_order_status"`
	// SuggestedMaxFee is the engine's recommended `max_fee` value
	// for the trade (per contract).
	SuggestedMaxFee Decimal `json:"suggested_max_fee"`
	// PreInitialMargin is the user's initial-margin balance
	// before the simulated order.
	PreInitialMargin Decimal `json:"pre_initial_margin"`
	// PostInitialMargin is the hypothetical balance after the
	// order would be accepted.
	PostInitialMargin Decimal `json:"post_initial_margin"`
	// PostLiquidationPrice is the subaccount's liquidation price
	// if the order were fully filled. The wire field is nullable;
	// a zero-value [Decimal] indicates no projected liquidation
	// (e.g. spot order or fully-collateralised perp).
	PostLiquidationPrice Decimal `json:"post_liquidation_price"`
}

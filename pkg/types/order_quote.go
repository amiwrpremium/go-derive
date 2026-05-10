// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the response shape of `private/order_quote`.
package types

import "github.com/amiwrpremium/go-derive/pkg/enums"

// OrderQuoteResult is the response of `private/order_quote`. The
// endpoint runs a hypothetical order through the matching engine
// without submitting and reports the engine's estimates for fill
// price, fee, and post-trade margin balance — useful for
// pre-flighting orders against thin books.
//
// The shape mirrors the response documented at
// docs.derive.xyz/reference/private-order_quote (and shared with
// the public counterpart at /reference/public-order_quote).
type OrderQuoteResult struct {
	// IsValid reports whether the order is expected to clear
	// margin requirements.
	IsValid bool `json:"is_valid"`
	// InvalidReason carries the engine's reason when IsValid is
	// false. Empty (zero-value) when the request is valid. The wire
	// field is nullable; the documented value set is shared with
	// [BestQuoteResult.InvalidReason] — see [enums.RFQInvalidReason].
	InvalidReason enums.RFQInvalidReason `json:"invalid_reason,omitempty"`
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
	// EstimatedRealizedPnLExclFees is the projected realized PnL with
	// the fee component of cost basis excluded. For orders the engine
	// projects as `cancelled`, this only includes PnL on the filled
	// portion.
	EstimatedRealizedPnLExclFees Decimal `json:"estimated_realized_pnl_excl_fees"`
	// MaxAmount is the engine's projected upper bound on the
	// fill amount given the simulated request, useful for sizing
	// follow-on orders.
	MaxAmount Decimal `json:"max_amount"`
	// EstimatedOrderStatus is the engine's projected lifecycle
	// state for the order on submission. Fully filled =
	// [enums.OrderStatusFilled]; limit/maker =
	// [enums.OrderStatusOpen]; partially filled IOC/market =
	// [enums.OrderStatusCancelled]; [enums.OrderStatusRejected] /
	// [enums.OrderStatusExpired] if margin or expiry rules trip.
	EstimatedOrderStatus enums.OrderStatus `json:"estimated_order_status"`
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

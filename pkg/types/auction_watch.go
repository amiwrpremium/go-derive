// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the payload shape emitted on the `auctions.watch`
// public WebSocket channel.
package types

import (
	"encoding/json"

	"github.com/amiwrpremium/go-derive/pkg/enums"
)

// AuctionWatchEvent is one notification on the `auctions.watch`
// channel. Derive emits one event per state transition for every
// liquidation auction it is currently running across the platform.
//
// When the auction is ongoing the [Details] pointer is populated
// with everything a bidder needs to construct a `private/liquidate`
// call; when the state transitions to ended the pointer is nil.
type AuctionWatchEvent struct {
	// SubaccountID is the subaccount being auctioned off.
	SubaccountID int64 `json:"subaccount_id"`
	// State is the lifecycle state of the auction.
	State enums.AuctionState `json:"state"`
	// Timestamp is when this event was emitted (millisecond Unix
	// epoch).
	Timestamp MillisTime `json:"timestamp"`
	// Details is populated when the auction is ongoing and is nil
	// once it has ended.
	Details *AuctionWatchDetails `json:"details"`
}

// AuctionWatchDetails is the per-event live state of an ongoing
// liquidation auction. Each numeric field is a wire decimal string;
// see the field doc comments for the precise meaning.
type AuctionWatchDetails struct {
	// Currency of the subaccount being liquidated. The bidder is
	// recommended to use the same currency to avoid unsupported-
	// currency errors.
	Currency string `json:"currency"`
	// EstimatedBidPrice is the discounted mark value of the whole
	// subaccount; it does NOT scale with the bid percent. Negative
	// values indicate insolvent auctions.
	EstimatedBidPrice Decimal `json:"estimated_bid_price"`
	// EstimatedDiscountPnL is the estimated profit relative to
	// EstimatedMTM if the liquidation succeeds at
	// EstimatedPercentBid and EstimatedBidPrice.
	EstimatedDiscountPnL Decimal `json:"estimated_discount_pnl"`
	// EstimatedMTM is the un-discounted mark-to-market value of the
	// subaccount being liquidated.
	EstimatedMTM Decimal `json:"estimated_mtm"`
	// EstimatedPercentBid is an estimate of the maximum percent of
	// the subaccount that can be liquidated.
	EstimatedPercentBid Decimal `json:"estimated_percent_bid"`
	// LastSeenTradeID is the most recent trade id for the
	// subaccount; pass it to `private/liquidate` so the engine can
	// detect drift from on-chain state.
	LastSeenTradeID int64 `json:"last_seen_trade_id"`
	// MarginType is the margin regime of the subaccount being
	// liquidated.
	MarginType enums.MarginType `json:"margin_type"`
	// MinCashTransfer is the suggested minimum amount of cash to
	// transfer to a newly created bidder subaccount; unused funds
	// are returned.
	MinCashTransfer Decimal `json:"min_cash_transfer"`
	// MinPriceLimit is the estimated minimum `price_limit` for
	// `private/liquidate`.
	MinPriceLimit Decimal `json:"min_price_limit"`
	// SubaccountBalances is the current balance map of the
	// subaccount being auctioned. The shape is left as raw JSON
	// because Derive does not pin its inner schema; callers can
	// re-decode per their margin model.
	SubaccountBalances json.RawMessage `json:"subaccount_balances"`
}

// Package types.
package types

import "github.com/amiwrpremium/go-derive/pkg/enums"

// AuctionBid is one bid placed against a liquidation auction.
type AuctionBid struct {
	// Bidder is the wallet that placed the bid.
	Bidder Address `json:"bidder,omitempty"`
	// Price is the price the bidder offered.
	Price Decimal `json:"price,omitempty"`
	// PercentLiquidated is how much of the position the bid covers.
	PercentLiquidated Decimal `json:"percent_liquidated,omitempty"`
	// Timestamp is when the bid was received.
	Timestamp MillisTime `json:"timestamp,omitempty"`
}

// Liquidation is a liquidation-auction event reported by the engine.
//
// The canonical shape mirrors `derivexyz/cockpit`'s
// `AuctionResultSchema`. Public payloads carry the auction outcome;
// private payloads on the affected subaccount add per-position breakdowns.
type Liquidation struct {
	// AuctionID is the unique server-side auction id.
	AuctionID string `json:"auction_id,omitempty"`
	// AuctionType is "solvent" or "insolvent" — see [enums.AuctionType].
	AuctionType enums.AuctionType `json:"auction_type,omitempty"`
	// SubaccountID is the subaccount being liquidated.
	SubaccountID int64 `json:"subaccount_id,omitempty"`
	// StartTimestamp is when the auction opened.
	StartTimestamp MillisTime `json:"start_timestamp,omitempty"`
	// EndTimestamp is when the auction closed (zero for ongoing).
	EndTimestamp MillisTime `json:"end_timestamp,omitempty"`
	// Bids is the chronological list of bids placed during the auction.
	Bids []AuctionBid `json:"bids,omitempty"`
	// CashReceived is the cash flow into the affected subaccount.
	CashReceived Decimal `json:"cash_received,omitempty"`
	// DiscountPnL is the realized PnL contribution from the auction discount.
	DiscountPnL Decimal `json:"discount_pnl,omitempty"`
	// Fee is the auction fee charged.
	Fee Decimal `json:"fee,omitempty"`
	// PercentLiquidated is how much of the subaccount was liquidated [0, 1].
	PercentLiquidated Decimal `json:"percent_liquidated,omitempty"`
	// RealizedPnL is the total realized PnL across the auction.
	RealizedPnL Decimal `json:"realized_pnl,omitempty"`
	// AmountsLiquidated is per-instrument size that was liquidated.
	AmountsLiquidated map[string]Decimal `json:"amounts_liquidated,omitempty"`
	// PositionsRealizedPnL is per-instrument realized PnL.
	PositionsRealizedPnL map[string]Decimal `json:"positions_realized_pnl,omitempty"`
	// Timestamp is when the engine recorded the event.
	Timestamp MillisTime `json:"timestamp"`
	// TxHash is the on-chain liquidation transaction hash, if available.
	TxHash TxHash `json:"tx_hash,omitempty"`
}

// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the liquidation-auction history shape returned by
// `private/get_liquidation_history`.
package types

// LiquidationAuction is one entry in `private/get_liquidation_history`.
// Derive's liquidation engine runs auctions to wind down
// undercollateralised subaccounts; each auction can include multiple
// bids.
//
// The shape mirrors `AuctionHistoryResultSchema` in Derive's v2.2
// OpenAPI spec.
type LiquidationAuction struct {
	// AuctionID is the unique auction id.
	AuctionID string `json:"auction_id"`
	// AuctionType is "solvent" or "insolvent" depending on whether
	// the subaccount's equity was positive or negative when the
	// auction kicked off.
	AuctionType string `json:"auction_type"`
	// SubaccountID is the liquidated subaccount.
	SubaccountID int64 `json:"subaccount_id"`
	// StartTimestamp is when the auction began (millisecond Unix
	// epoch).
	StartTimestamp MillisTime `json:"start_timestamp"`
	// EndTimestamp is when the auction ended. The wire field is
	// nullable; the zero-value here means the auction is still
	// open.
	EndTimestamp MillisTime `json:"end_timestamp"`
	// Fee is the fee the subaccount paid to the auction.
	Fee Decimal `json:"fee"`
	// TxHash is the auction-completion transaction hash.
	TxHash TxHash `json:"tx_hash"`
	// Bids is the per-bid breakdown of the auction.
	Bids []AuctionBid `json:"bids"`
}

// AuctionBid is one bid against a [LiquidationAuction]. Each bid
// liquidates a percentage of the subaccount's portfolio and reports
// the per-asset breakdown of what was closed.
//
// The shape mirrors `AuctionBidEventSchema` in Derive's v2.2 OpenAPI
// spec.
type AuctionBid struct {
	// AmountsLiquidated is the per-asset notional liquidated, keyed
	// by instrument or asset symbol.
	AmountsLiquidated map[string]Decimal `json:"amounts_liquidated"`
	// CashReceived is the cash flow on the bid. For the liquidated
	// subaccount it's the amount the liquidator paid; for the
	// liquidator it's positive when the security module paid them
	// or negative when they paid for the bid.
	CashReceived Decimal `json:"cash_received"`
	// DiscountPnL is the realized PnL from the bid pricing being at
	// a discount to mark portfolio value.
	DiscountPnL Decimal `json:"discount_pnl"`
	// PercentLiquidated is the fraction of the subaccount closed by
	// this bid, expressed as a decimal (e.g. "0.25" for 25 %).
	PercentLiquidated Decimal `json:"percent_liquidated"`
	// PositionsRealizedPnL is the per-position realized PnL on the
	// bid, keyed by instrument name.
	PositionsRealizedPnL map[string]Decimal `json:"positions_realized_pnl"`
	// PositionsRealizedPnLExclFees is the same as
	// PositionsRealizedPnL excluding the fee component of cost
	// basis.
	PositionsRealizedPnLExclFees map[string]Decimal `json:"positions_realized_pnl_excl_fees"`
	// RealizedPnL is the bid's net realized PnL assuming positions
	// close at mark.
	RealizedPnL Decimal `json:"realized_pnl"`
	// RealizedPnLExclFees is RealizedPnL excluding fees from cost
	// basis.
	RealizedPnLExclFees Decimal `json:"realized_pnl_excl_fees"`
	// Timestamp is when the bid was placed (millisecond Unix
	// epoch).
	Timestamp MillisTime `json:"timestamp"`
	// TxHash is the bid transaction hash.
	TxHash TxHash `json:"tx_hash"`
}

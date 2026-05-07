// Package types — see address.go for the overview.
package types

import "github.com/amiwrpremium/go-derive/pkg/enums"

// Trade is one filled execution.
//
// Public-trade payloads (from `public/get_trade_history` and the
// `trades.{instrument}` channel) populate the public fields. Private
// trade payloads additionally populate OrderID, SubaccountID, Fee,
// LiquidityRole, RealizedPnL, plus the on-chain settlement fields
// (TxStatus, TxHash, Wallet) and the optional QuoteID / RFQID linking
// when the fill came out of an RFQ flow.
type Trade struct {
	// TradeID is the unique server-side trade id.
	TradeID string `json:"trade_id"`
	// OrderID identifies the order on the user's side of the fill (private
	// only; empty for public-trade payloads).
	OrderID string `json:"order_id,omitempty"`
	// SubaccountID identifies the user's subaccount (private only).
	SubaccountID int64 `json:"subaccount_id,omitempty"`
	// Wallet is the signer wallet that owns the trade. Public-trade payloads
	// include this for tape attribution.
	Wallet Address `json:"wallet,omitempty"`
	// InstrumentName identifies the market.
	InstrumentName string `json:"instrument_name"`
	// Direction is the user's side: buy means the user took the ask;
	// sell means the user hit the bid.
	Direction enums.Direction `json:"direction"`
	// TradePrice is the printed fill price.
	TradePrice Decimal `json:"trade_price"`
	// TradeAmount is the filled size in base-currency units.
	TradeAmount Decimal `json:"trade_amount"`
	// MarkPrice is the engine's mark price at the time of the fill.
	MarkPrice Decimal `json:"mark_price"`
	// IndexPrice is the underlying index price at the time of the fill.
	IndexPrice Decimal `json:"index_price,omitempty"`
	// Fee is the fee paid (private only).
	Fee Decimal `json:"trade_fee,omitempty"`
	// ExpectedRebate is the maker rebate the engine projects for this fill,
	// if any. Public payloads carry this for tape transparency.
	ExpectedRebate Decimal `json:"expected_rebate,omitempty"`
	// ExtraFee is any additional fee applied (e.g. settlement).
	ExtraFee Decimal `json:"extra_fee,omitempty"`
	// LiquidityRole identifies whether the user was maker or taker (private only).
	LiquidityRole enums.LiquidityRole `json:"liquidity_role,omitempty"`
	// RealizedPnL is the realized PnL from this fill (private only).
	RealizedPnL Decimal `json:"realized_pnl,omitempty"`
	// RealizedPnLExclFees is realized PnL with fees excluded (private only).
	RealizedPnLExclFees Decimal `json:"realized_pnl_excl_fees,omitempty"`
	// QuoteID links the trade to the quote it executed against, when the
	// fill came out of an RFQ flow.
	QuoteID string `json:"quote_id,omitempty"`
	// RFQID links the trade back to the originating RFQ, when applicable.
	RFQID string `json:"rfq_id,omitempty"`
	// TxHash is the on-chain transaction hash that carried the trade
	// settlement (populated once the engine submits the tx).
	TxHash TxHash `json:"tx_hash,omitempty"`
	// TxStatus reports the on-chain settlement state — see
	// [enums.TxStatus] for the allowed values.
	TxStatus enums.TxStatus `json:"tx_status,omitempty"`
	// Timestamp is the engine-side execution time.
	Timestamp MillisTime `json:"timestamp"`
}

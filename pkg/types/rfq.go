// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the RFQ (request-for-quote) flow shapes:
//   - the per-leg shapes [RFQLeg] and [QuoteLeg];
//   - the maker-side [Quote] (signer-facing, full envelope);
//   - the taker-side [QuotePublic] (no signing fields, has wallet);
//   - the [RFQ] result envelope returned by `private/get_rfqs`;
//   - the per-method response wrappers [ExecuteQuoteResult],
//     [CancelBatchResult], and [BestQuoteResult].
package types

import "github.com/amiwrpremium/go-derive/pkg/enums"

// RFQLeg is one leg of a multi-leg RFQ.
//
// Multi-leg RFQs are how Derive supports option spreads, calendars, etc.
// Each leg references its own instrument, direction and amount; legs
// must be unique by instrument (see
// [github.com/amiwrpremium/go-derive/pkg/errors.CodeLegInstrumentsNotUnique]).
//
// Mirrors `LegUnpricedSchema` in the OAS — RFQ legs do not carry
// per-leg prices.
type RFQLeg struct {
	// InstrumentName identifies the leg's market.
	InstrumentName string `json:"instrument_name"`
	// Direction is buy or sell on this leg.
	Direction enums.Direction `json:"direction"`
	// Amount is the leg's size in base-currency units.
	Amount Decimal `json:"amount"`
}

// Validate performs schema-level checks on the receiver: instrument
// non-empty, direction in range, amount positive. Returns nil on
// success or a wrapped [ErrInvalidParams].
func (l RFQLeg) Validate() error {
	if l.InstrumentName == "" {
		return invalidParam("instrument_name", "required")
	}
	if err := l.Direction.Validate(); err != nil {
		return invalidParam("direction", err.Error())
	}
	if l.Amount.Sign() <= 0 {
		return invalidParam("amount", "must be positive")
	}
	return nil
}

// QuoteLeg is a priced leg attached to a maker's [Quote] response.
//
// Distinct from [RFQLeg] because RFQs don't carry per-leg prices —
// quotes do. Mirrors `LegPricedSchema` in the OAS.
type QuoteLeg struct {
	// InstrumentName identifies the leg's market.
	InstrumentName string `json:"instrument_name"`
	// Direction is the maker's side on this leg.
	Direction enums.Direction `json:"direction"`
	// Amount is the leg's size in base-currency units.
	Amount Decimal `json:"amount"`
	// Price is the per-leg price the maker is committing to.
	Price Decimal `json:"price"`
}

// RFQ is a Request-For-Quote initiated by a taker.
//
// The taker broadcasts the RFQ to whitelisted makers; makers respond
// with [Quote] objects. The taker selects a quote and executes via
// `private/execute_quote`.
//
// Mirrors `RFQResultSchema` in Derive's v2.2 OpenAPI spec.
type RFQ struct {
	// RFQID is the unique server-side id.
	RFQID string `json:"rfq_id"`
	// SubaccountID is the taker's subaccount.
	SubaccountID int64 `json:"subaccount_id"`
	// Wallet is the taker's wallet address.
	Wallet string `json:"wallet"`
	// Status is the current lifecycle state.
	Status enums.QuoteStatus `json:"status"`
	// CancelReason is set when the RFQ has been cancelled.
	CancelReason enums.CancelReason `json:"cancel_reason"`
	// Legs is the per-instrument breakdown (no per-leg prices on RFQs).
	Legs []RFQLeg `json:"legs"`
	// Counterparties is the explicit list of maker wallets the RFQ
	// was sent to. Empty for an open RFQ.
	Counterparties []string `json:"counterparties"`
	// Label is the taker's free-form per-RFQ tag.
	Label string `json:"label"`
	// PreferredDirection is the side the taker prefers makers to
	// quote when the RFQ is two-sided.
	PreferredDirection enums.Direction `json:"preferred_direction"`
	// ReducingDirection is set when the RFQ is risk-reducing on a
	// specific side.
	ReducingDirection enums.Direction `json:"reducing_direction"`
	// FilledDirection is set after the RFQ has been (partially or
	// fully) filled, indicating which side filled.
	FilledDirection enums.Direction `json:"filled_direction"`
	// FilledPct is the fraction of the RFQ filled, 0…1.
	FilledPct Decimal `json:"filled_pct"`
	// MaxTotalCost is the upper bound on what the taker is willing
	// to pay across the package.
	MaxTotalCost Decimal `json:"max_total_cost"`
	// MinTotalCost is the lower bound (for sells, the minimum
	// proceeds the taker accepts).
	MinTotalCost Decimal `json:"min_total_cost"`
	// TotalCost is the realised total cost for filled fraction so
	// far.
	TotalCost Decimal `json:"total_cost"`
	// AskTotalCost is the all-in ask price for the package.
	AskTotalCost Decimal `json:"ask_total_cost"`
	// BidTotalCost is the all-in bid price for the package.
	BidTotalCost Decimal `json:"bid_total_cost"`
	// MarkTotalCost is the engine's mark-to-market estimate for the
	// package.
	MarkTotalCost Decimal `json:"mark_total_cost"`
	// PartialFillStep is the minimum partial-fill increment.
	PartialFillStep Decimal `json:"partial_fill_step"`
	// ValidUntil is the RFQ's expiry (millisecond Unix epoch).
	ValidUntil MillisTime `json:"valid_until"`
	// CreationTimestamp is when the RFQ was first received.
	CreationTimestamp MillisTime `json:"creation_timestamp"`
	// LastUpdateTimestamp is the most recent state change.
	LastUpdateTimestamp MillisTime `json:"last_update_timestamp"`
}

// Quote is a market-maker's signer-side view of one quote response
// to an RFQ. The full canonical shape mirrors `QuoteResultSchema` in
// the OAS — every field below corresponds to one schema property.
//
// Use [QuotePublic] when consuming the taker-side `poll_quotes` view
// or the `best_quote` field of `rfq_get_best_quote`.
type Quote struct {
	// QuoteID is the unique server-side id.
	QuoteID string `json:"quote_id"`
	// RFQID identifies the RFQ this quote responds to.
	RFQID string `json:"rfq_id"`
	// SubaccountID is the maker's subaccount.
	SubaccountID int64 `json:"subaccount_id"`
	// Direction is the maker's side; the taker's fill is the
	// opposite.
	Direction enums.Direction `json:"direction"`
	// Legs is the per-instrument breakdown (must match the RFQ
	// legs).
	Legs []QuoteLeg `json:"legs"`
	// LegsHash is a server-side hash that ties the quote to a
	// specific leg ordering — useful for replay protection.
	LegsHash string `json:"legs_hash"`
	// Status is the current lifecycle state.
	Status enums.QuoteStatus `json:"status"`
	// CancelReason is set when Status is QuoteStatusCancelled.
	CancelReason enums.CancelReason `json:"cancel_reason"`
	// LiquidityRole identifies the quote as maker-side (always
	// "maker").
	LiquidityRole enums.LiquidityRole `json:"liquidity_role"`
	// Fee is the fee charged on the quote (when filled).
	Fee Decimal `json:"fee"`
	// MaxFee is the maker's cap on the per-fill fee.
	MaxFee Decimal `json:"max_fee"`
	// ExtraFee is any extra USDC fee added by the referring client
	// (included in Fee).
	ExtraFee Decimal `json:"extra_fee"`
	// FillPct is the fraction of the quote filled, 0…1.
	FillPct Decimal `json:"fill_pct"`
	// IsTransfer indicates the quote was settled as an internal
	// transfer rather than an on-chain trade.
	IsTransfer bool `json:"is_transfer"`
	// Label is the maker's free-form per-quote tag.
	Label string `json:"label"`
	// MMP indicates the quote participated in market-maker-protection
	// accounting on the maker's subaccount.
	MMP bool `json:"mmp"`
	// Nonce is the maker's signed nonce.
	Nonce uint64 `json:"nonce"`
	// Signer is the address that signed the quote action.
	Signer Address `json:"signer"`
	// Signature is the EIP-712 action signature.
	Signature string `json:"signature"`
	// SignatureExpiry is the Unix timestamp (seconds) past which
	// the signature is rejected.
	SignatureExpiry int64 `json:"signature_expiry_sec"`
	// TxHash is the on-chain settlement transaction hash, set
	// after the quote is executed.
	TxHash TxHash `json:"tx_hash"`
	// TxStatus is the on-chain settlement state.
	TxStatus enums.TxStatus `json:"tx_status"`
	// CreationTimestamp is when the quote was received.
	CreationTimestamp MillisTime `json:"creation_timestamp"`
	// LastUpdateTimestamp is the most recent state change.
	LastUpdateTimestamp MillisTime `json:"last_update_timestamp"`
}

// QuotePublic is the taker-side view of a quote — the same shape as
// [Quote] minus the maker's signing fields (`fee`, `extra_fee`,
// `is_transfer`, `label`, `max_fee`, `mmp`, `nonce`, `signature`,
// `signature_expiry_sec`, `signer`) and plus the maker's wallet.
//
// Returned by `private/poll_quotes` and as the `best_quote` field of
// `private/rfq_get_best_quote`. Mirrors `QuoteResultPublicSchema`.
type QuotePublic struct {
	// QuoteID is the unique server-side id.
	QuoteID string `json:"quote_id"`
	// RFQID identifies the RFQ this quote responds to.
	RFQID string `json:"rfq_id"`
	// SubaccountID is the maker's subaccount.
	SubaccountID int64 `json:"subaccount_id"`
	// Wallet is the maker's wallet address (the taker-visible
	// counterparty identifier).
	Wallet string `json:"wallet"`
	// Direction is the maker's side.
	Direction enums.Direction `json:"direction"`
	// Legs is the per-instrument priced breakdown.
	Legs []QuoteLeg `json:"legs"`
	// LegsHash is the engine's hash over the leg ordering.
	LegsHash string `json:"legs_hash"`
	// Status is the current lifecycle state.
	Status enums.QuoteStatus `json:"status"`
	// CancelReason is set when Status is QuoteStatusCancelled.
	CancelReason enums.CancelReason `json:"cancel_reason"`
	// LiquidityRole identifies the quote as maker-side.
	LiquidityRole enums.LiquidityRole `json:"liquidity_role"`
	// FillPct is the fraction of the quote filled, 0…1.
	FillPct Decimal `json:"fill_pct"`
	// TxHash is the on-chain settlement transaction hash.
	TxHash TxHash `json:"tx_hash"`
	// TxStatus is the on-chain settlement state.
	TxStatus enums.TxStatus `json:"tx_status"`
	// CreationTimestamp is when the quote was received.
	CreationTimestamp MillisTime `json:"creation_timestamp"`
	// LastUpdateTimestamp is the most recent state change.
	LastUpdateTimestamp MillisTime `json:"last_update_timestamp"`
}

// ExecuteQuoteResult is the response of `private/execute_quote`. It's
// the same as [Quote] (signer-side full envelope) plus the
// `rfq_filled_pct` field reporting what fraction of the RFQ this
// execution covered.
//
// Mirrors `PrivateExecuteQuoteResultSchema`.
type ExecuteQuoteResult struct {
	Quote
	// RFQFilledPct is the fraction of the originating RFQ filled
	// after this execution, 0…1.
	RFQFilledPct Decimal `json:"rfq_filled_pct"`
}

// ReplaceQuoteResult is the response of `private/replace_quote`. The
// endpoint cancels one outstanding maker quote and submits a
// replacement atomically; the response carries both the cancelled
// quote and the (optional) newly created quote, plus an optional
// engine error if the replacement was rejected.
//
// Mirrors `PrivateReplaceQuoteResultSchema` — the quote-side
// counterpart to [ReplaceResult] for orders.
type ReplaceQuoteResult struct {
	// CancelledQuote is the quote that was dropped.
	CancelledQuote Quote `json:"cancelled_quote"`
	// Quote is the replacement quote. Nil when CreateQuoteError is
	// non-nil — the engine cancelled the old quote but rejected the
	// new one.
	Quote *Quote `json:"quote,omitempty"`
	// CreateQuoteError is the engine error if the replacement quote
	// was rejected. Non-nil only on failure.
	CreateQuoteError *RPCError `json:"create_quote_error,omitempty"`
}

// CancelBatchResult is the response of both
// `private/cancel_batch_quotes` and `private/cancel_batch_rfqs`.
type CancelBatchResult struct {
	// CancelledIDs is the list of quote or RFQ ids that were
	// cancelled by the call.
	CancelledIDs []string `json:"cancelled_ids"`
}

// BestQuoteResult is the response of `private/rfq_get_best_quote`.
// The endpoint surveys live quotes against an RFQ shape and returns
// the best one along with the engine's margin-impact estimates. Many
// fields are nullable on the wire when the engine cannot compute an
// estimate; those decode to a zero-value [Decimal].
//
// Mirrors `PrivateRfqGetBestQuoteResultSchema`.
type BestQuoteResult struct {
	// BestQuote is the matched maker quote, or nil if no acceptable
	// quote exists.
	BestQuote *QuotePublic `json:"best_quote,omitempty"`
	// Direction is the RFQ direction.
	Direction enums.Direction `json:"direction"`
	// IsValid reports whether the RFQ is expected to clear margin.
	IsValid bool `json:"is_valid"`
	// InvalidReason carries a human-readable reason when IsValid
	// is false. Empty otherwise. The wire field is nullable.
	InvalidReason string `json:"invalid_reason,omitempty"`
	// EstimatedFee is the engine's estimate of the fee on the
	// trade.
	EstimatedFee Decimal `json:"estimated_fee"`
	// EstimatedRealizedPnL is the engine's estimate of realized
	// PnL on the trade.
	EstimatedRealizedPnL Decimal `json:"estimated_realized_pnl"`
	// EstimatedRealizedPnLExclFees is EstimatedRealizedPnL with the
	// fee component of cost basis excluded.
	EstimatedRealizedPnLExclFees Decimal `json:"estimated_realized_pnl_excl_fees"`
	// EstimatedTotalCost is the engine's estimate of the trade's
	// total dollar cost.
	EstimatedTotalCost Decimal `json:"estimated_total_cost"`
	// FilledPct is the fraction of the RFQ already filled, 0…1.
	FilledPct Decimal `json:"filled_pct"`
	// OrderbookTotalCost is the alternative total cost if filled
	// against the lit book; nullable on the wire when any leg
	// lacks orderbook depth.
	OrderbookTotalCost Decimal `json:"orderbook_total_cost"`
	// SuggestedMaxFee is the engine's recommended `max_fee` for
	// the trade.
	SuggestedMaxFee Decimal `json:"suggested_max_fee"`
	// PreInitialMargin is the user's initial margin before the
	// trade.
	PreInitialMargin Decimal `json:"pre_initial_margin"`
	// PostInitialMargin is the user's hypothetical initial margin
	// after the trade.
	PostInitialMargin Decimal `json:"post_initial_margin"`
	// PostLiquidationPrice is the post-trade liquidation price
	// (closest to the index). Nullable.
	PostLiquidationPrice Decimal `json:"post_liquidation_price"`
	// DownLiquidationPrice is the post-trade downside liquidation
	// price. Nullable.
	DownLiquidationPrice Decimal `json:"down_liquidation_price"`
	// UpLiquidationPrice is the post-trade upside liquidation
	// price. Nullable.
	UpLiquidationPrice Decimal `json:"up_liquidation_price"`
}

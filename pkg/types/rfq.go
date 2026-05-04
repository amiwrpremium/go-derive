package types

import "github.com/amiwrpremium/go-derive/pkg/enums"

// RFQLeg is one leg of a multi-leg RFQ.
//
// Multi-leg RFQs are how Derive supports option spreads, calendars, etc.
// Each leg references its own instrument, direction and amount; legs must
// be unique by instrument (see [github.com/amiwrpremium/go-derive/pkg/errors.CodeLegInstrumentsNotUnique]).
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
// quotes do.
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
// The taker broadcasts the RFQ to whitelisted makers; makers respond with
// [Quote] objects. The taker selects a quote and executes via
// `private/execute_quote`.
type RFQ struct {
	// RFQID is the unique server-side id.
	RFQID string `json:"rfq_id"`
	// SubaccountID is the taker's subaccount.
	SubaccountID int64 `json:"subaccount_id"`
	// Status is the current lifecycle state.
	Status enums.QuoteStatus `json:"status"`
	// CancelReason is set when Status is QuoteStatusCancelled.
	CancelReason enums.CancelReason `json:"cancel_reason,omitempty"`
	// Legs is the per-instrument breakdown (no per-leg prices on RFQs).
	Legs []RFQLeg `json:"legs"`
	// MaxFee is the cap on total fee the taker is willing to pay.
	MaxFee Decimal `json:"max_total_fee,omitempty"`
	// MinPrice and MaxPrice constrain the price band the taker accepts.
	MinPrice Decimal `json:"min_price,omitempty"`
	MaxPrice Decimal `json:"max_price,omitempty"`
	// CreationTimestamp is when the RFQ was first received.
	CreationTimestamp MillisTime `json:"creation_timestamp"`
	// LastUpdateTimestamp is the most recent state change.
	LastUpdateTimestamp MillisTime `json:"last_update_timestamp"`
}

// Quote is a market-maker response to an [RFQ].
//
// The full canonical shape mirrors `derivexyz/cockpit`'s
// `QuoteResultSchema` — every field below appears on the wire.
type Quote struct {
	// QuoteID is the unique server-side id.
	QuoteID string `json:"quote_id"`
	// RFQID identifies the RFQ this quote responds to.
	RFQID string `json:"rfq_id"`
	// SubaccountID is the maker's subaccount.
	SubaccountID int64 `json:"subaccount_id"`
	// Direction is the maker's side; the taker's fill is the opposite.
	Direction enums.Direction `json:"direction"`
	// Legs is the per-instrument breakdown (must match the RFQ legs).
	Legs []QuoteLeg `json:"legs"`
	// LegsHash is a server-side hash that ties the quote to a specific
	// leg ordering — useful for replay protection.
	LegsHash string `json:"legs_hash,omitempty"`
	// Price is the all-in net price for the package.
	Price Decimal `json:"price,omitempty"`
	// Fee is the fee charged on the quote (when filled).
	Fee Decimal `json:"fee,omitempty"`
	// MaxFee is the maker's cap on the per-fill fee.
	MaxFee Decimal `json:"max_fee,omitempty"`
	// LiquidityRole identifies the quote as maker-side (always "maker").
	LiquidityRole enums.LiquidityRole `json:"liquidity_role,omitempty"`
	// Status is the current lifecycle state.
	Status enums.QuoteStatus `json:"status"`
	// CancelReason is set when Status is QuoteStatusCancelled.
	CancelReason enums.CancelReason `json:"cancel_reason,omitempty"`
	// MMP indicates the quote participated in market-maker-protection
	// accounting on the maker's subaccount.
	MMP bool `json:"mmp,omitempty"`
	// Label is the maker's free-form per-quote tag.
	Label string `json:"label,omitempty"`
	// Nonce is the maker's signed nonce.
	Nonce uint64 `json:"nonce,omitempty"`
	// Signer is the address that signed the quote action.
	Signer Address `json:"signer,omitempty"`
	// Signature is the EIP-712 action signature.
	Signature string `json:"signature,omitempty"`
	// SignatureExpiry is the Unix timestamp (seconds) past which the
	// signature is rejected.
	SignatureExpiry int64 `json:"signature_expiry_sec,omitempty"`
	// TxHash is the on-chain settlement transaction hash, set after the
	// quote is executed.
	TxHash TxHash `json:"tx_hash,omitempty"`
	// TxStatus is the on-chain settlement state.
	TxStatus enums.TxStatus `json:"tx_status,omitempty"`
	// CreationTimestamp is when the quote was received.
	CreationTimestamp MillisTime `json:"creation_timestamp"`
	// LastUpdateTimestamp is the most recent state change.
	LastUpdateTimestamp MillisTime `json:"last_update_timestamp,omitempty"`
}

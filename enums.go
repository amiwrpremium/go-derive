// Package enums declares the named-string enums used across the SDK.
//
// Each enum is a defined string type — the simplest idiom in Go that gives
// you exhaustive switch warnings, free JSON round-trips, and a domain-specific
// receiver set without the heavyweight ceremony of an `iota` block plus
// custom marshalers. Aliases of underlying string types like:
//
//	type Direction string
//	const DirectionBuy Direction = "buy"
//
// match what big Go SDKs (aws-sdk-go-v2, stripe-go) use, and the wire format
// they produce is byte-for-byte what Derive expects.
//
// Every enum exposes a Valid method for cheap input validation. Some, like
// [Direction], expose extra domain helpers ([Direction.Sign],
// [Direction.Opposite], [OrderStatus.Terminal]).
//
// # Validating untrusted input
//
// Always check [Direction.Valid] (or the corresponding Valid method on the
// enum) before passing user-provided strings into the SDK. The Go type
// system can't prevent constructing an out-of-range value via `Direction("x")`,
// so the runtime check is the safety net.
package derive

import (
	"fmt"
)

// AssetType identifies the kind of asset that backs a [Balance] or
// [Collateral] entry. The set is the same as [InstrumentType] — Derive
// reuses the same three categories — but the wire field is named
// `asset_type` rather than `instrument_type` on those payloads.
type AssetType string

const (
	// AssetTypeERC20 is a spot ERC-20 token (typically used as collateral).
	AssetTypeERC20 AssetType = "erc20"
	// AssetTypeOption is an option contract.
	AssetTypeOption AssetType = "option"
	// AssetTypePerp is a perpetual futures contract.
	AssetTypePerp AssetType = "perp"
)

// Valid reports whether the receiver is one of the defined asset types.
func (a AssetType) Valid() bool {
	switch a {
	case AssetTypeERC20, AssetTypeOption, AssetTypePerp:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (a AssetType) Validate() error {
	if a.Valid() {
		return nil
	}
	return invalid("AssetType", string(a))
}

// AuctionType describes the regime of a liquidation auction Derive ran
// against an undercollateralised subaccount.
type AuctionType string

const (
	// AuctionTypeSolvent — subaccount equity is positive; auction
	// transfers positions to keep things solvent.
	AuctionTypeSolvent AuctionType = "solvent"
	// AuctionTypeInsolvent — subaccount equity is negative; insurance
	// fund / socialised-loss path activates.
	AuctionTypeInsolvent AuctionType = "insolvent"
)

// Valid reports whether the receiver is one of the defined auction types.
func (a AuctionType) Valid() bool {
	switch a {
	case AuctionTypeSolvent, AuctionTypeInsolvent:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (a AuctionType) Validate() error {
	if a.Valid() {
		return nil
	}
	return invalid("AuctionType", string(a))
}

// BalanceUpdateType is the wire enum that classifies one entry on the
// `subaccount.{id}.balances` channel. Each balance change carries one
// of these values so consumers know what bookkeeping caused it.
type BalanceUpdateType string

const (
	// BalanceUpdateTrade — balance change from a fill.
	BalanceUpdateTrade BalanceUpdateType = "trade"
	// BalanceUpdateAssetDeposit — ERC-20 deposited into a subaccount.
	BalanceUpdateAssetDeposit BalanceUpdateType = "asset_deposit"
	// BalanceUpdateAssetWithdrawal — ERC-20 withdrawn from a subaccount.
	BalanceUpdateAssetWithdrawal BalanceUpdateType = "asset_withdrawal"
	// BalanceUpdateTransfer — value moved between subaccounts on the same wallet.
	BalanceUpdateTransfer BalanceUpdateType = "transfer"
	// BalanceUpdateSubaccountDeposit — collateral moved into the subaccount.
	BalanceUpdateSubaccountDeposit BalanceUpdateType = "subaccount_deposit"
	// BalanceUpdateSubaccountWithdrawal — collateral moved out of the subaccount.
	BalanceUpdateSubaccountWithdrawal BalanceUpdateType = "subaccount_withdrawal"
	// BalanceUpdateLiquidation — liquidation auction settlement.
	BalanceUpdateLiquidation BalanceUpdateType = "liquidation"
	// BalanceUpdateOnchainDriftFix — reconciliation against on-chain state.
	BalanceUpdateOnchainDriftFix BalanceUpdateType = "onchain_drift_fix"
	// BalanceUpdatePerpSettlement — perpetual mark-to-market cash flow.
	BalanceUpdatePerpSettlement BalanceUpdateType = "perp_settlement"
	// BalanceUpdateOptionSettlement — option exercise/expiry cash flow.
	BalanceUpdateOptionSettlement BalanceUpdateType = "option_settlement"
	// BalanceUpdateInterestAccrual — interest accrual on borrowed/lent collateral.
	BalanceUpdateInterestAccrual BalanceUpdateType = "interest_accrual"
	// BalanceUpdateOnchainRevert — on-chain transaction reverted post-execution.
	BalanceUpdateOnchainRevert BalanceUpdateType = "onchain_revert"
	// BalanceUpdateDoubleRevert — double-revert recovery path.
	BalanceUpdateDoubleRevert BalanceUpdateType = "double_revert"
)

// Valid reports whether the receiver is one of the defined update types.
func (u BalanceUpdateType) Valid() bool {
	switch u {
	case BalanceUpdateTrade, BalanceUpdateAssetDeposit, BalanceUpdateAssetWithdrawal,
		BalanceUpdateTransfer, BalanceUpdateSubaccountDeposit, BalanceUpdateSubaccountWithdrawal,
		BalanceUpdateLiquidation, BalanceUpdateOnchainDriftFix, BalanceUpdatePerpSettlement,
		BalanceUpdateOptionSettlement, BalanceUpdateInterestAccrual,
		BalanceUpdateOnchainRevert, BalanceUpdateDoubleRevert:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (u BalanceUpdateType) Validate() error {
	if u.Valid() {
		return nil
	}
	return invalid("BalanceUpdateType", string(u))
}

// CancelReason carries the engine's explanation for why an order was
// cancelled. It is reported on the canonical Order record (and on quote
// updates) — `""` for orders that have not been cancelled, otherwise one
// of the values below.
//
// The set mirrors the canonical `derivexyz/cockpit` enum exactly.
type CancelReason string

const (
	// CancelReasonNone is the empty wire value used while the order is
	// still open (no cancel has happened).
	CancelReasonNone CancelReason = ""
	// CancelReasonUserRequest means the order was cancelled by an explicit
	// `private/cancel` (or label/instrument/all variant).
	CancelReasonUserRequest CancelReason = "user_request"
	// CancelReasonMMP means market-maker protection tripped and pulled
	// the order along with the rest of the maker's book.
	CancelReasonMMP CancelReason = "mmp_trigger"
	// CancelReasonInsufficientMargin means the engine cancelled the order
	// because filling it would breach the subaccount's margin rules.
	CancelReasonInsufficientMargin CancelReason = "insufficient_margin"
	// CancelReasonSignedMaxFeeTooLow means the signed `max_fee` field was
	// below the venue's required minimum at the time of fill.
	CancelReasonSignedMaxFeeTooLow CancelReason = "signed_max_fee_too_low"
	// CancelReasonIOC means an IOC or market order partially filled and
	// the remainder was cancelled by IOC semantics.
	CancelReasonIOC CancelReason = "ioc_or_market_partial_fill"
	// CancelReasonCancelOnDisconnect means the kill-switch fired because
	// the wallet's authenticated WebSocket session disconnected.
	CancelReasonCancelOnDisconnect CancelReason = "cancel_on_disconnect"
	// CancelReasonSessionKey means the signing session key was deregistered,
	// which invalidates all of its outstanding orders.
	CancelReasonSessionKey CancelReason = "session_key_deregistered"
	// CancelReasonSubaccountWithdrawn means the subaccount was withdrawn
	// from the venue, taking its outstanding orders with it.
	CancelReasonSubaccountWithdrawn CancelReason = "subaccount_withdrawn"
	// CancelReasonCompliance means the wallet was placed in a restricted
	// compliance state and its open orders were pulled.
	CancelReasonCompliance CancelReason = "compliance"
)

// Valid reports whether the receiver is one of the defined cancel reasons.
//
// `CancelReasonNone` (the empty string) counts as valid — that is the
// wire value for "still open, never cancelled".
func (c CancelReason) Valid() bool {
	switch c {
	case CancelReasonNone, CancelReasonUserRequest, CancelReasonMMP,
		CancelReasonInsufficientMargin, CancelReasonSignedMaxFeeTooLow,
		CancelReasonIOC, CancelReasonCancelOnDisconnect, CancelReasonSessionKey,
		CancelReasonSubaccountWithdrawn, CancelReasonCompliance:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (c CancelReason) Validate() error {
	if c.Valid() {
		return nil
	}
	return invalid("CancelReason", string(c))
}

// Direction is the side of a trade or order. Buy orders consume asks; sell
// orders consume bids.
type Direction string

const (
	// DirectionBuy means the order or trade is on the bid side.
	DirectionBuy Direction = "buy"
	// DirectionSell means the order or trade is on the ask side.
	DirectionSell Direction = "sell"
)

// Valid reports whether the receiver equals one of the defined directions.
// Use it to gate untrusted input before [Direction.Sign] or [Direction.Opposite]
// (both of which assume validity).
func (d Direction) Valid() bool {
	switch d {
	case DirectionBuy, DirectionSell:
		return true
	default:
		return false
	}
}

// Sign returns +1 for [DirectionBuy] and -1 for [DirectionSell]. It is
// useful when computing signed position deltas:
//
//	delta := amount.Mul(decimal.NewFromInt(int64(side.Sign())))
//
// Sign panics on values that haven't passed [Direction.Valid]. Validate
// untrusted input first.
func (d Direction) Sign() int {
	switch d {
	case DirectionBuy:
		return 1
	case DirectionSell:
		return -1
	default:
		panic("enums: Direction.Sign called on invalid value " + string(d))
	}
}

// Opposite returns the reverse of d. Used in cancel-and-reverse logic and
// when computing offsetting orders for hedges.
//
// Note: Opposite returns [DirectionBuy] for any non-[DirectionBuy] input,
// including invalid values. Combine with [Direction.Valid] when input
// trustworthiness matters.
func (d Direction) Opposite() Direction {
	if d == DirectionBuy {
		return DirectionSell
	}
	return DirectionBuy
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (d Direction) Validate() error {
	if d.Valid() {
		return nil
	}
	return invalid("Direction", string(d))
}

// Environment selects which Derive deployment a client talks to. It is
// surfaced to users via the With* options on each client; the SDK turns
// it into the corresponding [NetworkConfig].
type Environment string

const (
	// EnvironmentMainnet selects the production deployment (chain ID 957).
	EnvironmentMainnet Environment = "mainnet"
	// EnvironmentTestnet selects the demo/staging deployment (chain ID 901).
	EnvironmentTestnet Environment = "testnet"
)

// Valid reports whether the receiver is one of the defined environments.
func (e Environment) Valid() bool {
	switch e {
	case EnvironmentMainnet, EnvironmentTestnet:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (e Environment) Validate() error {
	if e.Valid() {
		return nil
	}
	return invalid("Environment", string(e))
}

// validationError is returned by every enum's Validate() method when the
// receiver carries a value not in the canonical set. It implements the
// error interface and is comparable with errors.Is against
// [ErrInvalidEnum].
type validationError struct {
	enum  string
	value string
}

func (e *validationError) Error() string {
	return fmt.Sprintf("enums: invalid %s value %q", e.enum, e.value)
}

// Is satisfies errors.Is so callers can match without unwrapping. Every
// enum-validation failure unwraps to [ErrInvalidEnum].
func (e *validationError) Is(target error) bool { return target == ErrInvalidEnum }

// ErrInvalidEnum is the sentinel returned from every enum's Validate
// method when the receiver isn't one of the defined wire values. Use
// errors.Is to detect it.
var ErrInvalidEnum = &validationError{enum: "<unknown>", value: ""}

// invalid is the package-internal helper each Validate method calls to
// build a concrete error if Valid() returned false.
func invalid(enum, value string) error {
	return &validationError{enum: enum, value: value}
}

// InstrumentType identifies the kind of contract a market quotes.
//
// Derive supports three: linear perpetuals, options (with expiries and
// strikes), and ERC-20 spot tokens used as collateral or for spot trading.
type InstrumentType string

const (
	// InstrumentTypePerp is a perpetual futures contract with continuous
	// funding payments and no fixed expiry.
	InstrumentTypePerp InstrumentType = "perp"
	// InstrumentTypeOption is a European-style option with a strike and
	// expiry; see [OptionDetails].
	InstrumentTypeOption InstrumentType = "option"
	// InstrumentTypeERC20 is a spot ERC-20 token (typically used as collateral).
	InstrumentTypeERC20 InstrumentType = "erc20"
)

// Valid reports whether the receiver is one of the defined instrument types.
func (k InstrumentType) Valid() bool {
	switch k {
	case InstrumentTypePerp, InstrumentTypeOption, InstrumentTypeERC20:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (k InstrumentType) Validate() error {
	if k.Valid() {
		return nil
	}
	return invalid("InstrumentType", string(k))
}

// LiquidityRole is the role of a fill from the perspective of one side of
// the trade. Maker fills sit on the book first; taker fills cross the book.
//
// Fees and rebates differ by role on every Derive market.
type LiquidityRole string

const (
	// LiquidityRoleMaker means the side provided liquidity (the resting
	// order).
	LiquidityRoleMaker LiquidityRole = "maker"
	// LiquidityRoleTaker means the side consumed liquidity (the crossing
	// order).
	LiquidityRoleTaker LiquidityRole = "taker"
)

// Valid reports whether the receiver is one of the defined roles.
func (r LiquidityRole) Valid() bool {
	switch r {
	case LiquidityRoleMaker, LiquidityRoleTaker:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (r LiquidityRole) Validate() error {
	if r.Valid() {
		return nil
	}
	return invalid("LiquidityRole", string(r))
}

// MarginType identifies which margin model a subaccount uses.
//
// Derive supports three: Standard Margin (`SM`), Portfolio Margin (`PM`),
// and the second-generation Portfolio Margin model (`PM2`). The wire
// values are uppercase abbreviations.
type MarginType string

const (
	// MarginTypeSM is Standard Margin — per-position margin computed
	// independently of the rest of the book. Most permissive accounts.
	MarginTypeSM MarginType = "SM"
	// MarginTypePM is the original Portfolio Margin model — netted
	// margin across the whole subaccount.
	MarginTypePM MarginType = "PM"
	// MarginTypePM2 is the second-generation Portfolio Margin model.
	MarginTypePM2 MarginType = "PM2"
)

// Valid reports whether the receiver is one of the defined margin types.
func (m MarginType) Valid() bool {
	switch m {
	case MarginTypeSM, MarginTypePM, MarginTypePM2:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (m MarginType) Validate() error {
	if m.Valid() {
		return nil
	}
	return invalid("MarginType", string(m))
}

// OptionType distinguishes calls from puts. It only applies when the
// surrounding instrument's [InstrumentType] is [InstrumentTypeOption].
//
// The wire format is single-letter — Derive emits `"C"` for calls and
// `"P"` for puts.
type OptionType string

const (
	// OptionTypeCall gives the holder the right to buy the underlying at
	// the strike on or before expiry.
	OptionTypeCall OptionType = "C"
	// OptionTypePut gives the holder the right to sell the underlying at
	// the strike on or before expiry.
	OptionTypePut OptionType = "P"
)

// Valid reports whether the receiver is one of the defined option types.
func (o OptionType) Valid() bool {
	switch o {
	case OptionTypeCall, OptionTypePut:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (o OptionType) Validate() error {
	if o.Valid() {
		return nil
	}
	return invalid("OptionType", string(o))
}

// OrderStatus is the lifecycle state of an order as reported by the
// matching engine. Use [OrderStatus.Terminal] to test for "no further
// updates expected".
//
// The set mirrors the canonical `derivexyz/cockpit` enum exactly:
// open, filled, cancelled, expired, rejected.
type OrderStatus string

const (
	// OrderStatusOpen means the order is resting on the book.
	OrderStatusOpen OrderStatus = "open"
	// OrderStatusFilled means the order has been completely matched.
	OrderStatusFilled OrderStatus = "filled"
	// OrderStatusCancelled means the order was cancelled by the user, the
	// session-key, or the engine before it filled. The associated
	// [CancelReason] explains which.
	OrderStatusCancelled OrderStatus = "cancelled"
	// OrderStatusExpired means the order's signature expiry passed before
	// it filled.
	OrderStatusExpired OrderStatus = "expired"
	// OrderStatusRejected means the engine rejected the order at submission
	// time (e.g. invalid price, post-only would cross).
	OrderStatusRejected OrderStatus = "rejected"
)

// Valid reports whether the receiver is one of the defined statuses.
func (s OrderStatus) Valid() bool {
	switch s {
	case OrderStatusOpen, OrderStatusFilled, OrderStatusCancelled,
		OrderStatusExpired, OrderStatusRejected:
		return true
	default:
		return false
	}
}

// Terminal reports whether the status is final — i.e. the order will not
// receive further updates and can be cleaned out of any in-memory cache.
//
// Only Open is non-terminal; everything else is.
func (s OrderStatus) Terminal() bool {
	switch s {
	case OrderStatusFilled, OrderStatusCancelled, OrderStatusExpired,
		OrderStatusRejected:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (s OrderStatus) Validate() error {
	if s.Valid() {
		return nil
	}
	return invalid("OrderStatus", string(s))
}

// OrderType describes how an order interacts with the order book.
//
// [OrderTypeLimit] orders rest on the book at a stated price; [OrderTypeMarket]
// orders cross the book immediately at the best available price subject to
// the user's slippage cap.
type OrderType string

const (
	// OrderTypeLimit is a price-limited order that rests on the book until
	// it crosses, expires, or is cancelled.
	OrderTypeLimit OrderType = "limit"
	// OrderTypeMarket is an order that crosses the book immediately,
	// constrained only by the caller's max-fee and limit-price cap.
	OrderTypeMarket OrderType = "market"
)

// Valid reports whether the receiver is one of the defined order types.
func (t OrderType) Valid() bool {
	switch t {
	case OrderTypeLimit, OrderTypeMarket:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (t OrderType) Validate() error {
	if t.Valid() {
		return nil
	}
	return invalid("OrderType", string(t))
}

// QuoteStatus is the lifecycle state of a maker [Quote] response to an
// open RFQ.
//
// The set is the same four values [OrderStatus] supports — open / filled
// / cancelled / expired — but exists as its own type so a `Quote.Status`
// field cannot be confused with an `Order.OrderStatus` at the type
// level.
type QuoteStatus string

const (
	// QuoteStatusOpen — quote is live and may still be executed.
	QuoteStatusOpen QuoteStatus = "open"
	// QuoteStatusFilled — taker executed the quote.
	QuoteStatusFilled QuoteStatus = "filled"
	// QuoteStatusCancelled — quote was cancelled (by maker or engine).
	QuoteStatusCancelled QuoteStatus = "cancelled"
	// QuoteStatusExpired — quote's signature_expiry_sec passed.
	QuoteStatusExpired QuoteStatus = "expired"
)

// Valid reports whether the receiver is one of the defined statuses.
func (s QuoteStatus) Valid() bool {
	switch s {
	case QuoteStatusOpen, QuoteStatusFilled, QuoteStatusCancelled, QuoteStatusExpired:
		return true
	default:
		return false
	}
}

// Terminal reports whether the status is final.
func (s QuoteStatus) Terminal() bool {
	switch s {
	case QuoteStatusFilled, QuoteStatusCancelled, QuoteStatusExpired:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (s QuoteStatus) Validate() error {
	if s.Valid() {
		return nil
	}
	return invalid("QuoteStatus", string(s))
}

// TimeInForce controls when, and under what conditions, an order becomes
// inactive. The matching engine consults the time-in-force as soon as the
// order is accepted; a TIF mismatch (e.g. PostOnly on a market order) yields
// a synchronous rejection.
type TimeInForce string

const (
	// TimeInForceGTC ("good-till-cancelled") keeps the order open until the
	// caller cancels it or it expires for another reason.
	TimeInForceGTC TimeInForce = "gtc"
	// TimeInForcePostOnly rejects the order if it would cross the book at
	// submission time. Used by makers to guarantee maker rebates.
	TimeInForcePostOnly TimeInForce = "post_only"
	// TimeInForceFOK ("fill-or-kill") requires the order to fill in full
	// immediately or be cancelled. Partial fills are not allowed.
	TimeInForceFOK TimeInForce = "fok"
	// TimeInForceIOC ("immediate-or-cancel") fills as much as it can right
	// now and cancels any remaining quantity. Partial fills are allowed.
	TimeInForceIOC TimeInForce = "ioc"
)

// Valid reports whether the receiver is one of the defined TIFs.
func (t TimeInForce) Valid() bool {
	switch t {
	case TimeInForceGTC, TimeInForcePostOnly, TimeInForceFOK, TimeInForceIOC:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (t TimeInForce) Validate() error {
	if t.Valid() {
		return nil
	}
	return invalid("TimeInForce", string(t))
}

// TxStatus is the on-chain transaction lifecycle as reported by Derive.
//
// Trades, transfers, and other actions that emit an on-chain transaction
// carry a `tx_status` field that walks through these values: requested
// → pending → settled (or reverted/ignored on failure).
type TxStatus string

const (
	// TxStatusRequested — request accepted, transaction not yet submitted.
	TxStatusRequested TxStatus = "requested"
	// TxStatusPending — transaction submitted on-chain, awaiting confirmation.
	TxStatusPending TxStatus = "pending"
	// TxStatusSettled — transaction confirmed; final.
	TxStatusSettled TxStatus = "settled"
	// TxStatusReverted — on-chain execution reverted; final.
	TxStatusReverted TxStatus = "reverted"
	// TxStatusIgnored — request superseded or de-duplicated; final.
	TxStatusIgnored TxStatus = "ignored"
)

// Valid reports whether the receiver is one of the defined statuses.
func (s TxStatus) Valid() bool {
	switch s {
	case TxStatusRequested, TxStatusPending, TxStatusSettled, TxStatusReverted, TxStatusIgnored:
		return true
	default:
		return false
	}
}

// Terminal reports whether the status is final.
func (s TxStatus) Terminal() bool {
	switch s {
	case TxStatusSettled, TxStatusReverted, TxStatusIgnored:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (s TxStatus) Validate() error {
	if s.Valid() {
		return nil
	}
	return invalid("TxStatus", string(s))
}

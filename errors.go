// Package derive — error types and sentinel values for the SDK. All
// errors are constructed so they work with errors.Is and errors.As from the
// standard library.
package derive

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/amiwrpremium/go-derive/internal/transport"
)

// APIError is a server-side JSON-RPC error returned by the Derive API. Use
// errors.Is to compare against a sentinel:
//
//	if errors.Is(err, derrors.ErrRateLimited) { ... }
//
// Use errors.As to inspect the raw code and data:
//
//	var apiErr *derrors.APIError
//	if errors.As(err, &apiErr) {
//	    log.Printf("derive code %d: %s", apiErr.Code, apiErr.Message)
//	}
type APIError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// Error implements the error interface. When the server's Message field is
// empty, the canonical description for the code (if known) is substituted —
// see [Description].
func (e *APIError) Error() string {
	msg := e.Message
	if msg == "" {
		msg = Description(e.Code)
	}
	if len(e.Data) > 0 {
		return fmt.Sprintf("derive: api error %d: %s (%s)", e.Code, msg, string(e.Data))
	}
	return fmt.Sprintf("derive: api error %d: %s", e.Code, msg)
}

// CanonicalMessage returns the canonical description for this error's Code.
// It is empty if the code is unknown to this SDK version.
func (e *APIError) CanonicalMessage() string {
	return Description(e.Code)
}

// Is reports whether this APIError corresponds to one of the package's
// sentinel errors. It enables `errors.Is(err, ErrXxx)` for the most
// common categories. Codes that don't map here can still be inspected via
// errors.As + APIError.Code.
//
//nolint:gocyclo // wide switch is the right shape for an exhaustive code map
func (e *APIError) Is(target error) bool {
	switch target {
	case ErrRateLimited:
		return e.Code == CodeRateLimitExceeded || e.Code == CodeConcurrentWSClientsLimitExceeded
	case ErrUnauthorized:
		switch e.Code {
		case CodeInvalidSignature,
			CodeHeaderWalletMissing,
			CodeHeaderWalletMismatch,
			CodeSignerNotOwner,
			CodeMissingPrivateParam,
			CodeSessionKeyNotFound,
			CodeSessionKeyExpired,
			CodeSessionKeyIPNotWhitelisted,
			CodeUnauthorizedKeyScope,
			CodeScopeNotAdmin,
			CodeUnauthorizedRFQMaker,
			CodeAccountNotWhitelistedAtomicOrders,
			CodeSentinelAuthInvalid:
			return true
		}
	case ErrInvalidSignature:
		return e.Code == CodeInvalidSignature
	case ErrSessionKeyExpired:
		return e.Code == CodeSessionKeyExpired
	case ErrSessionKeyNotFound:
		return e.Code == CodeSessionKeyNotFound
	case ErrInsufficientFunds:
		return e.Code == CodeInsufficientFundsOrder ||
			e.Code == CodeERC20InsufficientAllowance ||
			e.Code == CodeERC20InsufficientBalance
	case ErrOrderNotFound:
		return e.Code == CodeOrderNotExist
	case ErrAlreadyCancelled:
		return e.Code == CodeAlreadyCancelled
	case ErrAlreadyFilled:
		return e.Code == CodeAlreadyFilled
	case ErrAlreadyExpired:
		return e.Code == CodeAlreadyExpired
	case ErrInstrumentNotFound:
		return e.Code == CodeInstrumentNotFound || e.Code == CodeAssetNotFound
	case ErrSubaccountNotFound:
		return e.Code == CodeSubaccountNotFound
	case ErrAccountNotFound:
		return e.Code == CodeAccountNotFound
	case ErrChainIDMismatch:
		return e.Code == CodeChainIDMismatch
	case ErrMMPFrozen:
		return e.Code == CodeMMPFrozen
	case ErrRestrictedRegion:
		return e.Code == CodeRestrictedRegion ||
			e.Code == CodeAccountDisabledCompliance
	}
	return false
}

// Code* constants enumerate every JSON-RPC error code Derive returns.
//
// The list mirrors the canonical enum at Derive's published reference
// (https://docs.xyz/reference/error-codes) and the Python SDK's
// data_types/enums.py. Codes are grouped by topic and laid out in numeric
// order so it's easy to find a code starting from a server response.
//
// Compare with [APIError.Code] (or with the helpers in this package such as
// [APIError.Is]) to react to specific failures without string-matching the
// human message.
//
// # When to use a Code* constant directly
//
// Most callers should compare against the sentinels in errors.go via
// errors.Is — they handle several related codes in one check. Reach for a
// Code* constant when you need to disambiguate (e.g. to distinguish
// [CodePostOnlyReject] from a generic [CodeInsufficientFundsOrder]) or when
// you want to surface the canonical description ([Description]).
const (
	// CodeNoError is reserved for a successful response that nonetheless
	// embeds a JSON-RPC error object — rare, but documented by Derive.
	CodeNoError = 0

	// CodeParseError indicates the server could not parse the JSON it received.
	CodeParseError = -32700
	// CodeInvalidRequest indicates the JSON parsed but is not a valid
	// JSON-RPC request object.
	CodeInvalidRequest = -32600
	// CodeMethodNotFound indicates the server does not implement the
	// requested method name.
	CodeMethodNotFound = -32601
	// CodeInvalidParams indicates the method parameters are malformed.
	CodeInvalidParams = -32602
	// CodeInternalError indicates an unexpected server-side failure that
	// is not categorised as one of the more specific codes.
	CodeInternalError = -32603

	// CodeRateLimitExceeded indicates the per-IP request rate limit was
	// exceeded. Default sustained rate is 10 TPS with a 5x burst.
	CodeRateLimitExceeded = -32000
	// CodeConcurrentWSClientsLimitExceeded indicates the connection cap
	// for an account's concurrent WebSocket clients has been exceeded.
	CodeConcurrentWSClientsLimitExceeded = -32100

	// CodeDatabaseError indicates a generic database failure on the venue side.
	CodeDatabaseError = 8000
	// CodeDjangoError indicates a generic application-server failure.
	CodeDjangoError = 8001
	// CodeDuplicateCashAsset means the engine encountered a duplicate cash
	// asset registration.
	CodeDuplicateCashAsset = 8002
	// CodeNoOptionBalanceForSettlement means an option settlement event
	// referenced a balance that does not exist.
	CodeNoOptionBalanceForSettlement = 8003
	// CodeMultipleOptionBalancesForSettlement means an option settlement
	// event matched more than one balance row.
	CodeMultipleOptionBalancesForSettlement = 8004
	// CodeNoVacantInstruments means the venue has no vacant instrument slot
	// for the requested allocation.
	CodeNoVacantInstruments = 8100
	// CodeInvalidServiceType means a configured service-type identifier is
	// not recognised by the engine.
	CodeInvalidServiceType = 8101
	// CodeLatchNotRetained means an internal latch could not be retained
	// (transient locking issue).
	CodeLatchNotRetained = 8102
	// CodeFeedsNotFound means no oracle feeds exist for the requested
	// currency or asset.
	CodeFeedsNotFound = 8200
	// CodeScheduledDeactivationTooLate means an instrument's scheduled
	// deactivation is past the allowed window.
	CodeScheduledDeactivationTooLate = 8300
	// CodeInvalidHeartbeatInterval means an engine or publisher heartbeat
	// interval is invalid.
	CodeInvalidHeartbeatInterval = 8301
	// CodeInvalidMakerOrTakerFee means the maker and/or taker fees are
	// outside the accepted range.
	CodeInvalidMakerOrTakerFee = 8303
	// CodeInvalidInstrumentName means the supplied instrument name does
	// not parse against Derive's naming rules.
	CodeInvalidInstrumentName = 8304
	// CodeOptionSettlementPriceCouldNotBeSaved means an option settlement
	// price could not be persisted.
	CodeOptionSettlementPriceCouldNotBeSaved = 8402
	// CodeOptionSettlementPriceForNonOption means a settlement price was
	// pushed against a non-option asset.
	CodeOptionSettlementPriceForNonOption = 8403
	// CodeCounterpartyInsufficientFunds means a counterparty in an RFQ /
	// quote flow lacks the funds to settle.
	CodeCounterpartyInsufficientFunds = 8500
	// CodeCounterpartyMaxFeeTooLow means at least one counterparty's
	// signed `max_fee` is below the venue minimum.
	CodeCounterpartyMaxFeeTooLow = 8501

	// CodeOrderConfirmationTimeout means the engine did not acknowledge an
	// order within the expected window.
	CodeOrderConfirmationTimeout = 9000
	// CodeEngineConfirmationTimeout means the matching engine did not
	// confirm a state change within the expected window.
	CodeEngineConfirmationTimeout = 9001

	// CodeManagerNotFound means no account manager exists for the requested
	// wallet — a sign the wallet has never deposited.
	CodeManagerNotFound = 10000
	// CodeAssetNotERC20 means the supplied asset address is not an ERC-20.
	CodeAssetNotERC20 = 10001
	// CodeWalletMismatch means the wallet in the request does not match
	// the authenticated wallet (e.g. wrong X-LyraWallet header).
	CodeWalletMismatch = 10002
	// CodeSubaccountMismatch means the subaccount id does not belong to
	// the authenticated wallet.
	CodeSubaccountMismatch = 10003
	// CodeMultipleCurrenciesNotSupported means the operation rejected an
	// attempt to mix currencies in a single request.
	CodeMultipleCurrenciesNotSupported = 10004
	// CodeMaxSubaccountsReached means the wallet hit its subaccount cap.
	CodeMaxSubaccountsReached = 10005
	// CodeMaxSessionKeysReached means the wallet hit its session-key cap.
	CodeMaxSessionKeysReached = 10006
	// CodeMaxAssetsPerSubaccount means the subaccount hit its per-asset cap.
	CodeMaxAssetsPerSubaccount = 10007
	// CodeMaxExpiriesPerSubaccount means the subaccount hit its
	// per-expiry option cap.
	CodeMaxExpiriesPerSubaccount = 10008
	// CodeInvalidRecipientSubaccountID means a transfer's recipient
	// subaccount id is invalid.
	CodeInvalidRecipientSubaccountID = 10009
	// CodePMRMUSDCOnlyCollateral means a portfolio-margin-risk-managed
	// subaccount tried to use non-USDC collateral.
	CodePMRMUSDCOnlyCollateral = 10010
	// CodeERC20InsufficientAllowance means an on-chain ERC-20 allowance
	// was insufficient for the operation.
	CodeERC20InsufficientAllowance = 10011
	// CodeERC20InsufficientBalance means an on-chain ERC-20 balance was
	// insufficient for the operation.
	CodeERC20InsufficientBalance = 10012
	// CodePendingDeposit means a deposit has not yet finalized and the
	// funds aren't usable.
	CodePendingDeposit = 10013
	// CodePendingWithdrawal means a withdrawal is in flight.
	CodePendingWithdrawal = 10014
	// CodePM2CollateralConstraint means a portfolio-margin v2 collateral
	// rule was violated.
	CodePM2CollateralConstraint = 10015

	// CodeInsufficientFundsOrder means the order would have breached the
	// subaccount's margin / funds requirements.
	CodeInsufficientFundsOrder = 11000
	// CodeOrderRejectedFromQueue means the matching queue rejected the
	// order (typically a transient capacity issue).
	CodeOrderRejectedFromQueue = 11002
	// CodeAlreadyCancelled means the order is already in a cancelled state.
	CodeAlreadyCancelled = 11003
	// CodeAlreadyFilled means the order has already filled in full.
	CodeAlreadyFilled = 11004
	// CodeAlreadyExpired means the order's signature expiry already passed.
	CodeAlreadyExpired = 11005
	// CodeOrderNotExist means no order exists for the supplied id.
	CodeOrderNotExist = 11006
	// CodeSelfCrossDisallowed means the order would have crossed against
	// another order from the same subaccount.
	CodeSelfCrossDisallowed = 11007
	// CodePostOnlyReject means a post-only order would have crossed the
	// book and was rejected.
	CodePostOnlyReject = 11008
	// CodeZeroLiquidity means an IOC/market order found nothing on the
	// other side of the book.
	CodeZeroLiquidity = 11009
	// CodePostOnlyInvalidType means post-only is incompatible with the
	// requested order type (e.g. market).
	CodePostOnlyInvalidType = 11010
	// CodeInvalidSignatureExpiry means the signature_expiry_sec is too
	// short, in the past, or out of the accepted range.
	CodeInvalidSignatureExpiry = 11011
	// CodeInvalidAmount means the amount is zero, negative, or not
	// a multiple of the instrument's amount_step.
	CodeInvalidAmount = 11012
	// CodeInvalidLimitPrice means the limit price is off the tick size,
	// negative, or outside the accepted band.
	CodeInvalidLimitPrice = 11013
	// CodeFOKNotFilled means a fill-or-kill order could not be fully filled
	// at submission time and was cancelled.
	CodeFOKNotFilled = 11014
	// CodeMMPFrozen means market-maker protection has tripped for this
	// currency and orders are temporarily blocked. See [MMPConfig] in
	// (see methods.go).
	CodeMMPFrozen = 11015
	// CodeAlreadyConsumed means the engine has already processed this
	// nonce / request id.
	CodeAlreadyConsumed = 11016
	// CodeNonUniqueNonce means the supplied nonce has been used before
	// for this subaccount.
	CodeNonUniqueNonce = 11017
	// CodeInvalidNonceDate means the nonce timestamp is outside the
	// accepted window.
	CodeInvalidNonceDate = 11018
	// CodeOpenOrdersLimitExceeded means the per-subaccount open-orders cap
	// was exceeded.
	CodeOpenOrdersLimitExceeded = 11019
	// CodeNegativeERC20Balance means the operation would push an ERC-20
	// balance below zero.
	CodeNegativeERC20Balance = 11020
	// CodeInstrumentNotLive means the instrument is delisted, paused, or
	// pre-launch and not currently accepting orders.
	CodeInstrumentNotLive = 11021
	// CodeRejectTimestampExceeded means the request's reject_timestamp
	// passed before the engine could process it.
	CodeRejectTimestampExceeded = 11022
	// CodeMaxFeeTooLow means the order's max_fee is below the engine's
	// current minimum.
	CodeMaxFeeTooLow = 11023
	// CodeReduceOnlyNotSupported means reduce_only is incompatible with
	// the requested time-in-force.
	CodeReduceOnlyNotSupported = 11024
	// CodeReduceOnlyReject means a reduce-only order was rejected because
	// it would have grown or flipped the position.
	CodeReduceOnlyReject = 11025
	// CodeTransferReject means a sub-account transfer was rejected by
	// the engine.
	CodeTransferReject = 11026
	// CodeSubaccountUnderLiquidation means the subaccount is undergoing
	// a liquidation auction and cannot accept new orders.
	CodeSubaccountUnderLiquidation = 11027
	// CodeReplaceFilledAmountMismatch means the replaced order's filled
	// amount does not match the expected state.
	CodeReplaceFilledAmountMismatch = 11028

	// CodeTriggerOrderCancelled means a trigger order was cancelled before
	// its trigger fired.
	CodeTriggerOrderCancelled = 11050
	// CodeInvalidTriggerPrice means the trigger price is malformed.
	CodeInvalidTriggerPrice = 11051
	// CodeTriggerOrderLimitExceeded means the per-subaccount trigger-orders
	// cap was exceeded.
	CodeTriggerOrderLimitExceeded = 11052
	// CodeTriggerPriceTypeUnsupported means the chosen price type
	// (mark/index/last) is not supported on this instrument.
	CodeTriggerPriceTypeUnsupported = 11053
	// CodeTriggerOrderReplaceUnsupported means replace is not supported
	// for trigger orders.
	CodeTriggerOrderReplaceUnsupported = 11054
	// CodeMarketOrderInvalidTriggerPrice means a market order specified
	// an invalid trigger price.
	CodeMarketOrderInvalidTriggerPrice = 11055

	// CodeLegInstrumentsNotUnique means the legs of an RFQ or quote
	// contain duplicate instruments.
	CodeLegInstrumentsNotUnique = 11100
	// CodeRFQNotFound means no RFQ exists for the supplied id.
	CodeRFQNotFound = 11101
	// CodeQuoteNotFound means no quote exists for the supplied id.
	CodeQuoteNotFound = 11102
	// CodeRFQLegMismatch means a quote's legs do not match the RFQ's legs.
	CodeRFQLegMismatch = 11103
	// CodeRFQNotOpen means the RFQ is not in a quotable / executable state.
	CodeRFQNotOpen = 11104
	// CodeRFQIDMismatch means the rfq_id in the quote does not match the
	// referenced RFQ.
	CodeRFQIDMismatch = 11105
	// CodeInvalidRFQCounterparty means the counterparty is not authorized
	// to participate in this RFQ.
	CodeInvalidRFQCounterparty = 11106
	// CodeQuoteCostTooHigh means the quote's total cost exceeds the
	// caller's max_total_fee.
	CodeQuoteCostTooHigh = 11107

	// CodeAuctionNotOngoing means an auction-only operation was attempted
	// outside an active auction window.
	CodeAuctionNotOngoing = 11200
	// CodeOpenOrdersNotAllowed means orders cannot be placed during this
	// auction phase.
	CodeOpenOrdersNotAllowed = 11201
	// CodePriceLimitExceeded means the price exceeds the auction's allowed band.
	CodePriceLimitExceeded = 11202
	// CodeLastTradeIDMismatch means the supplied last_trade_id does not
	// match the auction's current state.
	CodeLastTradeIDMismatch = 11203

	// CodeAssetNotFound means no asset exists for the supplied identifier.
	CodeAssetNotFound = 12000
	// CodeInstrumentNotFound means no instrument exists for the supplied name.
	CodeInstrumentNotFound = 12001
	// CodeCurrencyNotFound means no currency exists for the supplied symbol.
	CodeCurrencyNotFound = 12002
	// CodeUSDCNoCaps means USDC caps are not configured for this account.
	CodeUSDCNoCaps = 12003

	// CodeInvalidChannels means one or more channel names in the subscribe
	// request are unrecognised.
	CodeInvalidChannels = 13000

	// CodeAccountNotFound means the wallet has no Derive account.
	CodeAccountNotFound = 14000
	// CodeSubaccountNotFound means the supplied subaccount id is unknown.
	CodeSubaccountNotFound = 14001
	// CodeSubaccountWithdrawn means the subaccount has been fully
	// withdrawn and is closed.
	CodeSubaccountWithdrawn = 14002
	// CodeUseDeregisterSessionKey means the caller tried to reduce a
	// session key's expiry via `register_session_key`; use the dedicated
	// deregister path instead.
	CodeUseDeregisterSessionKey = 14008
	// CodeSessionKeyExpiryTooLow means the session key expiry is below the
	// minimum required by the engine.
	CodeSessionKeyExpiryTooLow = 14009
	// CodeSessionKeyAlreadyRegistered means the session key already exists.
	CodeSessionKeyAlreadyRegistered = 14010
	// CodeSessionKeyRegisteredOtherAccount means the session key is
	// registered to a different account.
	CodeSessionKeyRegisteredOtherAccount = 14011
	// CodeAddressNotChecksummed means an address must be EIP-55 checksummed.
	CodeAddressNotChecksummed = 14012
	// CodeInvalidEthAddress means an address is not a valid Ethereum address.
	CodeInvalidEthAddress = 14013
	// CodeInvalidSignature means the signature did not verify.
	CodeInvalidSignature = 14014
	// CodeNonceMismatch means the request's nonce does not match the
	// on-chain account state.
	CodeNonceMismatch = 14015
	// CodeRawTxFunctionMismatch means the raw transaction function selector
	// does not match the expected one.
	CodeRawTxFunctionMismatch = 14016
	// CodeRawTxContractMismatch means the raw transaction's target contract
	// does not match the expected one.
	CodeRawTxContractMismatch = 14017
	// CodeRawTxParamsMismatch means the raw transaction's parameter shape
	// does not match.
	CodeRawTxParamsMismatch = 14018
	// CodeRawTxParamValuesMismatch means the raw transaction's parameter
	// values do not match.
	CodeRawTxParamValuesMismatch = 14019
	// CodeHeaderWalletMismatch means the X-LyraWallet header does not match
	// the request signer.
	CodeHeaderWalletMismatch = 14020
	// CodeHeaderWalletMissing means the X-LyraWallet header is absent.
	CodeHeaderWalletMissing = 14021
	// CodePrivateChannelSubscriptionFailed means a private channel
	// subscription failed and the WS connection must re-login.
	CodePrivateChannelSubscriptionFailed = 14022
	// CodeSignerNotOwner means the signer is neither the owner nor a
	// registered session key.
	CodeSignerNotOwner = 14023
	// CodeChainIDMismatch means the chain id in the signed action does not
	// match the network's chain id (signing against the wrong env).
	CodeChainIDMismatch = 14024
	// CodeMissingPrivateParam means a required authentication parameter is
	// absent on a private endpoint.
	CodeMissingPrivateParam = 14025
	// CodeSessionKeyNotFound means the session key is unknown to the engine.
	CodeSessionKeyNotFound = 14026
	// CodeUnauthorizedRFQMaker means the maker is not authorized to quote
	// for this counterparty.
	CodeUnauthorizedRFQMaker = 14027
	// CodeCrossCurrencyRFQNotSupported means cross-currency RFQs are not
	// available.
	CodeCrossCurrencyRFQNotSupported = 14028
	// CodeSessionKeyIPNotWhitelisted means the request originated from an
	// IP not whitelisted for this session key.
	CodeSessionKeyIPNotWhitelisted = 14029
	// CodeSessionKeyExpired means the session key's expiry has passed.
	CodeSessionKeyExpired = 14030
	// CodeUnauthorizedKeyScope means the session key's scope does not
	// authorise the requested operation.
	CodeUnauthorizedKeyScope = 14031
	// CodeScopeNotAdmin means the operation requires admin-scoped credentials.
	CodeScopeNotAdmin = 14032
	// CodeAccountNotWhitelistedAtomicOrders means the account is not
	// enrolled in atomic-signing orders.
	CodeAccountNotWhitelistedAtomicOrders = 14033
	// CodeReferralCodeNotFound means the supplied referral code is unknown.
	CodeReferralCodeNotFound = 14034

	// CodeRestrictedRegion means the request originates from a region
	// blocked by Derive's geo policy.
	CodeRestrictedRegion = 16000
	// CodeAccountDisabledCompliance means the account has been disabled
	// for compliance reasons.
	CodeAccountDisabledCompliance = 16001
	// CodeSentinelAuthInvalid means a Sentinel-class authentication check
	// failed.
	CodeSentinelAuthInvalid = 16100

	// CodeInvalidBlockNumber means the block number is invalid.
	CodeInvalidBlockNumber = 18000
	// CodeBlockEstimationFailed means block estimation failed for a vault
	// operation.
	CodeBlockEstimationFailed = 18001
	// CodeLightAccountOwnerMismatch means a light-account owner check failed.
	CodeLightAccountOwnerMismatch = 18002
	// CodeVaultERC20AssetNotExists means the vault has no ERC-20 asset
	// matching the request.
	CodeVaultERC20AssetNotExists = 18003
	// CodeVaultERC20PoolNotExists means the vault has no ERC-20 pool
	// matching the request.
	CodeVaultERC20PoolNotExists = 18004
	// CodeVaultAddAssetBeforeBalance means a vault asset must be added
	// before its balance can be updated.
	CodeVaultAddAssetBeforeBalance = 18005
	// CodeInvalidSwellSeason means the supplied Swell season is invalid.
	CodeInvalidSwellSeason = 18006
	// CodeVaultNotFound means no vault exists for the supplied identifier.
	CodeVaultNotFound = 18007

	// CodeMakerProgramNotFound means no maker program exists for the
	// supplied identifier.
	CodeMakerProgramNotFound = 19000
)

// Re-export stdlib helpers so callers don't need a second import.
var (
	Is     = errors.Is
	As     = errors.As
	Unwrap = errors.Unwrap
	New    = errors.New
)

// Sentinel errors. Compare with errors.Is. Each one maps to a category of
// JSON-RPC codes — see [APIError.Is] in api.go for the mapping.
//
//	if errors.Is(err, derrors.ErrRateLimited) { backoff(); return }
//	if errors.Is(err, derrors.ErrSessionKeyExpired) { reAuth(); return }
var (
	// ErrNotConnected is returned when a WebSocket call is attempted before
	// Connect() has succeeded or after the connection has terminated.
	// Declared in `internal/transport` to break the root↔transport cycle;
	// referenced here so `errors.Is(err, derrors.ErrNotConnected)` still
	// matches values produced by the transport pumps (same pointer).
	ErrNotConnected = transport.ErrNotConnected

	// ErrAlreadyConnected is returned when Connect() is called on a client
	// that is already connected. See ErrNotConnected for the indirection
	// rationale.
	ErrAlreadyConnected = transport.ErrAlreadyConnected

	// ErrUnauthorized is returned when the SDK has no signer configured or
	// the server rejects an authentication-class error code (invalid
	// signature, missing wallet header, expired session key, scope refusal).
	ErrUnauthorized = errors.New("derive: unauthorized")

	// ErrInvalidSignature maps specifically to code 14014 — the request was
	// rejected because the signature did not verify against the typed data.
	ErrInvalidSignature = errors.New("derive: invalid signature")

	// ErrSessionKeyExpired maps to code 14030.
	ErrSessionKeyExpired = errors.New("derive: session key expired")

	// ErrSessionKeyNotFound maps to code 14026.
	ErrSessionKeyNotFound = errors.New("derive: session key not found")

	// ErrRateLimited covers both per-IP request rate limiting (-32000) and
	// the WebSocket-concurrency cap (-32100).
	ErrRateLimited = errors.New("derive: rate limited")

	// ErrInsufficientFunds covers the order-side margin/funds rejection
	// (11000) and ERC-20 balance issues (10011, 10012).
	ErrInsufficientFunds = errors.New("derive: insufficient funds")

	// ErrOrderNotFound maps to 11006.
	ErrOrderNotFound = errors.New("derive: order not found")

	// ErrAlreadyCancelled, ErrAlreadyFilled, ErrAlreadyExpired correspond to
	// 11003, 11004, 11005.
	ErrAlreadyCancelled = errors.New("derive: order already cancelled")
	ErrAlreadyFilled    = errors.New("derive: order already filled")
	ErrAlreadyExpired   = errors.New("derive: order already expired")

	// ErrInstrumentNotFound covers 12001 (instrument) and 12000 (asset).
	ErrInstrumentNotFound = errors.New("derive: instrument not found")

	// ErrSubaccountNotFound maps to 14001.
	ErrSubaccountNotFound = errors.New("derive: subaccount not found")

	// ErrAccountNotFound maps to 14000.
	ErrAccountNotFound = errors.New("derive: account not found")

	// ErrChainIDMismatch maps to 14024 — signer's chain id doesn't match
	// the network's (almost always means signing for the wrong env).
	ErrChainIDMismatch = errors.New("derive: chain id mismatch")

	// ErrMMPFrozen maps to 11015 — market-maker protection has tripped.
	ErrMMPFrozen = errors.New("derive: market-maker protection frozen")

	// ErrRestrictedRegion maps to 16000 / 16001 / 16100 (compliance class).
	ErrRestrictedRegion = errors.New("derive: restricted region")

	// ErrSubscriptionClosed is returned by Subscription.Updates() once the
	// channel has been closed by either party. See ErrNotConnected for the
	// indirection rationale.
	ErrSubscriptionClosed = transport.ErrSubscriptionClosed

	// ErrSubaccountRequired is returned for private calls that need a
	// subaccount ID configured on the client.
	ErrSubaccountRequired = errors.New("derive: subaccount id required")

	// ErrInvalidConfig is returned by NewClient on malformed options.
	ErrInvalidConfig = errors.New("derive: invalid configuration")
)

// codeMessages maps each Derive error code to a canonical human-readable
// description. The text is derived from Derive's published error code
// reference at https://docs.xyz/reference/error-codes and the
// Python SDK's enum comments.
//
// Server responses always carry their own Message field, but it is sometimes
// empty or terse; [Description] returns the canonical text so callers and
// log lines have a stable, readable description regardless.
var codeMessages = map[int]string{
	CodeNoError: "no error",

	CodeParseError:                       "parse error: invalid JSON received by the server",
	CodeInvalidRequest:                   "invalid request: the JSON sent is not a valid request object",
	CodeMethodNotFound:                   "method not found: the requested method does not exist",
	CodeInvalidParams:                    "invalid params: the method parameters are malformed",
	CodeInternalError:                    "internal error: an unexpected server error occurred",
	CodeRateLimitExceeded:                "rate limit exceeded: too many requests in the rate-limit window",
	CodeConcurrentWSClientsLimitExceeded: "concurrent WebSocket clients limit exceeded for this account",

	CodeDatabaseError:                        "database error",
	CodeDjangoError:                          "application server error",
	CodeDuplicateCashAsset:                   "duplicate cash asset",
	CodeNoOptionBalanceForSettlement:         "no open option balance for the settlement event was found",
	CodeMultipleOptionBalancesForSettlement:  "more than one option balance for the settlement event was found",
	CodeNoVacantInstruments:                  "no vacant instruments",
	CodeInvalidServiceType:                   "invalid service type",
	CodeLatchNotRetained:                     "latch not retained",
	CodeFeedsNotFound:                        "feeds not found",
	CodeScheduledDeactivationTooLate:         "scheduled deactivation too late",
	CodeInvalidHeartbeatInterval:             "engine or publisher heartbeat interval is invalid",
	CodeInvalidMakerOrTakerFee:               "the maker and/or taker fees are invalid",
	CodeInvalidInstrumentName:                "instrument name is invalid",
	CodeOptionSettlementPriceCouldNotBeSaved: "option settlement price could not be saved",
	CodeOptionSettlementPriceForNonOption:    "option settlement price cannot be saved to a non-option asset",
	CodeCounterpartyInsufficientFunds:        "counterparty insufficient funds",
	CodeCounterpartyMaxFeeTooLow:             "max fee for one or more counterparties is too low",

	CodeOrderConfirmationTimeout:  "order confirmation timed out before the engine acknowledged it",
	CodeEngineConfirmationTimeout: "matching-engine confirmation timed out",

	CodeManagerNotFound:                "manager not found for the requested wallet",
	CodeAssetNotERC20:                  "asset is not an ERC-20 contract",
	CodeWalletMismatch:                 "wallet address in the request does not match the authenticated wallet",
	CodeSubaccountMismatch:             "subaccount id does not belong to the authenticated wallet",
	CodeMultipleCurrenciesNotSupported: "operation does not support multiple currencies",
	CodeMaxSubaccountsReached:          "maximum number of subaccounts reached for this wallet",
	CodeMaxSessionKeysReached:          "maximum number of session keys reached for this wallet",
	CodeMaxAssetsPerSubaccount:         "maximum number of assets per subaccount reached",
	CodeMaxExpiriesPerSubaccount:       "maximum number of option expiries per subaccount reached",
	CodeInvalidRecipientSubaccountID:   "recipient subaccount id is invalid",
	CodePMRMUSDCOnlyCollateral:         "portfolio-margin (PMRM) subaccounts only support USDC collateral",
	CodeERC20InsufficientAllowance:     "ERC-20 allowance is insufficient for this operation",
	CodeERC20InsufficientBalance:       "ERC-20 balance is insufficient for this operation",
	CodePendingDeposit:                 "deposit is still pending and cannot be used yet",
	CodePendingWithdrawal:              "withdrawal is still pending",
	CodePM2CollateralConstraint:        "portfolio-margin v2 collateral constraint not satisfied",

	CodeInsufficientFundsOrder:      "insufficient funds: order would breach the subaccount's margin requirements",
	CodeOrderRejectedFromQueue:      "order was rejected from the matching queue",
	CodeAlreadyCancelled:            "order is already cancelled",
	CodeAlreadyFilled:               "order is already filled",
	CodeAlreadyExpired:              "order is already expired",
	CodeOrderNotExist:               "order does not exist",
	CodeSelfCrossDisallowed:         "self-cross is disallowed",
	CodePostOnlyReject:              "post-only order would have crossed the book and was rejected",
	CodeZeroLiquidity:               "no liquidity is currently available to match this order",
	CodePostOnlyInvalidType:         "post-only is only valid for limit orders",
	CodeInvalidSignatureExpiry:      "signature expiry is invalid (too short or out of range)",
	CodeInvalidAmount:               "order amount is invalid (zero, negative, or below minimum step)",
	CodeInvalidLimitPrice:           "limit price is invalid (off tick size, negative, or out of band)",
	CodeFOKNotFilled:                "fill-or-kill order could not be fully filled and was cancelled",
	CodeMMPFrozen:                   "market-maker protection has frozen this currency",
	CodeAlreadyConsumed:             "order has already been consumed by a prior request",
	CodeNonUniqueNonce:              "nonce has already been used; nonces must be unique per subaccount",
	CodeInvalidNonceDate:            "nonce timestamp is outside the accepted window",
	CodeOpenOrdersLimitExceeded:     "open-orders-per-subaccount limit exceeded",
	CodeNegativeERC20Balance:        "operation would leave the ERC-20 balance negative",
	CodeInstrumentNotLive:           "instrument is not currently live (delisted, paused, or pre-launch)",
	CodeRejectTimestampExceeded:     "reject_timestamp exceeded before the engine processed the order",
	CodeMaxFeeTooLow:                "max_fee on the order is below the engine's current minimum",
	CodeReduceOnlyNotSupported:      "reduce_only is incompatible with the requested time-in-force",
	CodeReduceOnlyReject:            "reduce-only order rejected — would have grown or flipped the position",
	CodeTransferReject:              "subaccount transfer rejected by the engine",
	CodeSubaccountUnderLiquidation:  "subaccount is undergoing a liquidation auction",
	CodeReplaceFilledAmountMismatch: "replace order's filled amount does not match the expected state",

	CodeTriggerOrderCancelled:          "trigger order was cancelled before its trigger fired",
	CodeInvalidTriggerPrice:            "trigger price is invalid",
	CodeTriggerOrderLimitExceeded:      "trigger-orders-per-subaccount limit exceeded",
	CodeTriggerPriceTypeUnsupported:    "trigger price type is not supported for this instrument",
	CodeTriggerOrderReplaceUnsupported: "replace is not supported for trigger orders",
	CodeMarketOrderInvalidTriggerPrice: "market order has an invalid trigger price",

	CodeLegInstrumentsNotUnique: "RFQ/quote leg instruments must be unique",
	CodeRFQNotFound:             "RFQ not found",
	CodeQuoteNotFound:           "quote not found",
	CodeRFQLegMismatch:          "quote legs do not match the RFQ legs",
	CodeRFQNotOpen:              "RFQ is not open and cannot be quoted on or executed",
	CodeRFQIDMismatch:           "RFQ id in the request does not match the quote",
	CodeInvalidRFQCounterparty:  "RFQ counterparty is invalid or unauthorized",
	CodeQuoteCostTooHigh:        "quote cost exceeds the configured maximum",

	CodeAuctionNotOngoing:    "auction is not currently ongoing",
	CodeOpenOrdersNotAllowed: "open orders are not allowed during this auction phase",
	CodePriceLimitExceeded:   "price limit exceeded for the auction",
	CodeLastTradeIDMismatch:  "last trade id does not match the latest auction state",

	CodeAssetNotFound:      "asset not found",
	CodeInstrumentNotFound: "instrument not found",
	CodeCurrencyNotFound:   "currency not found",
	CodeUSDCNoCaps:         "USDC has no caps configured",

	CodeInvalidChannels: "invalid subscription channel name(s)",

	CodeAccountNotFound:                   "account not found",
	CodeSubaccountNotFound:                "subaccount not found",
	CodeSubaccountWithdrawn:               "subaccount has been fully withdrawn",
	CodeUseDeregisterSessionKey:           "cannot reduce session-key expiry via register_session_key — use the deregister endpoint",
	CodeSessionKeyExpiryTooLow:            "session key expiry is below the minimum allowed",
	CodeSessionKeyAlreadyRegistered:       "session key is already registered",
	CodeSessionKeyRegisteredOtherAccount:  "session key is already registered to a different account",
	CodeAddressNotChecksummed:             "address is not EIP-55 checksummed",
	CodeInvalidEthAddress:                 "address is not a valid Ethereum address",
	CodeInvalidSignature:                  "signature did not verify",
	CodeNonceMismatch:                     "nonce mismatch with the on-chain account state",
	CodeRawTxFunctionMismatch:             "raw transaction function selector does not match the expected one",
	CodeRawTxContractMismatch:             "raw transaction target contract does not match the expected one",
	CodeRawTxParamsMismatch:               "raw transaction parameter shape does not match",
	CodeRawTxParamValuesMismatch:          "raw transaction parameter values do not match",
	CodeHeaderWalletMismatch:              "X-LyraWallet header does not match the request signer",
	CodeHeaderWalletMissing:               "X-LyraWallet header is missing",
	CodePrivateChannelSubscriptionFailed:  "private-channel subscription failed: re-login required",
	CodeSignerNotOwner:                    "signer is not the registered owner or session key",
	CodeChainIDMismatch:                   "chain id in the signed action does not match the network",
	CodeMissingPrivateParam:               "private endpoint called without a required authentication parameter",
	CodeSessionKeyNotFound:                "session key not found",
	CodeUnauthorizedRFQMaker:              "unauthorized RFQ maker for this counterparty",
	CodeCrossCurrencyRFQNotSupported:      "cross-currency RFQ is not supported",
	CodeSessionKeyIPNotWhitelisted:        "session key's source IP is not whitelisted",
	CodeSessionKeyExpired:                 "session key has expired",
	CodeUnauthorizedKeyScope:              "session key scope does not authorize this operation",
	CodeScopeNotAdmin:                     "operation requires admin scope",
	CodeAccountNotWhitelistedAtomicOrders: "account is not whitelisted for atomic-signing orders",
	CodeReferralCodeNotFound:              "referral code not found",

	CodeRestrictedRegion:          "request originates from a restricted region",
	CodeAccountDisabledCompliance: "account has been disabled for compliance reasons",
	CodeSentinelAuthInvalid:       "sentinel authentication is invalid",

	CodeInvalidBlockNumber:         "block number is invalid",
	CodeBlockEstimationFailed:      "block estimation failed",
	CodeLightAccountOwnerMismatch:  "light-account owner does not match the request",
	CodeVaultERC20AssetNotExists:   "vault ERC-20 asset does not exist",
	CodeVaultERC20PoolNotExists:    "vault ERC-20 pool does not exist",
	CodeVaultAddAssetBeforeBalance: "vault asset must be added before balance updates",
	CodeInvalidSwellSeason:         "Swell season is invalid",
	CodeVaultNotFound:              "vault not found",

	CodeMakerProgramNotFound: "maker program not found",
}

// Description returns the canonical description for a Derive error code, or
// an empty string if the code is unknown to this version of the SDK.
//
// The text mirrors Derive's published error reference; it is appropriate to
// surface in user-facing messages and structured logs.
func Description(code int) string {
	return codeMessages[code]
}

// HasDescription reports whether the SDK knows a canonical description for
// the given code.
func HasDescription(code int) bool {
	_, ok := codeMessages[code]
	return ok
}

// ConnectionError wraps low-level transport failures so callers can
// distinguish network problems from API errors. Declared in
// `internal/transport` so the transport package can construct values
// without a `root → transport → root` import cycle; aliased here so
// external code keeps using `ConnectionError` exactly as before.
type ConnectionError = transport.ConnectionError

// TimeoutError signals a deadline expiry while waiting for a response.
// Declared in `internal/transport` for the same reason as
// [ConnectionError]; aliased here.
type TimeoutError = transport.TimeoutError

// SigningError wraps failures inside the signer (key parsing, hashing, ECDSA).
type SigningError struct {
	Op  string
	Err error
}

// Error implements the error interface.
func (e *SigningError) Error() string {
	return fmt.Sprintf("derive: signer: %s: %v", e.Op, e.Err)
}

// Unwrap returns the underlying signer error.
func (e *SigningError) Unwrap() error { return e.Err }

// ExpiredSignatureError indicates a signed action's expiry has passed before
// the server received it. Lengthen the expiry window or check the system
// clock skew.
type ExpiredSignatureError struct {
	ExpiryUnixSec int64
	NowUnixSec    int64
}

// Error implements the error interface.
func (e *ExpiredSignatureError) Error() string {
	return fmt.Sprintf("derive: signature expired at %d (now=%d)", e.ExpiryUnixSec, e.NowUnixSec)
}

// Package errors provides the SDK's error types and sentinel values. All
// errors are constructed so they work with errors.Is and errors.As from the
// standard library.
package errors

// Code* constants enumerate every JSON-RPC error code Derive returns.
//
// The list mirrors the canonical enum at Derive's published reference
// (https://docs.derive.xyz/reference/error-codes) and the Python SDK's
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

	// --- Standard JSON-RPC 2.0 (negative range) ---------------------------

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

	// --- Internal / database (8xxx) ---------------------------------------

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

	// --- Engine / order-confirmation timeouts (9xxx) ----------------------

	// CodeOrderConfirmationTimeout means the engine did not acknowledge an
	// order within the expected window.
	CodeOrderConfirmationTimeout = 9000
	// CodeEngineConfirmationTimeout means the matching engine did not
	// confirm a state change within the expected window.
	CodeEngineConfirmationTimeout = 9001
	// CodeCacheConnectionError means the engine's cache layer is
	// unreachable; the request body may include further diagnostic
	// details under `data`.
	CodeCacheConnectionError = 9002

	// --- Account / wallet (10xxx) -----------------------------------------

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
	// CodePM2CollateralUnsupported means PortfolioMargin v2 does not
	// support the requested collateral asset. Use a Portfolio Manager
	// or Standard Manager subaccount that does support the currency.
	CodePM2CollateralUnsupported = 10016

	// --- Order placement / lifecycle (11xxx) ------------------------------

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
	// internal/methods.
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
	// CodeOpenInterestCapExceeded means a trade or transfer was
	// rejected because it would have caused the open-interest to
	// exceed the cap for either the sender or recipient manager.
	CodeOpenInterestCapExceeded = 11029

	// --- Trigger orders (1105x) -------------------------------------------

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

	// --- RFQ / Quote (111xx) ----------------------------------------------

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
	// CodeRFQPartialFillTooHigh means an attempt was made to fill an
	// RFQ above the requested total size.
	CodeRFQPartialFillTooHigh = 11108
	// CodeRFQFilledDirectionLocked means an RFQ's filled direction
	// cannot be changed once the RFQ is partially filled.
	CodeRFQFilledDirectionLocked = 11109
	// CodeQuoteTakerCostTooHigh means the quote's taker-side total
	// cost exceeded the orderbook execution price band.
	CodeQuoteTakerCostTooHigh = 11110
	// CodeRFQDisabledForAccount means RFQ functionality is disabled
	// for this account, typically because of suspicious activity or
	// a terms-of-service violation.
	CodeRFQDisabledForAccount = 11111

	// --- Auction (112xx) --------------------------------------------------

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

	// --- Assets / instruments (12xxx) -------------------------------------

	// CodeAssetNotFound means no asset exists for the supplied identifier.
	CodeAssetNotFound = 12000
	// CodeInstrumentNotFound means no instrument exists for the supplied name.
	CodeInstrumentNotFound = 12001
	// CodeCurrencyNotFound means no currency exists for the supplied symbol.
	CodeCurrencyNotFound = 12002
	// CodeUSDCNoCaps means USDC caps are not configured for this account.
	CodeUSDCNoCaps = 12003

	// --- Subscriptions (13xxx) --------------------------------------------

	// CodeInvalidChannels means one or more channel names in the subscribe
	// request are unrecognised.
	CodeInvalidChannels = 13000

	// --- Account / auth (14xxx) -------------------------------------------

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

	// --- Compliance (16xxx) -----------------------------------------------

	// CodeRestrictedRegion means the request originates from a region
	// blocked by Derive's geo policy.
	CodeRestrictedRegion = 16000
	// CodeAccountDisabledCompliance means the account has been disabled
	// for compliance reasons.
	CodeAccountDisabledCompliance = 16001
	// CodeOFACBlocked means the account is blocked because of OFAC
	// compliance violations.
	CodeOFACBlocked = 16002
	// CodeSentinelAuthInvalid means a Sentinel-class authentication check
	// failed.
	CodeSentinelAuthInvalid = 16100

	// --- Vault / block (18xxx) --------------------------------------------

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

	// --- Maker programs (19xxx) -------------------------------------------

	// CodeMakerProgramNotFound means no maker program exists for the
	// supplied identifier.
	CodeMakerProgramNotFound = 19000
)

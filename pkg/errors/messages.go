// Package errors — see api.go for the overview.
package errors

// codeMessages maps each Derive error code to a canonical human-readable
// description. The text is derived from Derive's published error code
// reference at https://docs.derive.xyz/reference/error-codes and the
// Python SDK's enum comments.
//
// Server responses always carry their own Message field, but it is sometimes
// empty or terse; [Description] returns the canonical text so callers and
// log lines have a stable, readable description regardless.
var codeMessages = map[int]string{
	CodeNoError: "no error",

	// Standard JSON-RPC 2.0
	CodeParseError:                       "parse error: invalid JSON received by the server",
	CodeInvalidRequest:                   "invalid request: the JSON sent is not a valid request object",
	CodeMethodNotFound:                   "method not found: the requested method does not exist",
	CodeInvalidParams:                    "invalid params: the method parameters are malformed",
	CodeInternalError:                    "internal error: an unexpected server error occurred",
	CodeRateLimitExceeded:                "rate limit exceeded: too many requests in the rate-limit window",
	CodeConcurrentWSClientsLimitExceeded: "concurrent WebSocket clients limit exceeded for this account",

	// Internal / database (8xxx)
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

	// Engine / order-confirmation timeouts (9xxx)
	CodeOrderConfirmationTimeout:  "order confirmation timed out before the engine acknowledged it",
	CodeEngineConfirmationTimeout: "matching-engine confirmation timed out",

	// Account / wallet (10xxx)
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

	// Order placement / lifecycle (11xxx)
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

	// Trigger orders (1105x)
	CodeTriggerOrderCancelled:          "trigger order was cancelled before its trigger fired",
	CodeInvalidTriggerPrice:            "trigger price is invalid",
	CodeTriggerOrderLimitExceeded:      "trigger-orders-per-subaccount limit exceeded",
	CodeTriggerPriceTypeUnsupported:    "trigger price type is not supported for this instrument",
	CodeTriggerOrderReplaceUnsupported: "replace is not supported for trigger orders",
	CodeMarketOrderInvalidTriggerPrice: "market order has an invalid trigger price",

	// RFQ / Quote (111xx)
	CodeLegInstrumentsNotUnique: "RFQ/quote leg instruments must be unique",
	CodeRFQNotFound:             "RFQ not found",
	CodeQuoteNotFound:           "quote not found",
	CodeRFQLegMismatch:          "quote legs do not match the RFQ legs",
	CodeRFQNotOpen:              "RFQ is not open and cannot be quoted on or executed",
	CodeRFQIDMismatch:           "RFQ id in the request does not match the quote",
	CodeInvalidRFQCounterparty:  "RFQ counterparty is invalid or unauthorized",
	CodeQuoteCostTooHigh:        "quote cost exceeds the configured maximum",

	// Auction (112xx)
	CodeAuctionNotOngoing:    "auction is not currently ongoing",
	CodeOpenOrdersNotAllowed: "open orders are not allowed during this auction phase",
	CodePriceLimitExceeded:   "price limit exceeded for the auction",
	CodeLastTradeIDMismatch:  "last trade id does not match the latest auction state",

	// Assets / instruments (12xxx)
	CodeAssetNotFound:      "asset not found",
	CodeInstrumentNotFound: "instrument not found",
	CodeCurrencyNotFound:   "currency not found",
	CodeUSDCNoCaps:         "USDC has no caps configured",

	// Subscriptions (13xxx)
	CodeInvalidChannels: "invalid subscription channel name(s)",

	// Account / auth (14xxx)
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

	// Compliance (16xxx)
	CodeRestrictedRegion:          "request originates from a restricted region",
	CodeAccountDisabledCompliance: "account has been disabled for compliance reasons",
	CodeSentinelAuthInvalid:       "sentinel authentication is invalid",

	// Vault / block (18xxx)
	CodeInvalidBlockNumber:         "block number is invalid",
	CodeBlockEstimationFailed:      "block estimation failed",
	CodeLightAccountOwnerMismatch:  "light-account owner does not match the request",
	CodeVaultERC20AssetNotExists:   "vault ERC-20 asset does not exist",
	CodeVaultERC20PoolNotExists:    "vault ERC-20 pool does not exist",
	CodeVaultAddAssetBeforeBalance: "vault asset must be added before balance updates",
	CodeInvalidSwellSeason:         "Swell season is invalid",
	CodeVaultNotFound:              "vault not found",

	// Maker programs (19xxx)
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

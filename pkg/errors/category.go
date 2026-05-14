// Package errors provides the SDK's error types and sentinel values.
package errors

import (
	"errors"
)

// ErrorCategory buckets every error the SDK might surface into one of a
// small set of operationally-meaningful categories. Use [Category] to
// classify an error; use [IsRetryable] when the only question is
// whether a backoff+retry has a reasonable chance of helping.
//
// The set is intentionally small (ten values) so a switch over it stays
// maintainable. Callers who need finer granularity should drill down
// via [Code] and switch on the underlying [APIError.Code].
type ErrorCategory int

const (
	// CategoryUnknown is the fallback for nil or unrecognised errors.
	CategoryUnknown ErrorCategory = iota
	// CategoryRateLimit covers per-IP request rate limiting (-32000)
	// and the WebSocket concurrency cap (-32100). Retry after a
	// backoff.
	CategoryRateLimit
	// CategoryNetwork covers transport-level failures: dial errors,
	// dropped connections, "not connected" sentinel, subscription
	// closed. Retry after a backoff.
	CategoryNetwork
	// CategoryTimeout covers client-side request timeouts and the
	// engine's order/engine confirmation timeouts (9000, 9001).
	// Generally retryable but order-placement flows must consult the
	// order's nonce or call private/get_order before naive retry —
	// the order may actually have been placed.
	CategoryTimeout
	// CategoryAuth covers signer / session-key / scope failures plus
	// local signing errors. Not retryable; the caller must fix the
	// signer or session-key configuration first.
	CategoryAuth
	// CategoryInvalidRequest covers malformed JSON-RPC envelopes
	// (-32600..-32602, -32700), invalid subscription channels (13000),
	// and SDK-side sentinels for missing configuration. Not retryable;
	// the caller must fix the request.
	CategoryInvalidRequest
	// CategoryNotFound covers "no such record" errors across
	// orders, instruments, assets, accounts, subaccounts, session
	// keys, vaults, and maker programs. Not retryable.
	CategoryNotFound
	// CategoryEngineReject covers logical rejections by the matching
	// engine — insufficient funds, post-only crossed, MMP frozen,
	// instrument inactive, etc. Not retryable in any general sense;
	// the caller has to change inputs.
	CategoryEngineReject
	// CategoryCompliance covers terminal compliance blocks — OFAC,
	// restricted region, account disabled. Not retryable; requires
	// out-of-band escalation.
	CategoryCompliance
	// CategoryInternal covers engine internal errors (-32603, 8xxx,
	// 9002). Generally retryable after a backoff.
	CategoryInternal
)

// String returns the lowercase snake_case name. Useful as a label for
// metrics or structured logs.
func (c ErrorCategory) String() string {
	switch c {
	case CategoryRateLimit:
		return "rate_limit"
	case CategoryNetwork:
		return "network"
	case CategoryTimeout:
		return "timeout"
	case CategoryAuth:
		return "auth"
	case CategoryInvalidRequest:
		return "invalid_request"
	case CategoryNotFound:
		return "not_found"
	case CategoryEngineReject:
		return "engine_reject"
	case CategoryCompliance:
		return "compliance"
	case CategoryInternal:
		return "internal"
	default:
		return "unknown"
	}
}

// Category returns the best-fit [ErrorCategory] for an error. A nil
// error returns [CategoryUnknown].
//
// Resolution order:
//  1. Concrete SDK error types reachable via [errors.As]
//     ([*ConnectionError], [*TimeoutError], [*SigningError],
//     [*ExpiredSignatureError]).
//  2. [*APIError] reachable via [errors.As] — categorised by
//     [APIError.Code].
//  3. Package sentinels matched via [errors.Is].
//  4. [CategoryUnknown] otherwise.
func Category(err error) ErrorCategory {
	if err == nil {
		return CategoryUnknown
	}

	// Concrete SDK types first — they bypass the more expensive Is
	// chain and capture transport-level failures regardless of how
	// the engine code (if any) is categorised.
	var connErr *ConnectionError
	if errors.As(err, &connErr) {
		return CategoryNetwork
	}
	var toErr *TimeoutError
	if errors.As(err, &toErr) {
		return CategoryTimeout
	}
	var signErr *SigningError
	if errors.As(err, &signErr) {
		return CategoryAuth
	}
	var expErr *ExpiredSignatureError
	if errors.As(err, &expErr) {
		return CategoryAuth
	}

	// APIError carries the engine code.
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return categoryForCode(apiErr.Code)
	}

	// SDK sentinels (matched via errors.Is so wrapping works).
	switch {
	case errors.Is(err, ErrNotConnected),
		errors.Is(err, ErrAlreadyConnected),
		errors.Is(err, ErrSubscriptionClosed):
		return CategoryNetwork
	case errors.Is(err, ErrUnauthorized),
		errors.Is(err, ErrInvalidSignature),
		errors.Is(err, ErrSessionKeyExpired),
		errors.Is(err, ErrSessionKeyNotFound),
		errors.Is(err, ErrChainIDMismatch):
		return CategoryAuth
	case errors.Is(err, ErrRateLimited):
		return CategoryRateLimit
	case errors.Is(err, ErrInsufficientFunds),
		errors.Is(err, ErrAlreadyCancelled),
		errors.Is(err, ErrAlreadyFilled),
		errors.Is(err, ErrAlreadyExpired),
		errors.Is(err, ErrMMPFrozen):
		return CategoryEngineReject
	case errors.Is(err, ErrOrderNotFound),
		errors.Is(err, ErrInstrumentNotFound),
		errors.Is(err, ErrSubaccountNotFound),
		errors.Is(err, ErrAccountNotFound):
		return CategoryNotFound
	case errors.Is(err, ErrSubaccountRequired),
		errors.Is(err, ErrInvalidConfig):
		return CategoryInvalidRequest
	case errors.Is(err, ErrRestrictedRegion):
		return CategoryCompliance
	}

	return CategoryUnknown
}

// categoryForCode maps a JSON-RPC engine code to its category. The
// switch is wide on purpose — it's the single source of truth for the
// code → category mapping. Add new codes here as Derive publishes them.
//
//nolint:gocyclo // exhaustive switch is the right shape for an error-code map
func categoryForCode(code int) ErrorCategory {
	switch code {
	// Standard JSON-RPC 2.0
	case CodeRateLimitExceeded, CodeConcurrentWSClientsLimitExceeded:
		return CategoryRateLimit
	case CodeParseError, CodeInvalidRequest, CodeMethodNotFound, CodeInvalidParams:
		return CategoryInvalidRequest
	case CodeInternalError:
		return CategoryInternal

	// Engine / order-confirmation timeouts (9xxx)
	case CodeOrderConfirmationTimeout, CodeEngineConfirmationTimeout:
		return CategoryTimeout
	case CodeCacheConnectionError:
		return CategoryInternal

	// Subscription request errors (13xxx)
	case CodeInvalidChannels:
		return CategoryInvalidRequest

	// Assets / instruments not found (12xxx)
	case CodeAssetNotFound, CodeInstrumentNotFound, CodeCurrencyNotFound:
		return CategoryNotFound

	// Account-class not-founds (14xxx subset)
	case CodeAccountNotFound, CodeSubaccountNotFound, CodeSubaccountWithdrawn,
		CodeSessionKeyNotFound, CodeReferralCodeNotFound:
		return CategoryNotFound

	// Auth (14xxx subset)
	case CodeInvalidSignature,
		CodeHeaderWalletMismatch,
		CodeHeaderWalletMissing,
		CodePrivateChannelSubscriptionFailed,
		CodeSignerNotOwner,
		CodeMissingPrivateParam,
		CodeSessionKeyExpired,
		CodeSessionKeyIPNotWhitelisted,
		CodeUnauthorizedKeyScope,
		CodeScopeNotAdmin,
		CodeUnauthorizedRFQMaker,
		CodeAccountNotWhitelistedAtomicOrders:
		return CategoryAuth

	// Order-lifecycle not-found
	case CodeOrderNotExist:
		return CategoryNotFound

	// Compliance (16xxx)
	case CodeRestrictedRegion, CodeAccountDisabledCompliance, CodeOFACBlocked:
		return CategoryCompliance

	// Vault not-founds (18xxx subset)
	case CodeVaultNotFound, CodeVaultERC20AssetNotExists, CodeVaultERC20PoolNotExists:
		return CategoryNotFound

	// Maker programs (19xxx)
	case CodeMakerProgramNotFound:
		return CategoryNotFound
	}

	// Range-based fallthroughs for less-distinguished groups.
	switch {
	case code >= 8000 && code < 9000:
		// Internal / database errors.
		return CategoryInternal
	case code >= 10000 && code < 11000:
		// Account/wallet logical rejects (allowance, balance, caps).
		return CategoryEngineReject
	case code >= 11000 && code < 11200:
		// Order lifecycle + trigger orders + RFQ + quote.
		return CategoryEngineReject
	case code >= 11200 && code < 11300:
		// Liquidation auction.
		return CategoryEngineReject
	case code >= 12000 && code < 13000:
		// Asset/instrument other (e.g. 12003 USDC caps).
		return CategoryEngineReject
	case code >= 14000 && code < 15000:
		// Anything in the auth/account range we didn't catch above.
		return CategoryEngineReject
	case code >= 16000 && code < 17000:
		return CategoryCompliance
	case code >= 18000 && code < 19000:
		// Vault / smart-account logical rejects.
		return CategoryEngineReject
	}

	return CategoryUnknown
}

// IsRetryable reports whether the SDK or the engine deems an error
// transient — i.e. retrying the same request after a backoff has a
// reasonable chance of success. True for [CategoryRateLimit],
// [CategoryNetwork], [CategoryTimeout], and [CategoryInternal];
// false for everything else.
//
// Caveat: engine confirmation timeouts (codes
// [CodeOrderConfirmationTimeout] and [CodeEngineConfirmationTimeout])
// can mean "the order may have actually placed". Order-submission
// flows should consult the order nonce or call private/get_order
// before naively retrying — IsRetryable returning true is necessary
// but not sufficient for safety.
func IsRetryable(err error) bool {
	switch Category(err) {
	case CategoryRateLimit,
		CategoryNetwork,
		CategoryTimeout,
		CategoryInternal:
		return true
	default:
		return false
	}
}

// Code extracts the JSON-RPC engine code from any error in the SDK's
// error chain.
//
// Composable replacement for the explicit
//
//	var apiErr *APIError
//	if errors.As(err, &apiErr) { code := apiErr.Code; ... }
//
// idiom. Returns (0, false) for nil, network errors, signing errors,
// and SDK sentinels that don't carry a wire code.
func Code(err error) (int, bool) {
	if err == nil {
		return 0, false
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code, true
	}
	return 0, false
}

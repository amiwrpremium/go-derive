package errors_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
)

func TestErrorCategory_String(t *testing.T) {
	cases := []struct {
		c    derrors.ErrorCategory
		want string
	}{
		{derrors.CategoryUnknown, "unknown"},
		{derrors.CategoryRateLimit, "rate_limit"},
		{derrors.CategoryNetwork, "network"},
		{derrors.CategoryTimeout, "timeout"},
		{derrors.CategoryAuth, "auth"},
		{derrors.CategoryInvalidRequest, "invalid_request"},
		{derrors.CategoryNotFound, "not_found"},
		{derrors.CategoryEngineReject, "engine_reject"},
		{derrors.CategoryCompliance, "compliance"},
		{derrors.CategoryInternal, "internal"},
		{derrors.ErrorCategory(999), "unknown"},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, c.c.String())
	}
}

func TestCategory_NilReturnsUnknown(t *testing.T) {
	assert.Equal(t, derrors.CategoryUnknown, derrors.Category(nil))
}

func TestCategory_UnrecognisedReturnsUnknown(t *testing.T) {
	assert.Equal(t, derrors.CategoryUnknown, derrors.Category(errors.New("just some error")))
}

func TestCategory_ConcreteSDKErrorTypes(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want derrors.ErrorCategory
	}{
		{"ConnectionError", &derrors.ConnectionError{Op: "dial", Err: errors.New("EOF")}, derrors.CategoryNetwork},
		{"TimeoutError", &derrors.TimeoutError{Method: "private/order"}, derrors.CategoryTimeout},
		{"SigningError", &derrors.SigningError{Op: "trade module", Err: errors.New("bad key")}, derrors.CategoryAuth},
		{"ExpiredSignatureError", &derrors.ExpiredSignatureError{ExpiryUnixSec: 100, NowUnixSec: 200}, derrors.CategoryAuth},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, derrors.Category(c.err))
		})
	}
}

func TestCategory_APIErrorByCode(t *testing.T) {
	// One representative code per branch of the categorisation switch.
	cases := []struct {
		code int
		want derrors.ErrorCategory
		name string
	}{
		// Standard JSON-RPC
		{derrors.CodeRateLimitExceeded, derrors.CategoryRateLimit, "rate_limit_exceeded"},
		{derrors.CodeConcurrentWSClientsLimitExceeded, derrors.CategoryRateLimit, "concurrent_ws_clients"},
		{derrors.CodeParseError, derrors.CategoryInvalidRequest, "parse_error"},
		{derrors.CodeInvalidRequest, derrors.CategoryInvalidRequest, "invalid_request"},
		{derrors.CodeMethodNotFound, derrors.CategoryInvalidRequest, "method_not_found"},
		{derrors.CodeInvalidParams, derrors.CategoryInvalidRequest, "invalid_params"},
		{derrors.CodeInternalError, derrors.CategoryInternal, "internal_error"},

		// 8xxx — internal
		{derrors.CodeDatabaseError, derrors.CategoryInternal, "database_error"},
		{derrors.CodeFeedsNotFound, derrors.CategoryInternal, "feeds_not_found_internal"},
		{derrors.CodeCounterpartyInsufficientFunds, derrors.CategoryInternal, "counterparty_insufficient_funds"},

		// 9xxx — timeouts + cache
		{derrors.CodeOrderConfirmationTimeout, derrors.CategoryTimeout, "order_confirmation_timeout"},
		{derrors.CodeEngineConfirmationTimeout, derrors.CategoryTimeout, "engine_confirmation_timeout"},
		{derrors.CodeCacheConnectionError, derrors.CategoryInternal, "cache_connection_error"},

		// 10xxx — engine reject
		{derrors.CodeERC20InsufficientBalance, derrors.CategoryEngineReject, "erc20_insufficient_balance"},
		{derrors.CodeManagerNotFound, derrors.CategoryEngineReject, "manager_not_found"},

		// 11xxx — order lifecycle + RFQ + auction
		{derrors.CodeInsufficientFundsOrder, derrors.CategoryEngineReject, "insufficient_funds_order"},
		{derrors.CodeOrderNotExist, derrors.CategoryNotFound, "order_not_exist"},
		{derrors.CodePostOnlyReject, derrors.CategoryEngineReject, "post_only_reject"},
		{derrors.CodeMMPFrozen, derrors.CategoryEngineReject, "mmp_frozen"},
		{derrors.CodeQuoteTakerCostTooHigh, derrors.CategoryEngineReject, "rfq_quote_taker_cost_too_high"},
		{derrors.CodeAuctionNotOngoing, derrors.CategoryEngineReject, "auction_not_ongoing"},

		// 12xxx — asset/instrument
		{derrors.CodeAssetNotFound, derrors.CategoryNotFound, "asset_not_found"},
		{derrors.CodeInstrumentNotFound, derrors.CategoryNotFound, "instrument_not_found"},
		{derrors.CodeCurrencyNotFound, derrors.CategoryNotFound, "currency_not_found"},
		{derrors.CodeUSDCNoCaps, derrors.CategoryEngineReject, "usdc_no_caps"},

		// 13xxx
		{derrors.CodeInvalidChannels, derrors.CategoryInvalidRequest, "invalid_channels"},

		// 14xxx — auth subset
		{derrors.CodeInvalidSignature, derrors.CategoryAuth, "invalid_signature"},
		{derrors.CodeHeaderWalletMissing, derrors.CategoryAuth, "header_wallet_missing"},
		{derrors.CodeSessionKeyExpired, derrors.CategoryAuth, "session_key_expired"},
		{derrors.CodeUnauthorizedKeyScope, derrors.CategoryAuth, "unauthorized_key_scope"},
		// 14xxx — not-found subset
		{derrors.CodeAccountNotFound, derrors.CategoryNotFound, "account_not_found"},
		{derrors.CodeSubaccountNotFound, derrors.CategoryNotFound, "subaccount_not_found"},
		{derrors.CodeSessionKeyNotFound, derrors.CategoryNotFound, "session_key_not_found"},
		// 14xxx — engine reject fallback
		{derrors.CodeChainIDMismatch, derrors.CategoryEngineReject, "chain_id_mismatch"},

		// 16xxx — compliance
		{derrors.CodeRestrictedRegion, derrors.CategoryCompliance, "restricted_region"},
		{derrors.CodeAccountDisabledCompliance, derrors.CategoryCompliance, "account_disabled"},
		{derrors.CodeOFACBlocked, derrors.CategoryCompliance, "ofac_blocked"},
		{derrors.CodeSentinelAuthInvalid, derrors.CategoryCompliance, "sentinel_auth_invalid"},

		// 18xxx — vault
		{derrors.CodeVaultNotFound, derrors.CategoryNotFound, "vault_not_found"},
		{derrors.CodeInvalidBlockNumber, derrors.CategoryEngineReject, "invalid_block_number"},

		// 19xxx — maker programs
		{derrors.CodeMakerProgramNotFound, derrors.CategoryNotFound, "maker_program_not_found"},

		// unknown code → unknown category
		{99999, derrors.CategoryUnknown, "unknown_code"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := &derrors.APIError{Code: c.code}
			assert.Equal(t, c.want, derrors.Category(err),
				"code %d", c.code)
		})
	}
}

func TestCategory_UnwrapsThroughFmtErrorf(t *testing.T) {
	apiErr := &derrors.APIError{Code: derrors.CodeRateLimitExceeded}
	wrapped := fmt.Errorf("rest: call failed: %w", apiErr)
	assert.Equal(t, derrors.CategoryRateLimit, derrors.Category(wrapped))
}

func TestCategory_Sentinels(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want derrors.ErrorCategory
	}{
		{"ErrNotConnected", derrors.ErrNotConnected, derrors.CategoryNetwork},
		{"ErrAlreadyConnected", derrors.ErrAlreadyConnected, derrors.CategoryNetwork},
		{"ErrSubscriptionClosed", derrors.ErrSubscriptionClosed, derrors.CategoryNetwork},
		{"ErrUnauthorized", derrors.ErrUnauthorized, derrors.CategoryAuth},
		{"ErrInvalidSignature", derrors.ErrInvalidSignature, derrors.CategoryAuth},
		{"ErrSessionKeyExpired", derrors.ErrSessionKeyExpired, derrors.CategoryAuth},
		{"ErrSessionKeyNotFound", derrors.ErrSessionKeyNotFound, derrors.CategoryAuth},
		{"ErrChainIDMismatch", derrors.ErrChainIDMismatch, derrors.CategoryAuth},
		{"ErrRateLimited", derrors.ErrRateLimited, derrors.CategoryRateLimit},
		{"ErrInsufficientFunds", derrors.ErrInsufficientFunds, derrors.CategoryEngineReject},
		{"ErrAlreadyCancelled", derrors.ErrAlreadyCancelled, derrors.CategoryEngineReject},
		{"ErrAlreadyFilled", derrors.ErrAlreadyFilled, derrors.CategoryEngineReject},
		{"ErrAlreadyExpired", derrors.ErrAlreadyExpired, derrors.CategoryEngineReject},
		{"ErrMMPFrozen", derrors.ErrMMPFrozen, derrors.CategoryEngineReject},
		{"ErrOrderNotFound", derrors.ErrOrderNotFound, derrors.CategoryNotFound},
		{"ErrInstrumentNotFound", derrors.ErrInstrumentNotFound, derrors.CategoryNotFound},
		{"ErrSubaccountNotFound", derrors.ErrSubaccountNotFound, derrors.CategoryNotFound},
		{"ErrAccountNotFound", derrors.ErrAccountNotFound, derrors.CategoryNotFound},
		{"ErrSubaccountRequired", derrors.ErrSubaccountRequired, derrors.CategoryInvalidRequest},
		{"ErrInvalidConfig", derrors.ErrInvalidConfig, derrors.CategoryInvalidRequest},
		{"ErrRestrictedRegion", derrors.ErrRestrictedRegion, derrors.CategoryCompliance},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, derrors.Category(c.err))
		})
	}
}

func TestIsRetryable(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"unknown", errors.New("foo"), false},
		{"rate_limit", derrors.ErrRateLimited, true},
		{"network", derrors.ErrNotConnected, true},
		{"timeout_sentinel_via_type", &derrors.TimeoutError{Method: "x"}, true},
		{"engine_timeout_via_code", &derrors.APIError{Code: derrors.CodeOrderConfirmationTimeout}, true},
		{"internal_via_code", &derrors.APIError{Code: derrors.CodeInternalError}, true},
		{"cache_error", &derrors.APIError{Code: derrors.CodeCacheConnectionError}, true},
		{"auth", derrors.ErrUnauthorized, false},
		{"engine_reject", derrors.ErrInsufficientFunds, false},
		{"not_found", derrors.ErrOrderNotFound, false},
		{"compliance", derrors.ErrRestrictedRegion, false},
		{"invalid_request", &derrors.APIError{Code: derrors.CodeInvalidParams}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, derrors.IsRetryable(c.err))
		})
	}
}

func TestCode(t *testing.T) {
	// Plain APIError.
	apiErr := &derrors.APIError{Code: derrors.CodeOrderNotExist, Message: "order does not exist"}
	c, ok := derrors.Code(apiErr)
	assert.True(t, ok)
	assert.Equal(t, derrors.CodeOrderNotExist, c)

	// Wrapped APIError — errors.As must unwrap.
	wrapped := fmt.Errorf("transport: %w", apiErr)
	c, ok = derrors.Code(wrapped)
	assert.True(t, ok)
	assert.Equal(t, derrors.CodeOrderNotExist, c)

	// Non-APIError types return (0, false).
	c, ok = derrors.Code(&derrors.ConnectionError{Op: "dial", Err: errors.New("EOF")})
	assert.False(t, ok)
	assert.Equal(t, 0, c)

	c, ok = derrors.Code(errors.New("nothing special"))
	assert.False(t, ok)
	assert.Equal(t, 0, c)

	c, ok = derrors.Code(nil)
	assert.False(t, ok)
	assert.Equal(t, 0, c)
}

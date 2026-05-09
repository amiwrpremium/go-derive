package derive_test

import (
	"encoding/json"
	stderrors "errors"
	"github.com/amiwrpremium/go-derive"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// FuzzAPIError_UnmarshalJSON guards the on-the-wire error decoder against
// panics on adversarial input.
func FuzzAPIError_UnmarshalJSON(f *testing.F) {
	f.Add([]byte(`{"code":-32000,"message":"throttled"}`))
	f.Add([]byte(`{"code":11015,"message":"frozen","data":{"why":"none"}}`))
	f.Add([]byte(`{"code":"not-a-number"}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`null`))
	f.Add([]byte(``))

	f.Fuzz(func(t *testing.T, raw []byte) {
		var e derive.APIError
		_ = json.Unmarshal(raw, &e)
		_ = e.Error()
		_ = e.CanonicalMessage()
	})
}
func TestAPIError_Error_WithoutData(t *testing.T) {
	e := &derive.APIError{Code: 11000, Message: "insufficient"}
	got := e.Error()
	assert.Contains(t, got, "11000")
	assert.Contains(t, got, "insufficient")
	assert.NotContains(t, got, "(")
}

func TestAPIError_Error_WithData(t *testing.T) {
	e := &derive.APIError{
		Code:    11000,
		Message: "insufficient margin",
		Data:    []byte(`{"need":"500"}`),
	}
	got := e.Error()
	assert.Contains(t, got, "11000")
	assert.Contains(t, got, `{"need":"500"}`)
}

func TestAPIError_Is_RateLimited_StandardCode(t *testing.T) {
	e := &derive.APIError{Code: derive.CodeRateLimitExceeded}
	assert.True(t, derive.Is(e, derive.ErrRateLimited))
}

func TestAPIError_Is_RateLimited_WSConcurrencyCode(t *testing.T) {
	e := &derive.APIError{Code: derive.CodeConcurrentWSClientsLimitExceeded}
	assert.True(t, derive.Is(e, derive.ErrRateLimited))
}

func TestAPIError_Is_RateLimited_RejectsOtherCodes(t *testing.T) {
	e := &derive.APIError{Code: derive.CodeOrderNotExist}
	assert.False(t, derive.Is(e, derive.ErrRateLimited))
}

func TestAPIError_Is_Unauthorized_AllAuthCodes(t *testing.T) {
	cases := []int{
		derive.CodeInvalidSignature,
		derive.CodeHeaderWalletMissing,
		derive.CodeHeaderWalletMismatch,
		derive.CodeSignerNotOwner,
		derive.CodeMissingPrivateParam,
		derive.CodeSessionKeyNotFound,
		derive.CodeSessionKeyExpired,
		derive.CodeSessionKeyIPNotWhitelisted,
		derive.CodeUnauthorizedKeyScope,
		derive.CodeScopeNotAdmin,
		derive.CodeUnauthorizedRFQMaker,
		derive.CodeAccountNotWhitelistedAtomicOrders,
		derive.CodeSentinelAuthInvalid,
	}
	for _, code := range cases {
		t.Run("code", func(t *testing.T) {
			e := &derive.APIError{Code: code}
			assert.True(t, derive.Is(e, derive.ErrUnauthorized), "code %d", code)
		})
	}
}

func TestAPIError_Is_InvalidSignature(t *testing.T) {
	e := &derive.APIError{Code: derive.CodeInvalidSignature}
	assert.True(t, derive.Is(e, derive.ErrInvalidSignature))
	assert.True(t, derive.Is(e, derive.ErrUnauthorized), "should also match the broader auth bucket")
}

func TestAPIError_Is_SessionKeyExpired(t *testing.T) {
	e := &derive.APIError{Code: derive.CodeSessionKeyExpired}
	assert.True(t, derive.Is(e, derive.ErrSessionKeyExpired))
	assert.True(t, derive.Is(e, derive.ErrUnauthorized))
}

func TestAPIError_Is_SessionKeyNotFound(t *testing.T) {
	e := &derive.APIError{Code: derive.CodeSessionKeyNotFound}
	assert.True(t, derive.Is(e, derive.ErrSessionKeyNotFound))
}

func TestAPIError_Is_InsufficientFunds_OrderCode(t *testing.T) {
	assert.True(t, derive.Is(
		&derive.APIError{Code: derive.CodeInsufficientFundsOrder},
		derive.ErrInsufficientFunds))
}

func TestAPIError_Is_InsufficientFunds_ERC20Codes(t *testing.T) {
	for _, code := range []int{derive.CodeERC20InsufficientAllowance, derive.CodeERC20InsufficientBalance} {
		e := &derive.APIError{Code: code}
		assert.True(t, derive.Is(e, derive.ErrInsufficientFunds), "code %d", code)
	}
}

func TestAPIError_Is_OrderNotFound(t *testing.T) {
	e := &derive.APIError{Code: derive.CodeOrderNotExist}
	assert.True(t, derive.Is(e, derive.ErrOrderNotFound))
}

func TestAPIError_Is_AlreadyCancelled(t *testing.T) {
	assert.True(t, derive.Is(&derive.APIError{Code: derive.CodeAlreadyCancelled}, derive.ErrAlreadyCancelled))
}

func TestAPIError_Is_AlreadyFilled(t *testing.T) {
	assert.True(t, derive.Is(&derive.APIError{Code: derive.CodeAlreadyFilled}, derive.ErrAlreadyFilled))
}

func TestAPIError_Is_AlreadyExpired(t *testing.T) {
	assert.True(t, derive.Is(&derive.APIError{Code: derive.CodeAlreadyExpired}, derive.ErrAlreadyExpired))
}

func TestAPIError_Is_InstrumentNotFound_BothCodes(t *testing.T) {
	for _, code := range []int{derive.CodeInstrumentNotFound, derive.CodeAssetNotFound} {
		e := &derive.APIError{Code: code}
		assert.True(t, derive.Is(e, derive.ErrInstrumentNotFound), "code %d", code)
	}
}

func TestAPIError_Is_SubaccountNotFound(t *testing.T) {
	e := &derive.APIError{Code: derive.CodeSubaccountNotFound}
	assert.True(t, derive.Is(e, derive.ErrSubaccountNotFound))
}

func TestAPIError_Is_AccountNotFound(t *testing.T) {
	e := &derive.APIError{Code: derive.CodeAccountNotFound}
	assert.True(t, derive.Is(e, derive.ErrAccountNotFound))
}

func TestAPIError_Is_ChainIDMismatch(t *testing.T) {
	e := &derive.APIError{Code: derive.CodeChainIDMismatch}
	assert.True(t, derive.Is(e, derive.ErrChainIDMismatch))
}

func TestAPIError_Is_MMPFrozen(t *testing.T) {
	e := &derive.APIError{Code: derive.CodeMMPFrozen}
	assert.True(t, derive.Is(e, derive.ErrMMPFrozen))
}

func TestAPIError_Is_RestrictedRegion_AllArms(t *testing.T) {
	for _, code := range []int{derive.CodeRestrictedRegion, derive.CodeAccountDisabledCompliance} {
		e := &derive.APIError{Code: code}
		assert.True(t, derive.Is(e, derive.ErrRestrictedRegion), "code %d", code)
	}
}

func TestAPIError_Is_UnknownCode_FallsThrough(t *testing.T) {
	e := &derive.APIError{Code: 99999}
	assert.False(t, derive.Is(e, derive.ErrRateLimited))
	assert.False(t, derive.Is(e, derive.ErrUnauthorized))
	assert.False(t, derive.Is(e, derive.ErrInsufficientFunds))
}

func TestAPIError_Is_NonSentinelTarget(t *testing.T) {
	e := &derive.APIError{Code: derive.CodeRateLimitExceeded}
	assert.False(t, derive.Is(e, stderrors.New("unrelated")))
}

func TestAPIError_Is_DefaultArm_SwitchFallthrough(t *testing.T) {

	e := &derive.APIError{Code: derive.CodeTriggerOrderCancelled}
	for _, sentinel := range []error{
		derive.ErrRateLimited, derive.ErrUnauthorized,
		derive.ErrInvalidSignature, derive.ErrSessionKeyExpired,
		derive.ErrSessionKeyNotFound, derive.ErrInsufficientFunds,
		derive.ErrOrderNotFound, derive.ErrAlreadyCancelled,
		derive.ErrAlreadyFilled, derive.ErrAlreadyExpired,
		derive.ErrInstrumentNotFound, derive.ErrSubaccountNotFound,
		derive.ErrAccountNotFound, derive.ErrChainIDMismatch,
		derive.ErrMMPFrozen, derive.ErrRestrictedRegion,
	} {
		assert.False(t, derive.Is(e, sentinel),
			"trigger-cancelled (%d) should not match %v", e.Code, sentinel)
	}
}

func TestAPIError_Implements_ErrorInterface(_ *testing.T) {
	var _ error = &derive.APIError{}
}

// allCodes is the canonical map of every Code* constant exported by the
// package. Adding a new code requires adding it here; the duplicate-detection
// test below catches drift.
func allCodes() map[string]int {
	return map[string]int{
		"NoError":                           derive.CodeNoError,
		"ParseError":                        derive.CodeParseError,
		"InvalidRequest":                    derive.CodeInvalidRequest,
		"MethodNotFound":                    derive.CodeMethodNotFound,
		"InvalidParams":                     derive.CodeInvalidParams,
		"InternalError":                     derive.CodeInternalError,
		"RateLimitExceeded":                 derive.CodeRateLimitExceeded,
		"ConcurrentWSClientsLimitExceeded":  derive.CodeConcurrentWSClientsLimitExceeded,
		"OrderConfirmationTimeout":          derive.CodeOrderConfirmationTimeout,
		"EngineConfirmationTimeout":         derive.CodeEngineConfirmationTimeout,
		"ManagerNotFound":                   derive.CodeManagerNotFound,
		"AssetNotERC20":                     derive.CodeAssetNotERC20,
		"WalletMismatch":                    derive.CodeWalletMismatch,
		"SubaccountMismatch":                derive.CodeSubaccountMismatch,
		"MultipleCurrenciesNotSupported":    derive.CodeMultipleCurrenciesNotSupported,
		"MaxSubaccountsReached":             derive.CodeMaxSubaccountsReached,
		"MaxSessionKeysReached":             derive.CodeMaxSessionKeysReached,
		"MaxAssetsPerSubaccount":            derive.CodeMaxAssetsPerSubaccount,
		"MaxExpiriesPerSubaccount":          derive.CodeMaxExpiriesPerSubaccount,
		"InvalidRecipientSubaccountID":      derive.CodeInvalidRecipientSubaccountID,
		"PMRMUSDCOnlyCollateral":            derive.CodePMRMUSDCOnlyCollateral,
		"ERC20InsufficientAllowance":        derive.CodeERC20InsufficientAllowance,
		"ERC20InsufficientBalance":          derive.CodeERC20InsufficientBalance,
		"PendingDeposit":                    derive.CodePendingDeposit,
		"PendingWithdrawal":                 derive.CodePendingWithdrawal,
		"PM2CollateralConstraint":           derive.CodePM2CollateralConstraint,
		"InsufficientFundsOrder":            derive.CodeInsufficientFundsOrder,
		"OrderRejectedFromQueue":            derive.CodeOrderRejectedFromQueue,
		"AlreadyCancelled":                  derive.CodeAlreadyCancelled,
		"AlreadyFilled":                     derive.CodeAlreadyFilled,
		"AlreadyExpired":                    derive.CodeAlreadyExpired,
		"OrderNotExist":                     derive.CodeOrderNotExist,
		"SelfCrossDisallowed":               derive.CodeSelfCrossDisallowed,
		"PostOnlyReject":                    derive.CodePostOnlyReject,
		"ZeroLiquidity":                     derive.CodeZeroLiquidity,
		"PostOnlyInvalidType":               derive.CodePostOnlyInvalidType,
		"InvalidSignatureExpiry":            derive.CodeInvalidSignatureExpiry,
		"InvalidAmount":                     derive.CodeInvalidAmount,
		"InvalidLimitPrice":                 derive.CodeInvalidLimitPrice,
		"FOKNotFilled":                      derive.CodeFOKNotFilled,
		"MMPFrozen":                         derive.CodeMMPFrozen,
		"AlreadyConsumed":                   derive.CodeAlreadyConsumed,
		"NonUniqueNonce":                    derive.CodeNonUniqueNonce,
		"InvalidNonceDate":                  derive.CodeInvalidNonceDate,
		"OpenOrdersLimitExceeded":           derive.CodeOpenOrdersLimitExceeded,
		"NegativeERC20Balance":              derive.CodeNegativeERC20Balance,
		"InstrumentNotLive":                 derive.CodeInstrumentNotLive,
		"TriggerOrderCancelled":             derive.CodeTriggerOrderCancelled,
		"InvalidTriggerPrice":               derive.CodeInvalidTriggerPrice,
		"TriggerOrderLimitExceeded":         derive.CodeTriggerOrderLimitExceeded,
		"TriggerPriceTypeUnsupported":       derive.CodeTriggerPriceTypeUnsupported,
		"TriggerOrderReplaceUnsupported":    derive.CodeTriggerOrderReplaceUnsupported,
		"MarketOrderInvalidTriggerPrice":    derive.CodeMarketOrderInvalidTriggerPrice,
		"LegInstrumentsNotUnique":           derive.CodeLegInstrumentsNotUnique,
		"RFQNotFound":                       derive.CodeRFQNotFound,
		"QuoteNotFound":                     derive.CodeQuoteNotFound,
		"RFQLegMismatch":                    derive.CodeRFQLegMismatch,
		"RFQNotOpen":                        derive.CodeRFQNotOpen,
		"RFQIDMismatch":                     derive.CodeRFQIDMismatch,
		"InvalidRFQCounterparty":            derive.CodeInvalidRFQCounterparty,
		"QuoteCostTooHigh":                  derive.CodeQuoteCostTooHigh,
		"AuctionNotOngoing":                 derive.CodeAuctionNotOngoing,
		"OpenOrdersNotAllowed":              derive.CodeOpenOrdersNotAllowed,
		"PriceLimitExceeded":                derive.CodePriceLimitExceeded,
		"LastTradeIDMismatch":               derive.CodeLastTradeIDMismatch,
		"AssetNotFound":                     derive.CodeAssetNotFound,
		"InstrumentNotFound":                derive.CodeInstrumentNotFound,
		"CurrencyNotFound":                  derive.CodeCurrencyNotFound,
		"USDCNoCaps":                        derive.CodeUSDCNoCaps,
		"InvalidChannels":                   derive.CodeInvalidChannels,
		"AccountNotFound":                   derive.CodeAccountNotFound,
		"SubaccountNotFound":                derive.CodeSubaccountNotFound,
		"SubaccountWithdrawn":               derive.CodeSubaccountWithdrawn,
		"SessionKeyExpiryTooLow":            derive.CodeSessionKeyExpiryTooLow,
		"SessionKeyAlreadyRegistered":       derive.CodeSessionKeyAlreadyRegistered,
		"SessionKeyRegisteredOtherAccount":  derive.CodeSessionKeyRegisteredOtherAccount,
		"AddressNotChecksummed":             derive.CodeAddressNotChecksummed,
		"InvalidEthAddress":                 derive.CodeInvalidEthAddress,
		"InvalidSignature":                  derive.CodeInvalidSignature,
		"NonceMismatch":                     derive.CodeNonceMismatch,
		"RawTxFunctionMismatch":             derive.CodeRawTxFunctionMismatch,
		"RawTxContractMismatch":             derive.CodeRawTxContractMismatch,
		"RawTxParamsMismatch":               derive.CodeRawTxParamsMismatch,
		"RawTxParamValuesMismatch":          derive.CodeRawTxParamValuesMismatch,
		"HeaderWalletMismatch":              derive.CodeHeaderWalletMismatch,
		"HeaderWalletMissing":               derive.CodeHeaderWalletMissing,
		"PrivateChannelSubscriptionFailed":  derive.CodePrivateChannelSubscriptionFailed,
		"SignerNotOwner":                    derive.CodeSignerNotOwner,
		"ChainIDMismatch":                   derive.CodeChainIDMismatch,
		"MissingPrivateParam":               derive.CodeMissingPrivateParam,
		"SessionKeyNotFound":                derive.CodeSessionKeyNotFound,
		"UnauthorizedRFQMaker":              derive.CodeUnauthorizedRFQMaker,
		"CrossCurrencyRFQNotSupported":      derive.CodeCrossCurrencyRFQNotSupported,
		"SessionKeyIPNotWhitelisted":        derive.CodeSessionKeyIPNotWhitelisted,
		"SessionKeyExpired":                 derive.CodeSessionKeyExpired,
		"UnauthorizedKeyScope":              derive.CodeUnauthorizedKeyScope,
		"ScopeNotAdmin":                     derive.CodeScopeNotAdmin,
		"AccountNotWhitelistedAtomicOrders": derive.CodeAccountNotWhitelistedAtomicOrders,
		"ReferralCodeNotFound":              derive.CodeReferralCodeNotFound,
		"RestrictedRegion":                  derive.CodeRestrictedRegion,
		"AccountDisabledCompliance":         derive.CodeAccountDisabledCompliance,
		"SentinelAuthInvalid":               derive.CodeSentinelAuthInvalid,
		"InvalidBlockNumber":                derive.CodeInvalidBlockNumber,
		"BlockEstimationFailed":             derive.CodeBlockEstimationFailed,
		"LightAccountOwnerMismatch":         derive.CodeLightAccountOwnerMismatch,
		"VaultERC20AssetNotExists":          derive.CodeVaultERC20AssetNotExists,
		"VaultERC20PoolNotExists":           derive.CodeVaultERC20PoolNotExists,
		"VaultAddAssetBeforeBalance":        derive.CodeVaultAddAssetBeforeBalance,
		"InvalidSwellSeason":                derive.CodeInvalidSwellSeason,
		"VaultNotFound":                     derive.CodeVaultNotFound,
		"MakerProgramNotFound":              derive.CodeMakerProgramNotFound,
	}
}

func TestCodes_AllUnique(t *testing.T) {
	all := allCodes()
	seen := map[int]string{}
	for name, code := range all {
		if other, dup := seen[code]; dup {
			t.Errorf("duplicate code %d: %s and %s", code, name, other)
		}
		seen[code] = name
	}
}

func TestCodes_StandardJSONRPC_AllNegative(t *testing.T) {
	for _, c := range []int{
		derive.CodeParseError,
		derive.CodeInvalidRequest,
		derive.CodeMethodNotFound,
		derive.CodeInvalidParams,
		derive.CodeInternalError,
		derive.CodeRateLimitExceeded,
		derive.CodeConcurrentWSClientsLimitExceeded,
	} {
		assert.Less(t, c, 0)
	}
}

func TestCodes_DeriveSpecific_NonNegative(t *testing.T) {
	exempt := map[int]bool{
		derive.CodeParseError:                       true,
		derive.CodeInvalidRequest:                   true,
		derive.CodeMethodNotFound:                   true,
		derive.CodeInvalidParams:                    true,
		derive.CodeInternalError:                    true,
		derive.CodeRateLimitExceeded:                true,
		derive.CodeConcurrentWSClientsLimitExceeded: true,
	}
	for name, c := range allCodes() {
		if exempt[c] {
			continue
		}
		assert.GreaterOrEqual(t, c, 0, "code %s (%d) should be non-negative", name, c)
	}
}

func TestCodes_KnownValues(t *testing.T) {
	assert.Equal(t, -32000, derive.CodeRateLimitExceeded)
	assert.Equal(t, 11006, derive.CodeOrderNotExist)
	assert.Equal(t, 14014, derive.CodeInvalidSignature)
	assert.Equal(t, 14030, derive.CodeSessionKeyExpired)
	assert.Equal(t, 11015, derive.CodeMMPFrozen)
	assert.Equal(t, 16000, derive.CodeRestrictedRegion)
	assert.Equal(t, 19000, derive.CodeMakerProgramNotFound)
}

func TestCodes_TotalCount(t *testing.T) {

	assert.Equal(t, 111, len(allCodes()))
}
func TestSentinels_DistinctValues(t *testing.T) {
	assert.NotEqual(t, derive.ErrUnauthorized, derive.ErrRateLimited)
	assert.NotEqual(t, derive.ErrNotConnected, derive.ErrAlreadyConnected)
	assert.NotEqual(t, derive.ErrSubaccountRequired, derive.ErrInvalidConfig)
	assert.NotEqual(t, derive.ErrSubscriptionClosed, derive.ErrNotConnected)
}

func TestSentinels_HaveMessages(t *testing.T) {
	for _, e := range []error{
		derive.ErrNotConnected,
		derive.ErrAlreadyConnected,
		derive.ErrUnauthorized,
		derive.ErrRateLimited,
		derive.ErrSubscriptionClosed,
		derive.ErrSubaccountRequired,
		derive.ErrInvalidConfig,
	} {
		assert.NotEmpty(t, e.Error())
		assert.Contains(t, e.Error(), "derive")
	}
}

func TestExportedStdlibHelpers(t *testing.T) {
	assert.NotNil(t, derive.Is)
	assert.NotNil(t, derive.As)
	assert.NotNil(t, derive.Unwrap)
	assert.NotNil(t, derive.New)

	e := derive.New("hello")
	assert.Equal(t, "hello", e.Error())
}
func TestDescription_KnownCodeReturnsText(t *testing.T) {
	got := derive.Description(derive.CodeRateLimitExceeded)
	assert.NotEmpty(t, got)
	assert.Contains(t, got, "rate")
}

func TestDescription_UnknownCodeReturnsEmpty(t *testing.T) {
	assert.Equal(t, "", derive.Description(99999))
}

func TestHasDescription_KnownCode(t *testing.T) {
	assert.True(t, derive.HasDescription(derive.CodeInvalidSignature))
}

func TestHasDescription_UnknownCode(t *testing.T) {
	assert.False(t, derive.HasDescription(99999))
}

// TestDescription_AllCodesHaveText runs a description coverage check —
// every Code* constant declared in codes.go must have a non-empty entry
// in the message map. If a future PR adds a new code without a message,
// this test fails immediately.
func TestDescription_AllCodesHaveText(t *testing.T) {
	for name, code := range allCodes() {
		t.Run(name, func(t *testing.T) {
			got := derive.Description(code)
			assert.NotEmpty(t, got, "code %d (%s) has no description", code, name)
		})
	}
}

func TestDescription_NoTrailingPunctuation(t *testing.T) {

	for name, code := range allCodes() {
		desc := derive.Description(code)
		assert.False(t, strings.HasSuffix(desc, "."),
			"description for %s ends with a period: %q", name, desc)
	}
}

func TestDescription_LowercaseStart(t *testing.T) {
	for name, code := range allCodes() {
		desc := derive.Description(code)
		if desc == "" {
			continue
		}
		first := desc[0]

		if first >= 'A' && first <= 'Z' {

			allowed := []string{"USDC", "X-LyraWallet", "RFQ", "PMRM", "WebSocket", "ERC-20", "Swell"}
			ok := false
			for _, p := range allowed {
				if strings.HasPrefix(desc, p) {
					ok = true
					break
				}
			}
			assert.True(t, ok, "%s: description starts with uppercase but isn't a known proper noun: %q", name, desc)
		}
	}
}

func TestAPIError_Error_FillsInCanonicalWhenMessageEmpty(t *testing.T) {
	e := &derive.APIError{Code: derive.CodeRateLimitExceeded}
	got := e.Error()
	assert.Contains(t, got, "rate")
	assert.Contains(t, got, "-32000")
}

func TestAPIError_Error_KeepsServerMessageWhenPresent(t *testing.T) {
	e := &derive.APIError{
		Code:    derive.CodeRateLimitExceeded,
		Message: "Custom server message",
	}
	got := e.Error()
	assert.Contains(t, got, "Custom server message")
	assert.NotContains(t, got, "rate limit window")
}

func TestAPIError_Error_UnknownCodeAndEmptyMessage(t *testing.T) {
	e := &APIError99999{}
	_ = e
	apiErr := &derive.APIError{Code: 99999}
	got := apiErr.Error()

	assert.Contains(t, got, "99999")
}

func TestAPIError_CanonicalMessage_KnownCode(t *testing.T) {
	e := &derive.APIError{Code: derive.CodeMMPFrozen}
	assert.Contains(t, e.CanonicalMessage(), "market-maker protection")
}

func TestAPIError_CanonicalMessage_UnknownCode(t *testing.T) {
	e := &derive.APIError{Code: 99999}
	assert.Equal(t, "", e.CanonicalMessage())
}

// dummy private type to silence unused-import diagnostics if test file
// reorganisation drops the apiError synonym; harmless.
type APIError99999 struct{}

func TestConnectionError_Error(t *testing.T) {
	e := &derive.ConnectionError{Op: "dial", Err: stderrors.New("refused")}
	got := e.Error()
	assert.Contains(t, got, "dial")
	assert.Contains(t, got, "refused")
}

func TestConnectionError_Unwrap(t *testing.T) {
	inner := stderrors.New("dial timeout")
	e := &derive.ConnectionError{Op: "dial", Err: inner}
	assert.True(t, stderrors.Is(e, inner))
	assert.Same(t, inner, stderrors.Unwrap(e))
}

func TestConnectionError_NilInner(t *testing.T) {

	e := &derive.ConnectionError{Op: "dial"}
	assert.Contains(t, e.Error(), "dial")
	assert.Nil(t, stderrors.Unwrap(e))
}

func TestTimeoutError_Error(t *testing.T) {
	e := &derive.TimeoutError{Method: "private/order"}
	got := e.Error()
	assert.Contains(t, got, "private/order")
	assert.Contains(t, got, "timeout")
}

func TestTimeoutError_Implements_Error(_ *testing.T) {
	var _ error = &derive.TimeoutError{}
}
func TestSigningError_Error(t *testing.T) {
	e := &derive.SigningError{Op: "parse", Err: stderrors.New("bad key")}
	got := e.Error()
	assert.Contains(t, got, "parse")
	assert.Contains(t, got, "bad key")
}

func TestSigningError_Unwrap(t *testing.T) {
	inner := stderrors.New("bad key")
	e := &derive.SigningError{Op: "parse", Err: inner}
	assert.True(t, stderrors.Is(e, inner))
	assert.Same(t, inner, stderrors.Unwrap(e))
}

func TestSigningError_NilInner(t *testing.T) {
	e := &derive.SigningError{Op: "noop"}
	assert.Contains(t, e.Error(), "noop")
	assert.Nil(t, stderrors.Unwrap(e))
}

func TestExpiredSignatureError_Error(t *testing.T) {
	e := &derive.ExpiredSignatureError{ExpiryUnixSec: 100, NowUnixSec: 200}
	got := e.Error()
	assert.Contains(t, got, "100")
	assert.Contains(t, got, "200")
	assert.Contains(t, got, "expired")
}

func TestExpiredSignatureError_FutureExpiry(t *testing.T) {

	e := &derive.ExpiredSignatureError{ExpiryUnixSec: 1000, NowUnixSec: 500}
	assert.Contains(t, e.Error(), "1000")
}

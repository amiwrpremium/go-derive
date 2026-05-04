package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
)

// allCodes is the canonical map of every Code* constant exported by the
// package. Adding a new code requires adding it here; the duplicate-detection
// test below catches drift.
func allCodes() map[string]int {
	return map[string]int{
		"NoError":                           derrors.CodeNoError,
		"ParseError":                        derrors.CodeParseError,
		"InvalidRequest":                    derrors.CodeInvalidRequest,
		"MethodNotFound":                    derrors.CodeMethodNotFound,
		"InvalidParams":                     derrors.CodeInvalidParams,
		"InternalError":                     derrors.CodeInternalError,
		"RateLimitExceeded":                 derrors.CodeRateLimitExceeded,
		"ConcurrentWSClientsLimitExceeded":  derrors.CodeConcurrentWSClientsLimitExceeded,
		"OrderConfirmationTimeout":          derrors.CodeOrderConfirmationTimeout,
		"EngineConfirmationTimeout":         derrors.CodeEngineConfirmationTimeout,
		"ManagerNotFound":                   derrors.CodeManagerNotFound,
		"AssetNotERC20":                     derrors.CodeAssetNotERC20,
		"WalletMismatch":                    derrors.CodeWalletMismatch,
		"SubaccountMismatch":                derrors.CodeSubaccountMismatch,
		"MultipleCurrenciesNotSupported":    derrors.CodeMultipleCurrenciesNotSupported,
		"MaxSubaccountsReached":             derrors.CodeMaxSubaccountsReached,
		"MaxSessionKeysReached":             derrors.CodeMaxSessionKeysReached,
		"MaxAssetsPerSubaccount":            derrors.CodeMaxAssetsPerSubaccount,
		"MaxExpiriesPerSubaccount":          derrors.CodeMaxExpiriesPerSubaccount,
		"InvalidRecipientSubaccountID":      derrors.CodeInvalidRecipientSubaccountID,
		"PMRMUSDCOnlyCollateral":            derrors.CodePMRMUSDCOnlyCollateral,
		"ERC20InsufficientAllowance":        derrors.CodeERC20InsufficientAllowance,
		"ERC20InsufficientBalance":          derrors.CodeERC20InsufficientBalance,
		"PendingDeposit":                    derrors.CodePendingDeposit,
		"PendingWithdrawal":                 derrors.CodePendingWithdrawal,
		"PM2CollateralConstraint":           derrors.CodePM2CollateralConstraint,
		"InsufficientFundsOrder":            derrors.CodeInsufficientFundsOrder,
		"OrderRejectedFromQueue":            derrors.CodeOrderRejectedFromQueue,
		"AlreadyCancelled":                  derrors.CodeAlreadyCancelled,
		"AlreadyFilled":                     derrors.CodeAlreadyFilled,
		"AlreadyExpired":                    derrors.CodeAlreadyExpired,
		"OrderNotExist":                     derrors.CodeOrderNotExist,
		"SelfCrossDisallowed":               derrors.CodeSelfCrossDisallowed,
		"PostOnlyReject":                    derrors.CodePostOnlyReject,
		"ZeroLiquidity":                     derrors.CodeZeroLiquidity,
		"PostOnlyInvalidType":               derrors.CodePostOnlyInvalidType,
		"InvalidSignatureExpiry":            derrors.CodeInvalidSignatureExpiry,
		"InvalidAmount":                     derrors.CodeInvalidAmount,
		"InvalidLimitPrice":                 derrors.CodeInvalidLimitPrice,
		"FOKNotFilled":                      derrors.CodeFOKNotFilled,
		"MMPFrozen":                         derrors.CodeMMPFrozen,
		"AlreadyConsumed":                   derrors.CodeAlreadyConsumed,
		"NonUniqueNonce":                    derrors.CodeNonUniqueNonce,
		"InvalidNonceDate":                  derrors.CodeInvalidNonceDate,
		"OpenOrdersLimitExceeded":           derrors.CodeOpenOrdersLimitExceeded,
		"NegativeERC20Balance":              derrors.CodeNegativeERC20Balance,
		"InstrumentNotLive":                 derrors.CodeInstrumentNotLive,
		"TriggerOrderCancelled":             derrors.CodeTriggerOrderCancelled,
		"InvalidTriggerPrice":               derrors.CodeInvalidTriggerPrice,
		"TriggerOrderLimitExceeded":         derrors.CodeTriggerOrderLimitExceeded,
		"TriggerPriceTypeUnsupported":       derrors.CodeTriggerPriceTypeUnsupported,
		"TriggerOrderReplaceUnsupported":    derrors.CodeTriggerOrderReplaceUnsupported,
		"MarketOrderInvalidTriggerPrice":    derrors.CodeMarketOrderInvalidTriggerPrice,
		"LegInstrumentsNotUnique":           derrors.CodeLegInstrumentsNotUnique,
		"RFQNotFound":                       derrors.CodeRFQNotFound,
		"QuoteNotFound":                     derrors.CodeQuoteNotFound,
		"RFQLegMismatch":                    derrors.CodeRFQLegMismatch,
		"RFQNotOpen":                        derrors.CodeRFQNotOpen,
		"RFQIDMismatch":                     derrors.CodeRFQIDMismatch,
		"InvalidRFQCounterparty":            derrors.CodeInvalidRFQCounterparty,
		"QuoteCostTooHigh":                  derrors.CodeQuoteCostTooHigh,
		"AuctionNotOngoing":                 derrors.CodeAuctionNotOngoing,
		"OpenOrdersNotAllowed":              derrors.CodeOpenOrdersNotAllowed,
		"PriceLimitExceeded":                derrors.CodePriceLimitExceeded,
		"LastTradeIDMismatch":               derrors.CodeLastTradeIDMismatch,
		"AssetNotFound":                     derrors.CodeAssetNotFound,
		"InstrumentNotFound":                derrors.CodeInstrumentNotFound,
		"CurrencyNotFound":                  derrors.CodeCurrencyNotFound,
		"USDCNoCaps":                        derrors.CodeUSDCNoCaps,
		"InvalidChannels":                   derrors.CodeInvalidChannels,
		"AccountNotFound":                   derrors.CodeAccountNotFound,
		"SubaccountNotFound":                derrors.CodeSubaccountNotFound,
		"SubaccountWithdrawn":               derrors.CodeSubaccountWithdrawn,
		"SessionKeyExpiryTooLow":            derrors.CodeSessionKeyExpiryTooLow,
		"SessionKeyAlreadyRegistered":       derrors.CodeSessionKeyAlreadyRegistered,
		"SessionKeyRegisteredOtherAccount":  derrors.CodeSessionKeyRegisteredOtherAccount,
		"AddressNotChecksummed":             derrors.CodeAddressNotChecksummed,
		"InvalidEthAddress":                 derrors.CodeInvalidEthAddress,
		"InvalidSignature":                  derrors.CodeInvalidSignature,
		"NonceMismatch":                     derrors.CodeNonceMismatch,
		"RawTxFunctionMismatch":             derrors.CodeRawTxFunctionMismatch,
		"RawTxContractMismatch":             derrors.CodeRawTxContractMismatch,
		"RawTxParamsMismatch":               derrors.CodeRawTxParamsMismatch,
		"RawTxParamValuesMismatch":          derrors.CodeRawTxParamValuesMismatch,
		"HeaderWalletMismatch":              derrors.CodeHeaderWalletMismatch,
		"HeaderWalletMissing":               derrors.CodeHeaderWalletMissing,
		"PrivateChannelSubscriptionFailed":  derrors.CodePrivateChannelSubscriptionFailed,
		"SignerNotOwner":                    derrors.CodeSignerNotOwner,
		"ChainIDMismatch":                   derrors.CodeChainIDMismatch,
		"MissingPrivateParam":               derrors.CodeMissingPrivateParam,
		"SessionKeyNotFound":                derrors.CodeSessionKeyNotFound,
		"UnauthorizedRFQMaker":              derrors.CodeUnauthorizedRFQMaker,
		"CrossCurrencyRFQNotSupported":      derrors.CodeCrossCurrencyRFQNotSupported,
		"SessionKeyIPNotWhitelisted":        derrors.CodeSessionKeyIPNotWhitelisted,
		"SessionKeyExpired":                 derrors.CodeSessionKeyExpired,
		"UnauthorizedKeyScope":              derrors.CodeUnauthorizedKeyScope,
		"ScopeNotAdmin":                     derrors.CodeScopeNotAdmin,
		"AccountNotWhitelistedAtomicOrders": derrors.CodeAccountNotWhitelistedAtomicOrders,
		"ReferralCodeNotFound":              derrors.CodeReferralCodeNotFound,
		"RestrictedRegion":                  derrors.CodeRestrictedRegion,
		"AccountDisabledCompliance":         derrors.CodeAccountDisabledCompliance,
		"SentinelAuthInvalid":               derrors.CodeSentinelAuthInvalid,
		"InvalidBlockNumber":                derrors.CodeInvalidBlockNumber,
		"BlockEstimationFailed":             derrors.CodeBlockEstimationFailed,
		"LightAccountOwnerMismatch":         derrors.CodeLightAccountOwnerMismatch,
		"VaultERC20AssetNotExists":          derrors.CodeVaultERC20AssetNotExists,
		"VaultERC20PoolNotExists":           derrors.CodeVaultERC20PoolNotExists,
		"VaultAddAssetBeforeBalance":        derrors.CodeVaultAddAssetBeforeBalance,
		"InvalidSwellSeason":                derrors.CodeInvalidSwellSeason,
		"VaultNotFound":                     derrors.CodeVaultNotFound,
		"MakerProgramNotFound":              derrors.CodeMakerProgramNotFound,
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
		derrors.CodeParseError,
		derrors.CodeInvalidRequest,
		derrors.CodeMethodNotFound,
		derrors.CodeInvalidParams,
		derrors.CodeInternalError,
		derrors.CodeRateLimitExceeded,
		derrors.CodeConcurrentWSClientsLimitExceeded,
	} {
		assert.Less(t, c, 0)
	}
}

func TestCodes_DeriveSpecific_NonNegative(t *testing.T) {
	exempt := map[int]bool{
		derrors.CodeParseError:                       true,
		derrors.CodeInvalidRequest:                   true,
		derrors.CodeMethodNotFound:                   true,
		derrors.CodeInvalidParams:                    true,
		derrors.CodeInternalError:                    true,
		derrors.CodeRateLimitExceeded:                true,
		derrors.CodeConcurrentWSClientsLimitExceeded: true,
	}
	for name, c := range allCodes() {
		if exempt[c] {
			continue
		}
		assert.GreaterOrEqual(t, c, 0, "code %s (%d) should be non-negative", name, c)
	}
}

func TestCodes_KnownValues(t *testing.T) {
	assert.Equal(t, -32000, derrors.CodeRateLimitExceeded)
	assert.Equal(t, 11006, derrors.CodeOrderNotExist)
	assert.Equal(t, 14014, derrors.CodeInvalidSignature)
	assert.Equal(t, 14030, derrors.CodeSessionKeyExpired)
	assert.Equal(t, 11015, derrors.CodeMMPFrozen)
	assert.Equal(t, 16000, derrors.CodeRestrictedRegion)
	assert.Equal(t, 19000, derrors.CodeMakerProgramNotFound)
}

func TestCodes_TotalCount(t *testing.T) {
	// Sanity check that all documented codes are present.
	assert.Equal(t, 111, len(allCodes()))
}

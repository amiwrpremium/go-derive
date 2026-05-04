package errors_test

import (
	stderrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"

	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
)

func TestAPIError_Error_WithoutData(t *testing.T) {
	e := &derrors.APIError{Code: 11000, Message: "insufficient"}
	got := e.Error()
	assert.Contains(t, got, "11000")
	assert.Contains(t, got, "insufficient")
	assert.NotContains(t, got, "(")
}

func TestAPIError_Error_WithData(t *testing.T) {
	e := &derrors.APIError{
		Code:    11000,
		Message: "insufficient margin",
		Data:    []byte(`{"need":"500"}`),
	}
	got := e.Error()
	assert.Contains(t, got, "11000")
	assert.Contains(t, got, `{"need":"500"}`)
}

// --- Is() mapping per sentinel ----------------------------------------------

func TestAPIError_Is_RateLimited_StandardCode(t *testing.T) {
	e := &derrors.APIError{Code: derrors.CodeRateLimitExceeded}
	assert.True(t, derrors.Is(e, derrors.ErrRateLimited))
}

func TestAPIError_Is_RateLimited_WSConcurrencyCode(t *testing.T) {
	e := &derrors.APIError{Code: derrors.CodeConcurrentWSClientsLimitExceeded}
	assert.True(t, derrors.Is(e, derrors.ErrRateLimited))
}

func TestAPIError_Is_RateLimited_RejectsOtherCodes(t *testing.T) {
	e := &derrors.APIError{Code: derrors.CodeOrderNotExist}
	assert.False(t, derrors.Is(e, derrors.ErrRateLimited))
}

func TestAPIError_Is_Unauthorized_AllAuthCodes(t *testing.T) {
	cases := []int{
		derrors.CodeInvalidSignature,
		derrors.CodeHeaderWalletMissing,
		derrors.CodeHeaderWalletMismatch,
		derrors.CodeSignerNotOwner,
		derrors.CodeMissingPrivateParam,
		derrors.CodeSessionKeyNotFound,
		derrors.CodeSessionKeyExpired,
		derrors.CodeSessionKeyIPNotWhitelisted,
		derrors.CodeUnauthorizedKeyScope,
		derrors.CodeScopeNotAdmin,
		derrors.CodeUnauthorizedRFQMaker,
		derrors.CodeAccountNotWhitelistedAtomicOrders,
		derrors.CodeSentinelAuthInvalid,
	}
	for _, code := range cases {
		t.Run("code", func(t *testing.T) {
			e := &derrors.APIError{Code: code}
			assert.True(t, derrors.Is(e, derrors.ErrUnauthorized), "code %d", code)
		})
	}
}

func TestAPIError_Is_InvalidSignature(t *testing.T) {
	e := &derrors.APIError{Code: derrors.CodeInvalidSignature}
	assert.True(t, derrors.Is(e, derrors.ErrInvalidSignature))
	assert.True(t, derrors.Is(e, derrors.ErrUnauthorized), "should also match the broader auth bucket")
}

func TestAPIError_Is_SessionKeyExpired(t *testing.T) {
	e := &derrors.APIError{Code: derrors.CodeSessionKeyExpired}
	assert.True(t, derrors.Is(e, derrors.ErrSessionKeyExpired))
	assert.True(t, derrors.Is(e, derrors.ErrUnauthorized))
}

func TestAPIError_Is_SessionKeyNotFound(t *testing.T) {
	e := &derrors.APIError{Code: derrors.CodeSessionKeyNotFound}
	assert.True(t, derrors.Is(e, derrors.ErrSessionKeyNotFound))
}

func TestAPIError_Is_InsufficientFunds_OrderCode(t *testing.T) {
	assert.True(t, derrors.Is(
		&derrors.APIError{Code: derrors.CodeInsufficientFundsOrder},
		derrors.ErrInsufficientFunds))
}

func TestAPIError_Is_InsufficientFunds_ERC20Codes(t *testing.T) {
	for _, code := range []int{derrors.CodeERC20InsufficientAllowance, derrors.CodeERC20InsufficientBalance} {
		e := &derrors.APIError{Code: code}
		assert.True(t, derrors.Is(e, derrors.ErrInsufficientFunds), "code %d", code)
	}
}

func TestAPIError_Is_OrderNotFound(t *testing.T) {
	e := &derrors.APIError{Code: derrors.CodeOrderNotExist}
	assert.True(t, derrors.Is(e, derrors.ErrOrderNotFound))
}

func TestAPIError_Is_AlreadyCancelled(t *testing.T) {
	assert.True(t, derrors.Is(&derrors.APIError{Code: derrors.CodeAlreadyCancelled}, derrors.ErrAlreadyCancelled))
}

func TestAPIError_Is_AlreadyFilled(t *testing.T) {
	assert.True(t, derrors.Is(&derrors.APIError{Code: derrors.CodeAlreadyFilled}, derrors.ErrAlreadyFilled))
}

func TestAPIError_Is_AlreadyExpired(t *testing.T) {
	assert.True(t, derrors.Is(&derrors.APIError{Code: derrors.CodeAlreadyExpired}, derrors.ErrAlreadyExpired))
}

func TestAPIError_Is_InstrumentNotFound_BothCodes(t *testing.T) {
	for _, code := range []int{derrors.CodeInstrumentNotFound, derrors.CodeAssetNotFound} {
		e := &derrors.APIError{Code: code}
		assert.True(t, derrors.Is(e, derrors.ErrInstrumentNotFound), "code %d", code)
	}
}

func TestAPIError_Is_SubaccountNotFound(t *testing.T) {
	e := &derrors.APIError{Code: derrors.CodeSubaccountNotFound}
	assert.True(t, derrors.Is(e, derrors.ErrSubaccountNotFound))
}

func TestAPIError_Is_AccountNotFound(t *testing.T) {
	e := &derrors.APIError{Code: derrors.CodeAccountNotFound}
	assert.True(t, derrors.Is(e, derrors.ErrAccountNotFound))
}

func TestAPIError_Is_ChainIDMismatch(t *testing.T) {
	e := &derrors.APIError{Code: derrors.CodeChainIDMismatch}
	assert.True(t, derrors.Is(e, derrors.ErrChainIDMismatch))
}

func TestAPIError_Is_MMPFrozen(t *testing.T) {
	e := &derrors.APIError{Code: derrors.CodeMMPFrozen}
	assert.True(t, derrors.Is(e, derrors.ErrMMPFrozen))
}

func TestAPIError_Is_RestrictedRegion_AllArms(t *testing.T) {
	for _, code := range []int{derrors.CodeRestrictedRegion, derrors.CodeAccountDisabledCompliance} {
		e := &derrors.APIError{Code: code}
		assert.True(t, derrors.Is(e, derrors.ErrRestrictedRegion), "code %d", code)
	}
}

func TestAPIError_Is_UnknownCode_FallsThrough(t *testing.T) {
	e := &derrors.APIError{Code: 99999}
	assert.False(t, derrors.Is(e, derrors.ErrRateLimited))
	assert.False(t, derrors.Is(e, derrors.ErrUnauthorized))
	assert.False(t, derrors.Is(e, derrors.ErrInsufficientFunds))
}

func TestAPIError_Is_NonSentinelTarget(t *testing.T) {
	e := &derrors.APIError{Code: derrors.CodeRateLimitExceeded}
	assert.False(t, derrors.Is(e, stderrors.New("unrelated")))
}

func TestAPIError_Is_DefaultArm_SwitchFallthrough(t *testing.T) {
	// A code that no sentinel covers (e.g. 11050 TriggerOrderCancelled);
	// every Is target should return false.
	e := &derrors.APIError{Code: derrors.CodeTriggerOrderCancelled}
	for _, sentinel := range []error{
		derrors.ErrRateLimited, derrors.ErrUnauthorized,
		derrors.ErrInvalidSignature, derrors.ErrSessionKeyExpired,
		derrors.ErrSessionKeyNotFound, derrors.ErrInsufficientFunds,
		derrors.ErrOrderNotFound, derrors.ErrAlreadyCancelled,
		derrors.ErrAlreadyFilled, derrors.ErrAlreadyExpired,
		derrors.ErrInstrumentNotFound, derrors.ErrSubaccountNotFound,
		derrors.ErrAccountNotFound, derrors.ErrChainIDMismatch,
		derrors.ErrMMPFrozen, derrors.ErrRestrictedRegion,
	} {
		assert.False(t, derrors.Is(e, sentinel),
			"trigger-cancelled (%d) should not match %v", e.Code, sentinel)
	}
}

func TestAPIError_Implements_ErrorInterface(_ *testing.T) {
	var _ error = &derrors.APIError{}
}

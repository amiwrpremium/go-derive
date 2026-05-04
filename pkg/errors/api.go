package errors

import (
	"encoding/json"
	"fmt"
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

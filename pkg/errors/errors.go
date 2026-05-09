// Package errors provides the SDK's error types and sentinel values. All
// errors are constructed so they work with errors.Is and errors.As from the
// standard library.
package errors

import (
	"errors"

	"github.com/amiwrpremium/go-derive/internal/transport"
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

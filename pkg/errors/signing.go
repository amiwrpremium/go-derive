// Package errors provides the SDK's error types and sentinel values. All
// errors are constructed so they work with errors.Is and errors.As from the
// standard library.
package errors

import "fmt"

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

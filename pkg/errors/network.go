// Package errors provides the SDK's error types and sentinel values. All
// errors are constructed so they work with errors.Is and errors.As from the
// standard library.
package errors

import "fmt"

// ConnectionError wraps low-level transport failures so callers can
// distinguish network problems from API errors.
type ConnectionError struct {
	Op  string
	Err error
}

// Error implements the error interface.
func (e *ConnectionError) Error() string {
	return fmt.Sprintf("derive: connection: %s: %v", e.Op, e.Err)
}

// Unwrap returns the underlying transport error.
func (e *ConnectionError) Unwrap() error { return e.Err }

// TimeoutError signals a deadline expiry while waiting for a response.
type TimeoutError struct {
	Method string
}

// Error implements the error interface.
func (e *TimeoutError) Error() string {
	return fmt.Sprintf("derive: timeout waiting for response to %q", e.Method)
}

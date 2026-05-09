// Package transport — network-error types.
//
// `ConnectionError` and `TimeoutError` are declared here (rather than at
// root or in `pkg/errors`) to break the `root → transport → errors → root`
// import cycle that would otherwise form once `pkg/errors` lifts to root.
// The public-facing copies live as Go type aliases in `pkg/errors`
// (and, post-migration, in root `errors.go`); external code keeps using
// `derive.ConnectionError` / `derive.TimeoutError` exactly as before.
package transport

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

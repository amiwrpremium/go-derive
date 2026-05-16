// Package types declares the domain types used in REST and WebSocket
// requests and responses.
package types

import (
	"errors"
	"fmt"
)

// ErrInvalidParams is the sentinel returned (wrapped) by every
// `*Input.Validate()` / `*Query.Validate()` in this package when a
// caller-supplied field fails its constraint. Match it with
// errors.Is.
var ErrInvalidParams = errors.New("types: invalid params")

// invalidParam is the package-internal helper for assembling consistent
// validation errors. The returned error wraps [ErrInvalidParams] so
// callers can match without unwrapping each kind separately.
func invalidParam(field, reason string) error {
	return fmt.Errorf("%w: %s: %s", ErrInvalidParams, field, reason)
}

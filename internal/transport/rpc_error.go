// Package transport — wire-level RPC error.
//
// `JSONRPCError` is the transport-internal representation of a server-side
// JSON-RPC error envelope. The methods layer (root, post-migration) wraps
// each `*JSONRPCError` it receives into a `*derive.APIError` at the
// boundary; keeping the rich-typed APIError out of transport breaks the
// `root↔transport` import cycle that would otherwise form once
// `pkg/errors` lifts to root.
package transport

import (
	"encoding/json"
	"fmt"
)

// JSONRPCError carries the three fields of a JSON-RPC error envelope.
// Callers receive this type from transport.Call and translate it into
// the public `*derive.APIError` at the methods boundary.
type JSONRPCError struct {
	Code    int
	Message string
	Data    json.RawMessage
}

// Error implements the error interface.
func (e *JSONRPCError) Error() string {
	return fmt.Sprintf("derive: rpc error %d: %s", e.Code, e.Message)
}

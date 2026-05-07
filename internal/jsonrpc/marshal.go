// Package jsonrpc implements JSON-RPC 2.0 framing.
//
// The package is transport-agnostic: callers serialise a [Request] to
// bytes, hand them to whatever transport (HTTP body, WebSocket frame, …),
// and deserialise the reply into a [Response]. The HTTP and WebSocket
// transports in internal/transport both build on this package.
//
// The package also defines [Notification], which is the unsolicited
// frame Derive emits for subscription updates, and [IsNotification] for
// distinguishing it from a request reply on the read path.
package jsonrpc

import (
	"encoding/json"
	"fmt"
)

// DecodeResult unmarshals [Response.Result] into out, propagating server
// errors as Go errors.
//
// Behaviour by case:
//
//   - resp == nil: returns "jsonrpc: nil response".
//   - resp.Error != nil: returns the *[Error] (which satisfies the error
//     interface) and leaves out untouched.
//   - resp.Result is empty or out is nil: returns nil with no decode.
//   - otherwise: json.Unmarshal(Result, out), wrapping any error.
func DecodeResult(resp *Response, out any) error {
	if resp == nil {
		return fmt.Errorf("jsonrpc: nil response")
	}
	if resp.Error != nil {
		return resp.Error
	}
	if out == nil || len(resp.Result) == 0 {
		return nil
	}
	if err := json.Unmarshal(resp.Result, out); err != nil {
		return fmt.Errorf("jsonrpc: decode result: %w", err)
	}
	return nil
}

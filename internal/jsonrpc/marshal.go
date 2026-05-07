// Package jsonrpc — see envelope.go for the overview.
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

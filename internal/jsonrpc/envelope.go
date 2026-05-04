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

// Version is the JSON-RPC protocol version this package implements.
const Version = "2.0"

// Request is an outgoing JSON-RPC 2.0 request envelope.
type Request struct {
	// JSONRPC is the protocol version, always "2.0".
	JSONRPC string `json:"jsonrpc"`
	// ID is the caller-assigned correlation id.
	ID uint64 `json:"id"`
	// Method is the dotted method name (e.g. "public/get_instruments").
	Method string `json:"method"`
	// Params is the pre-marshalled parameter object; omitted when nil.
	Params json.RawMessage `json:"params,omitempty"`
}

// NewRequest builds a [Request] with id and method, marshalling params via
// [json.Marshal]. Pass nil for params to omit the field from the wire
// format.
//
// Returns an error if params is not JSON-marshalable.
func NewRequest(id uint64, method string, params any) (*Request, error) {
	req := &Request{JSONRPC: Version, ID: id, Method: method}
	if params == nil {
		return req, nil
	}
	raw, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("jsonrpc: marshal params: %w", err)
	}
	req.Params = raw
	return req, nil
}

// Response is an incoming JSON-RPC 2.0 response envelope.
//
// Exactly one of Result or Error is populated for a well-formed reply.
// Use [DecodeResult] to extract a typed result while propagating server
// errors as errors.
type Response struct {
	// JSONRPC is the protocol version, always "2.0".
	JSONRPC string `json:"jsonrpc"`
	// ID is the response correlation id as a raw JSON token.
	//
	// Over WebSocket Derive echoes the numeric request id, so [Response.IDUint64]
	// recovers a uint64 for dispatcher correlation. Over REST Derive
	// always replies with a freshly-generated UUID string regardless of
	// what was sent — the HTTP transport doesn't correlate by id, so the
	// raw token is preserved here without forcing a type.
	ID json.RawMessage `json:"id,omitempty"`
	// Result is the success payload as raw JSON; nil on error responses.
	Result json.RawMessage `json:"result,omitempty"`
	// Error is the failure object; nil on success responses.
	Error *Error `json:"error,omitempty"`
}

// IDUint64 returns the response id as a uint64 if it was a JSON number,
// or (0, false) if missing or not numeric (e.g. a string UUID).
//
// Used by the WebSocket transport to route replies to pending calls.
func (r *Response) IDUint64() (uint64, bool) {
	if len(r.ID) == 0 {
		return 0, false
	}
	var n uint64
	if err := json.Unmarshal(r.ID, &n); err != nil {
		return 0, false
	}
	return n, true
}

// Error is the standard JSON-RPC 2.0 error object.
//
// Code follows the spec ranges (negative numbers for protocol errors,
// otherwise application-specific). Derive uses [Data] to carry structured
// details on validation errors and rate-limit responses (the
// pkg/errors.APIError type carries the same shape with sentinel mapping).
type Error struct {
	// Code identifies the failure class.
	Code int `json:"code"`
	// Message is the human-readable summary.
	Message string `json:"message"`
	// Data is the optional structured payload.
	Data json.RawMessage `json:"data,omitempty"`
}

// Error implements the error interface, formatting the code, message and
// optional data payload.
func (e *Error) Error() string {
	if len(e.Data) > 0 {
		return fmt.Sprintf("jsonrpc error %d: %s (%s)", e.Code, e.Message, string(e.Data))
	}
	return fmt.Sprintf("jsonrpc error %d: %s", e.Code, e.Message)
}

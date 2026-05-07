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

import "encoding/json"

// Notification is an unsolicited message from the server.
//
// Derive uses notifications for subscription channel updates: Method is
// "subscription" and Params is a [SubscriptionParams] envelope. There is
// no [Request.ID] (the absence of the id field is what distinguishes a
// notification from a reply — see [IsNotification]).
type Notification struct {
	// JSONRPC is the protocol version, always "2.0".
	JSONRPC string `json:"jsonrpc"`
	// Method names the notification ("subscription" for channel updates).
	Method string `json:"method"`
	// Params is the inner payload as raw JSON.
	Params json.RawMessage `json:"params"`
}

// SubscriptionParams is the canonical shape of [Notification.Params] when
// Method is "subscription". Channel identifies the subscribed channel and
// Data carries the typed payload defined by that channel's descriptor.
type SubscriptionParams struct {
	// Channel is the dotted server-side channel name (e.g.
	// "trades.BTC-PERP").
	Channel string `json:"channel"`
	// Data is the typed payload as raw JSON; the channel descriptor's
	// Decode method turns it into a Go value.
	Data json.RawMessage `json:"data"`
}

// IsNotification reports whether a raw frame is a server-initiated
// notification rather than a reply to a request.
//
// It probes only the "id" and "method" fields, avoiding a full unmarshal —
// hot-path code on the read pump calls this on every frame.
func IsNotification(raw []byte) bool {
	var probe struct {
		ID     *uint64 `json:"id"`
		Method string  `json:"method"`
	}
	if err := json.Unmarshal(raw, &probe); err != nil {
		return false
	}
	return probe.ID == nil && probe.Method != ""
}

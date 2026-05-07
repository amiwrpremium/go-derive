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

import "sync/atomic"

// IDGen produces strictly-increasing JSON-RPC request IDs.
//
// IDs are produced by an atomic counter, so the type is safe for
// concurrent use across many goroutines. The zero value is not usable;
// construct via [NewIDGen].
type IDGen struct{ n atomic.Uint64 }

// NewIDGen returns a generator whose first [IDGen.Next] yields 1.
//
// ID 0 is deliberately skipped so callers can use it as a "missing"
// sentinel when threading IDs through internal data structures.
func NewIDGen() *IDGen {
	g := &IDGen{}
	g.n.Store(0)
	return g
}

// Next returns the next ID. It is safe for concurrent use.
func (g *IDGen) Next() uint64 { return g.n.Add(1) }

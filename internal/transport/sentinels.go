// Package transport — sentinel values used by the transport pumps.
//
// Like network_errors.go, these declarations sit here rather than in
// `pkg/errors` (or, post-migration, root) to break the
// `root↔transport` import cycle. `pkg/errors` re-declares each
// sentinel as `var X = transport.X`; the var initialiser preserves
// pointer identity so `errors.Is(err, derive.X)` matches when the
// underlying error came from this package.
package transport

import "errors"

// ErrNotConnected is returned when a WebSocket call is attempted before
// Connect() has succeeded or after the connection has terminated.
var ErrNotConnected = errors.New("derive: not connected")

// ErrAlreadyConnected is returned when Connect() is called on a client
// that is already connected.
var ErrAlreadyConnected = errors.New("derive: already connected")

// ErrSubscriptionClosed is returned by Subscription.Updates() once the
// channel has been closed by either party.
var ErrSubscriptionClosed = errors.New("derive: subscription closed")

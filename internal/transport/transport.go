// Package transport — see http.go for the overview.
package transport

import (
	"context"
	"encoding/json"
)

// Transport is the minimal interface used to issue JSON-RPC calls.
//
// Implementations are responsible for ID assignment, framing,
// rate-limiting, and request/response correlation. They must be safe for
// concurrent use; callers may share a Transport across goroutines.
type Transport interface {
	// Call issues a single JSON-RPC request and decodes the result into
	// out. Pass nil for out to discard the result body. params may be nil.
	//
	// Server-side errors arrive as
	// *[github.com/amiwrpremium/go-derive/pkg/errors.APIError]. Transport
	// failures arrive as
	// *[github.com/amiwrpremium/go-derive/pkg/errors.ConnectionError].
	Call(ctx context.Context, method string, params, out any) error

	// Close releases any resources held by the transport. After Close the
	// implementation may treat further Call invocations as errors.
	Close() error
}

// Decoder turns a raw notification payload (the bytes inside
// [github.com/amiwrpremium/go-derive/internal/jsonrpc.SubscriptionParams.Data])
// into a typed value. Channel descriptors in pkg/channels supply concrete
// decoders.
type Decoder func(raw json.RawMessage) (any, error)

// Subscription is the consumer-facing handle returned by
// [Subscriber.Subscribe].
//
// Updates flow on the channel returned by [Subscription.Updates] until the
// subscription terminates — either by an explicit [Subscription.Close] or
// by the underlying connection failing. After Updates closes, call
// [Subscription.Err] to learn why.
type Subscription interface {
	// Channel returns the dotted server-side channel name (e.g.
	// "trades.BTC-PERP").
	Channel() string

	// Updates returns a receive-only channel of decoded events. The
	// channel is closed when the subscription ends; check [Subscription.Err]
	// to learn the terminal error (or nil for a clean close).
	Updates() <-chan any

	// Close terminates the subscription and best-effort sends an
	// unsubscribe RPC to the server.
	Close() error

	// Err returns the terminal error after [Subscription.Updates] is
	// closed. It returns nil for a clean shutdown.
	Err() error
}

// Subscriber is implemented by transports that support pub/sub.
//
// The HTTP transport does not satisfy this interface — it has no way to
// deliver server-initiated notifications. The WebSocket transport does.
type Subscriber interface {
	// Subscribe registers a server-side channel and returns a
	// [Subscription] whose Updates yield values produced by decode.
	Subscribe(ctx context.Context, channel string, decode Decoder) (Subscription, error)
}

// Package ws is the WebSocket-backed client for Derive's JSON-RPC API.
//
// This file declares the functional-options surface for [Subscribe]
// and [SubscribeFunc] — buffer size, drop policy, and an error
// handler for events the pump cannot deliver.
package ws

import "errors"

// DropPolicy controls what [Subscribe] does when the per-subscription
// buffer is full and a new event arrives.
type DropPolicy int

const (
	// DropNewest discards the incoming event when the buffer is full.
	// The default. Suits pub/sub workloads where a slow consumer
	// must not back-pressure the WebSocket read pump and starve other
	// subscriptions on the same client.
	DropNewest DropPolicy = iota

	// DropOldest discards the oldest buffered event and pushes the
	// new one. Best-effort — under contention the pump may briefly
	// race with the consumer and end up dropping the new event
	// instead. Suits live-state feeds where freshness matters more
	// than completeness.
	DropOldest

	// Block back-pressures the read pump until the consumer reads.
	// Use with care: a stalled consumer will stall every subscription
	// on the same client because the WebSocket read pump is shared.
	// Suits workloads that must not lose events and can tolerate the
	// back-pressure trade-off.
	Block
)

// ErrBufferFull is reported to [WithErrorHandler] when the
// per-subscription buffer is full and the configured [DropPolicy]
// caused an event to be dropped.
var ErrBufferFull = errors.New("ws: subscription buffer full; event dropped")

// ErrTypeMismatch is reported to [WithErrorHandler] when the
// underlying decoder produced a value the typed pump could not
// assert into T. Usually a bug in the caller (wrong T for the
// channel) or a docs/wire schema drift.
var ErrTypeMismatch = errors.New("ws: subscription type mismatch; event dropped")

// SubscribeOption tunes a single Subscribe / SubscribeFunc call. The
// zero default is bufferSize=256, dropPolicy=DropNewest, no error
// handler.
type SubscribeOption func(*subscribeConfig)

type subscribeConfig struct {
	bufferSize   int
	dropPolicy   DropPolicy
	errorHandler func(error)
}

const defaultBufferSize = 256

func defaultSubscribeConfig() subscribeConfig {
	return subscribeConfig{
		bufferSize: defaultBufferSize,
		dropPolicy: DropNewest,
	}
}

func applySubscribeOpts(opts []SubscribeOption) subscribeConfig {
	cfg := defaultSubscribeConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.bufferSize <= 0 {
		cfg.bufferSize = defaultBufferSize
	}
	return cfg
}

// WithBufferSize sets the in-memory event buffer for the typed
// channel. Values <= 0 fall back to the default (256). Larger
// buffers absorb consumer pauses at the cost of memory.
func WithBufferSize(n int) SubscribeOption {
	return func(c *subscribeConfig) { c.bufferSize = n }
}

// WithDropPolicy sets the policy used when the buffer is full. See
// [DropPolicy] for the trade-offs.
func WithDropPolicy(p DropPolicy) SubscribeOption {
	return func(c *subscribeConfig) { c.dropPolicy = p }
}

// WithErrorHandler installs a callback invoked whenever the pump
// cannot deliver an event — buffer-full drops (see [ErrBufferFull])
// and type-assertion mismatches (see [ErrTypeMismatch]).
//
// The handler runs synchronously on the read-pump goroutine; keep
// it non-blocking. To process errors asynchronously, push them onto
// a buffered channel from inside the handler.
//
// Terminal errors that close the subscription are not delivered
// here — read [Subscription.Err] after [Subscription.Updates] closes.
func WithErrorHandler(fn func(error)) SubscribeOption {
	return func(c *subscribeConfig) { c.errorHandler = fn }
}

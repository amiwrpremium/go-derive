package testutil

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/amiwrpremium/go-derive/internal/transport"
)

// FakeTransport implements transport.Transport without any network access.
// Tests register handlers per-method and inspect captured calls afterwards.
type FakeTransport struct {
	mu       sync.Mutex
	handlers map[string]FakeHandler
	calls    []FakeCall
	closed   bool
}

// FakeCall is a record of one Call invocation.
type FakeCall struct {
	Method string
	Params json.RawMessage
}

// FakeHandler returns the result and/or error for one method.
//
// To reply with a JSON-RPC application error, return ( nil, &derrors.APIError{...} ).
// If the handler is unset, Call returns an "unhandled method" error.
type FakeHandler func(params json.RawMessage) (any, error)

// NewFakeTransport returns an empty FakeTransport.
func NewFakeTransport() *FakeTransport {
	return &FakeTransport{handlers: map[string]FakeHandler{}}
}

// Handle registers a handler for one method name. Replaces any previous
// handler for that method.
func (f *FakeTransport) Handle(method string, h FakeHandler) {
	f.mu.Lock()
	f.handlers[method] = h
	f.mu.Unlock()
}

// HandleResult is a shortcut for Handle that always returns the given value.
func (f *FakeTransport) HandleResult(method string, result any) {
	f.Handle(method, func(json.RawMessage) (any, error) { return result, nil })
}

// HandleError is a shortcut for Handle that always returns the given error.
func (f *FakeTransport) HandleError(method string, err error) {
	f.Handle(method, func(json.RawMessage) (any, error) { return nil, err })
}

// Calls returns a snapshot of recorded calls in order.
func (f *FakeTransport) Calls() []FakeCall {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]FakeCall, len(f.calls))
	copy(out, f.calls)
	return out
}

// LastCall returns the most recent call, or zero-value if none.
func (f *FakeTransport) LastCall() FakeCall {
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.calls) == 0 {
		return FakeCall{}
	}
	return f.calls[len(f.calls)-1]
}

// Reset clears recorded calls but keeps registered handlers.
func (f *FakeTransport) Reset() {
	f.mu.Lock()
	f.calls = nil
	f.mu.Unlock()
}

// Call satisfies transport.Transport.
func (f *FakeTransport) Call(_ context.Context, method string, params, out any) error {
	f.mu.Lock()
	if f.closed {
		f.mu.Unlock()
		return fmt.Errorf("faketransport: closed")
	}

	var raw json.RawMessage
	if params != nil {
		b, err := json.Marshal(params)
		if err != nil {
			f.mu.Unlock()
			return fmt.Errorf("faketransport: marshal params: %w", err)
		}
		raw = b
	}
	f.calls = append(f.calls, FakeCall{Method: method, Params: raw})
	h, ok := f.handlers[method]
	f.mu.Unlock()

	if !ok {
		return fmt.Errorf("faketransport: unhandled method %q", method)
	}
	result, err := h(raw)
	if err != nil {
		return err
	}
	if out == nil || result == nil {
		return nil
	}
	rb, mErr := json.Marshal(result)
	if mErr != nil {
		return fmt.Errorf("faketransport: marshal result: %w", mErr)
	}
	return json.Unmarshal(rb, out)
}

// Close satisfies transport.Transport.
func (f *FakeTransport) Close() error {
	f.mu.Lock()
	f.closed = true
	f.mu.Unlock()
	return nil
}

// compile-time interface check
var _ transport.Transport = (*FakeTransport)(nil)

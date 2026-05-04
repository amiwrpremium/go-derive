package testutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/amiwrpremium/go-derive/internal/jsonrpc"
)

// MockWSServer is an in-process JSON-RPC WebSocket server for testing.
// Register handlers per method; push notifications via Notify.
type MockWSServer struct {
	server   *httptest.Server
	mu       sync.Mutex
	handlers map[string]MockWSHandler
	conns    []*mockConn
	closed   bool
}

// MockWSHandler returns the typed response for one JSON-RPC method.
// Returning err != nil writes a JSON-RPC error frame.
type MockWSHandler func(params json.RawMessage) (result any, err *jsonrpc.Error)

// mockConn is one accepted connection. The writeMu serialises all writes
// (data and control) through one goroutine's worth of write calls —
// gorilla/websocket forbids concurrent WriteMessage calls.
type mockConn struct {
	c       *websocket.Conn
	subs    map[string]bool
	mu      sync.Mutex
	writeMu sync.Mutex
}

// upgrader is shared across connections; gorilla doesn't make Upgrader
// goroutine-safe explicitly, but Upgrade only reads its fields, so
// concurrent calls are safe in practice and standard practice.
var upgrader = websocket.Upgrader{
	// Allow any origin; this is a test server.
	CheckOrigin: func(*http.Request) bool { return true },
	// Match the client side default (8 MiB).
	ReadBufferSize:  8 << 10,
	WriteBufferSize: 8 << 10,
}

// NewMockWSServer starts the server. Defer Close().
func NewMockWSServer() *MockWSServer {
	m := &MockWSServer{handlers: map[string]MockWSHandler{}}
	m.server = httptest.NewServer(http.HandlerFunc(m.handle))
	return m
}

// URL returns the ws:// URL for clients to dial.
func (m *MockWSServer) URL() string {
	return strings.Replace(m.server.URL, "http://", "ws://", 1)
}

// Close shuts down all connections and the underlying HTTP server.
func (m *MockWSServer) Close() {
	m.mu.Lock()
	m.closed = true
	conns := m.conns
	m.conns = nil
	m.mu.Unlock()
	for _, c := range conns {
		c.writeMu.Lock()
		_ = c.c.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseGoingAway, "shutdown"),
			time.Now().Add(time.Second))
		c.writeMu.Unlock()
		_ = c.c.Close()
	}
	m.server.Close()
}

// Handle registers a handler for one method.
func (m *MockWSServer) Handle(method string, h MockWSHandler) {
	m.mu.Lock()
	m.handlers[method] = h
	m.mu.Unlock()
}

// HandleResult is a shortcut for Handle that always returns result.
func (m *MockWSServer) HandleResult(method string, result any) {
	m.Handle(method, func(json.RawMessage) (any, *jsonrpc.Error) { return result, nil })
}

// Subscribed reports whether at least one connection has subscribed to channel.
func (m *MockWSServer) Subscribed(channel string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, c := range m.conns {
		c.mu.Lock()
		ok := c.subs[channel]
		c.mu.Unlock()
		if ok {
			return true
		}
	}
	return false
}

// WaitSubscribed blocks until any connection has subscribed to channel
// or the timeout fires.
func (m *MockWSServer) WaitSubscribed(channel string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if m.Subscribed(channel) {
			return true
		}
		time.Sleep(5 * time.Millisecond)
	}
	return false
}

// Notify pushes a subscription notification to every connection currently
// subscribed to the channel.
func (m *MockWSServer) Notify(channel string, data any) {
	raw, err := json.Marshal(data)
	if err != nil {
		return
	}
	notif := jsonrpc.Notification{
		JSONRPC: jsonrpc.Version,
		Method:  "subscription",
	}
	params, _ := json.Marshal(jsonrpc.SubscriptionParams{Channel: channel, Data: raw})
	notif.Params = params
	frame, _ := json.Marshal(notif)

	m.mu.Lock()
	conns := append([]*mockConn{}, m.conns...)
	m.mu.Unlock()
	for _, c := range conns {
		c.mu.Lock()
		ok := c.subs[channel]
		c.mu.Unlock()
		if !ok {
			continue
		}
		c.writeMu.Lock()
		_ = c.c.SetWriteDeadline(time.Now().Add(2 * time.Second))
		_ = c.c.WriteMessage(websocket.TextMessage, frame)
		c.writeMu.Unlock()
	}
}

// handle is the HTTP upgrade handler.
func (m *MockWSServer) handle(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	mc := &mockConn{c: c, subs: map[string]bool{}}
	m.mu.Lock()
	m.conns = append(m.conns, mc)
	m.mu.Unlock()

	defer func() {
		mc.writeMu.Lock()
		_ = c.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			time.Now().Add(time.Second))
		mc.writeMu.Unlock()
		_ = c.Close()
	}()

	for {
		_, data, err := c.ReadMessage()
		if err != nil {
			return
		}

		var req jsonrpc.Request
		if err := json.Unmarshal(data, &req); err != nil {
			continue
		}

		resp := m.handleOne(mc, &req)
		out, _ := json.Marshal(resp)
		mc.writeMu.Lock()
		_ = c.SetWriteDeadline(time.Now().Add(2 * time.Second))
		_ = c.WriteMessage(websocket.TextMessage, out)
		mc.writeMu.Unlock()
	}
}

func (m *MockWSServer) handleOne(mc *mockConn, req *jsonrpc.Request) jsonrpc.Response {
	idBytes, _ := json.Marshal(req.ID) // json.Marshal of uint64 never errors
	resp := jsonrpc.Response{JSONRPC: jsonrpc.Version, ID: idBytes}

	// A custom handler — registered via Handle — takes precedence over the
	// built-in subscribe/unsubscribe defaults. This lets tests simulate
	// server-side errors on subscribe.
	m.mu.Lock()
	override, hasOverride := m.handlers[req.Method]
	m.mu.Unlock()
	if hasOverride && (req.Method == "subscribe" || req.Method == "unsubscribe") {
		v, errResp := override(req.Params)
		if errResp != nil {
			resp.Error = errResp
			return resp
		}
		if v != nil {
			if b, err := json.Marshal(v); err == nil {
				resp.Result = b
			}
		} else {
			resp.Result = json.RawMessage("null")
		}
		return resp
	}

	// Built-in subscribe / unsubscribe handling.
	if req.Method == "subscribe" {
		var p struct {
			Channels []string `json:"channels"`
		}
		_ = json.Unmarshal(req.Params, &p)
		mc.mu.Lock()
		for _, ch := range p.Channels {
			mc.subs[ch] = true
		}
		mc.mu.Unlock()
		resp.Result, _ = json.Marshal(map[string]any{
			"status":                map[string]string{},
			"current_subscriptions": p.Channels,
		})
		return resp
	}
	if req.Method == "unsubscribe" {
		var p struct {
			Channels []string `json:"channels"`
		}
		_ = json.Unmarshal(req.Params, &p)
		mc.mu.Lock()
		for _, ch := range p.Channels {
			delete(mc.subs, ch)
		}
		mc.mu.Unlock()
		resp.Result, _ = json.Marshal(map[string]any{"status": map[string]string{}})
		return resp
	}

	m.mu.Lock()
	h, ok := m.handlers[req.Method]
	m.mu.Unlock()
	if !ok {
		resp.Error = &jsonrpc.Error{Code: -32601, Message: "Method not found: " + req.Method}
		return resp
	}
	v, errResp := h(req.Params)
	if errResp != nil {
		resp.Error = errResp
		return resp
	}
	if v != nil {
		b, err := json.Marshal(v)
		if err != nil {
			resp.Error = &jsonrpc.Error{Code: -32603, Message: err.Error()}
			return resp
		}
		resp.Result = b
	} else {
		resp.Result = json.RawMessage("null")
	}
	return resp
}

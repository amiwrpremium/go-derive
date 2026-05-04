// Package testutil exposes test fixtures and mock servers used by the SDK's
// internal tests. It is internal-only.
package testutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"

	"github.com/amiwrpremium/go-derive/internal/jsonrpc"
)

// MockServer is a tiny JSON-RPC HTTP server. Register one or more handlers
// per method name; unhandled methods return MethodNotFound. The HTTP method
// must always be POST.
type MockServer struct {
	server   *httptest.Server
	mu       sync.Mutex
	handlers map[string]MockHandler
	requests []MockRequest
}

// MockRequest is a record of one received call.
type MockRequest struct {
	Method  string
	Params  json.RawMessage
	Headers http.Header
}

// MockHandler returns either a result value (any json-marshalable) or an error.
type MockHandler func(req MockRequest) (any, *jsonrpc.Error)

// NewMockServer starts an httptest server and returns it.
func NewMockServer() *MockServer {
	m := &MockServer{handlers: map[string]MockHandler{}}
	m.server = httptest.NewServer(http.HandlerFunc(m.handle))
	return m
}

// URL returns the server's base URL.
func (m *MockServer) URL() string { return m.server.URL }

// Close shuts down the server.
func (m *MockServer) Close() { m.server.Close() }

// Handle registers a handler for one JSON-RPC method.
func (m *MockServer) Handle(method string, h MockHandler) {
	m.mu.Lock()
	m.handlers[method] = h
	m.mu.Unlock()
}

// Requests returns a snapshot of received requests.
func (m *MockServer) Requests() []MockRequest {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]MockRequest, len(m.requests))
	copy(out, m.requests)
	return out
}

func (m *MockServer) handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Derive's REST API routes by URL path: `<base>/public/get_time` etc.
	// The body is now the params object (not a JSON-RPC envelope), so the
	// method comes from the path.
	method := strings.TrimPrefix(r.URL.Path, "/")
	var params json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	m.mu.Lock()
	h, ok := m.handlers[method]
	m.requests = append(m.requests, MockRequest{Method: method, Params: params, Headers: r.Header.Clone()})
	m.mu.Unlock()

	resp := jsonrpc.Response{JSONRPC: "2.0", ID: json.RawMessage(`"mock-id"`)}
	if !ok {
		resp.Error = &jsonrpc.Error{Code: -32601, Message: "Method not found: " + method}
	} else {
		v, errResp := h(MockRequest{Method: method, Params: params, Headers: r.Header.Clone()})
		if errResp != nil {
			resp.Error = errResp
		} else {
			b, err := json.Marshal(v)
			if err != nil {
				resp.Error = &jsonrpc.Error{Code: -32603, Message: err.Error()}
			} else {
				resp.Result = b
			}
		}
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

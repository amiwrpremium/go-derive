// Package transport defines the JSON-RPC wire interface that pkg/rest and
// pkg/ws share, plus the HTTP and WebSocket implementations that satisfy
// it.
//
// # Layered design
//
// pkg/rest and pkg/ws both consume a [Transport] through the embedded
// API struct. The same method definition (e.g.
// [github.com/amiwrpremium/go-derive.API.GetInstruments])
// works against either transport because the only thing it needs is
// [Transport.Call].
//
// The WebSocket transport additionally implements [Subscriber] so the
// pkg/ws layer can drive subscriptions without a separate connection.
package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/amiwrpremium/go-derive/internal/jsonrpc"
)

// HTTPHeaderProvider is called once per request and returns the headers
// to attach (in addition to Content-Type / Accept / User-Agent which the
// transport sets automatically).
//
// The auth layer uses this to inject per-request EIP-191 signed headers.
// May return nil with no error to skip header injection on a given call.
type HTTPHeaderProvider func(ctx context.Context, method string, body []byte) (http.Header, error)

// HTTPTransport is a JSON-RPC transport over HTTP POST. Each Call performs
// one round-trip; there is no multiplexing or keep-alive correlation.
//
// Derive's REST API is HTTP-with-method-in-URL: every call POSTs the JSON
// body to `<base>/<method>` (e.g. `https://api.lyra.finance/public/get_time`).
// The body uses the JSON-RPC envelope so that the auth signing path can
// share one canonical request shape with the WebSocket transport — Derive
// ignores `jsonrpc`/`method` fields in the body and routes by URL.
type HTTPTransport struct {
	url     string
	client  *http.Client
	idgen   *jsonrpc.IDGen
	limiter *RateLimiter
	headers HTTPHeaderProvider
	ua      string
}

// HTTPConfig configures a new [HTTPTransport].
type HTTPConfig struct {
	// URL is the absolute JSON-RPC endpoint URL. Required.
	URL string
	// Client is the HTTP client to use. Optional; defaults to one with a
	// 30-second timeout.
	Client *http.Client
	// UserAgent is sent as the User-Agent request header. Optional.
	UserAgent string
	// Limiter is the rate limiter; nil disables limiting.
	Limiter *RateLimiter
	// Headers, when non-nil, is consulted on every request to inject auth
	// headers. Pass nil for public-only access.
	Headers HTTPHeaderProvider
}

// NewHTTP returns an HTTP-backed [Transport].
//
// URL must be absolute; client is optional and defaults to a *http.Client
// with a 30-second Timeout.
func NewHTTP(cfg HTTPConfig) (*HTTPTransport, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("transport: HTTP url is required")
	}
	c := cfg.Client
	if c == nil {
		c = &http.Client{Timeout: 30 * time.Second}
	}
	return &HTTPTransport{
		url:     cfg.URL,
		client:  c,
		idgen:   jsonrpc.NewIDGen(),
		limiter: cfg.Limiter,
		headers: cfg.Headers,
		ua:      cfg.UserAgent,
	}, nil
}

// Call issues a JSON-RPC request and decodes the result into out.
//
// The flow:
//   - block on the rate limiter (if configured) until a token is available
//   - assign a fresh request ID and marshal the body
//   - build the HTTP POST and inject auth headers via [HTTPHeaderProvider]
//   - send the request and parse the response
//   - on success decode the result; on JSON-RPC error return a
//     [APIError]
//   - on transport failure return a
//     [ConnectionError]
func (t *HTTPTransport) Call(ctx context.Context, method string, params, out any) error {
	if err := t.limiter.Wait(ctx); err != nil {
		return err
	}

	// Derive's REST API rejects unknown body fields strictly; it does NOT
	// accept the JSON-RPC envelope on most endpoints. Send only the params
	// object as the body, with `{}` substituted when params is empty.
	req, err := jsonrpc.NewRequest(t.idgen.Next(), method, params)
	if err != nil {
		return err
	}
	body := []byte(req.Params)
	if len(body) == 0 {
		body = []byte("{}")
	}

	// Derive routes by URL path, not by JSON-RPC method field. Each call
	// must POST to `<base>/<method>` (e.g. `<base>/public/get_time`).
	fullURL, err := url.JoinPath(t.url, method)
	if err != nil {
		return &ConnectionError{Op: "build url", Err: err}
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(body))
	if err != nil {
		return &ConnectionError{Op: "build request", Err: err}
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	if t.ua != "" {
		httpReq.Header.Set("User-Agent", t.ua)
	}
	if t.headers != nil {
		extra, hErr := t.headers(ctx, method, body)
		if hErr != nil {
			return hErr
		}
		for k, v := range extra {
			for _, vv := range v {
				httpReq.Header.Add(k, vv)
			}
		}
	}

	resp, err := t.client.Do(httpReq)
	if err != nil {
		return &ConnectionError{Op: "do request", Err: err}
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &ConnectionError{Op: "read response", Err: err}
	}

	var rpcResp jsonrpc.Response
	if err := json.Unmarshal(respBody, &rpcResp); err != nil {
		return fmt.Errorf("transport: decode response (status %d): %w: body=%s",
			resp.StatusCode, err, string(respBody))
	}
	if rpcResp.Error != nil {
		return &JSONRPCError{
			Code:    rpcResp.Error.Code,
			Message: rpcResp.Error.Message,
			Data:    rpcResp.Error.Data,
		}
	}
	return jsonrpc.DecodeResult(&rpcResp, out)
}

// Close is a no-op for the HTTP transport — *http.Client cleans up
// connections itself. Implemented to satisfy [Transport].
func (*HTTPTransport) Close() error { return nil }

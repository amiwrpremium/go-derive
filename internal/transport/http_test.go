package transport_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/internal/jsonrpc"
	"github.com/amiwrpremium/go-derive/internal/transport"
	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
)

// rpcServer spins up a small Derive-shaped HTTP server with a handler hook.
//
// Derive's REST routes by URL path (`<base>/public/get_time`) and the body
// is just the params object — no JSON-RPC envelope. The fixture mirrors
// that: it parses `Method` from the URL path and `Params` from the body.
func rpcServer(t *testing.T, handle func(req jsonrpc.Request, hdr http.Header) (any, *jsonrpc.Error)) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		method := strings.TrimPrefix(r.URL.Path, "/")
		var params json.RawMessage
		require.NoError(t, json.NewDecoder(r.Body).Decode(&params))
		v, errResp := handle(jsonrpc.Request{Method: method, Params: params}, r.Header.Clone())
		resp := jsonrpc.Response{JSONRPC: jsonrpc.Version, ID: json.RawMessage(`"mock-id"`)}
		if errResp != nil {
			resp.Error = errResp
		} else if v != nil {
			b, _ := json.Marshal(v)
			resp.Result = b
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
}

func TestHTTPTransport_Success(t *testing.T) {
	srv := rpcServer(t, func(req jsonrpc.Request, _ http.Header) (any, *jsonrpc.Error) {
		assert.Equal(t, "public/get_time", req.Method)
		return 1700000000000, nil
	})
	defer srv.Close()

	tt, err := transport.NewHTTP(transport.HTTPConfig{URL: srv.URL})
	require.NoError(t, err)

	var got int64
	require.NoError(t, tt.Call(context.Background(), "public/get_time", nil, &got))
	assert.Equal(t, int64(1700000000000), got)
	require.NoError(t, tt.Close())
}

func TestHTTPTransport_APIError(t *testing.T) {
	srv := rpcServer(t, func(_ jsonrpc.Request, _ http.Header) (any, *jsonrpc.Error) {
		return nil, &jsonrpc.Error{Code: 10002, Message: "rate"}
	})
	defer srv.Close()

	tt, err := transport.NewHTTP(transport.HTTPConfig{URL: srv.URL})
	require.NoError(t, err)
	err = tt.Call(context.Background(), "x", nil, nil)
	require.Error(t, err)
	var apiErr *derrors.APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, 10002, apiErr.Code)
}

func TestHTTPTransport_HeaderProvider(t *testing.T) {
	var seen atomic.Value
	srv := rpcServer(t, func(_ jsonrpc.Request, h http.Header) (any, *jsonrpc.Error) {
		seen.Store(h)
		return "ok", nil
	})
	defer srv.Close()

	hdrs := func(_ context.Context, method string, _ []byte) (http.Header, error) {
		h := http.Header{}
		h.Set("X-LyraWallet", "0xabc")
		h.Set("X-Test-Method", method)
		return h, nil
	}
	tt, err := transport.NewHTTP(transport.HTTPConfig{URL: srv.URL, Headers: hdrs, UserAgent: "go-derive-test/1"})
	require.NoError(t, err)

	var got string
	require.NoError(t, tt.Call(context.Background(), "public/get_time", nil, &got))
	hdr, ok := seen.Load().(http.Header)
	require.True(t, ok)
	assert.Equal(t, "0xabc", hdr.Get("X-LyraWallet"))
	assert.Equal(t, "public/get_time", hdr.Get("X-Test-Method"))
	assert.Equal(t, "go-derive-test/1", hdr.Get("User-Agent"))
}

func TestHTTPTransport_HeaderProviderError(t *testing.T) {
	srv := rpcServer(t, func(_ jsonrpc.Request, _ http.Header) (any, *jsonrpc.Error) { return "ok", nil })
	defer srv.Close()
	tt, err := transport.NewHTTP(transport.HTTPConfig{
		URL: srv.URL,
		Headers: func(context.Context, string, []byte) (http.Header, error) {
			return nil, errors.New("hdr boom")
		},
	})
	require.NoError(t, err)
	err = tt.Call(context.Background(), "x", nil, nil)
	assert.ErrorContains(t, err, "hdr boom")
}

func TestHTTPTransport_ServerNotJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("definitely not json"))
	}))
	defer srv.Close()
	tt, err := transport.NewHTTP(transport.HTTPConfig{URL: srv.URL})
	require.NoError(t, err)
	err = tt.Call(context.Background(), "x", nil, nil)
	assert.Error(t, err)
}

func TestHTTPTransport_NetworkErrorWrapped(t *testing.T) {
	tt, err := transport.NewHTTP(transport.HTTPConfig{
		URL:    "http://127.0.0.1:1", // nothing listens here
		Client: &http.Client{Timeout: 200 * time.Millisecond},
	})
	require.NoError(t, err)
	err = tt.Call(context.Background(), "x", nil, nil)
	require.Error(t, err)
	var connErr *derrors.ConnectionError
	assert.True(t, errors.As(err, &connErr))
}

func TestHTTPTransport_RequiresURL(t *testing.T) {
	_, err := transport.NewHTTP(transport.HTTPConfig{})
	assert.Error(t, err)
}

func TestHTTPTransport_RateLimited(t *testing.T) {
	hits := atomic.Int64{}
	srv := rpcServer(t, func(_ jsonrpc.Request, _ http.Header) (any, *jsonrpc.Error) {
		hits.Add(1)
		return "ok", nil
	})
	defer srv.Close()

	rl := transport.NewRateLimiter(1, 1)
	tt, err := transport.NewHTTP(transport.HTTPConfig{URL: srv.URL, Limiter: rl})
	require.NoError(t, err)

	require.NoError(t, tt.Call(context.Background(), "x", nil, nil))

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	err = tt.Call(ctx, "x", nil, nil)
	assert.True(t, strings.Contains(err.Error(), "deadline") || errors.Is(err, context.DeadlineExceeded))
	assert.Equal(t, int64(1), hits.Load(), "second call should never have hit the server")
}

func TestHTTPTransport_DiscardsResultWhenOutNil(t *testing.T) {
	srv := rpcServer(t, func(_ jsonrpc.Request, _ http.Header) (any, *jsonrpc.Error) {
		return map[string]string{"x": "y"}, nil
	})
	defer srv.Close()
	tt, err := transport.NewHTTP(transport.HTTPConfig{URL: srv.URL})
	require.NoError(t, err)
	require.NoError(t, tt.Call(context.Background(), "x", nil, nil))
}

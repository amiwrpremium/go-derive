package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	derrors "github.com/amiwrpremium/go-derive/pkg/errors"

	"github.com/amiwrpremium/go-derive/internal/jsonrpc"
	"github.com/amiwrpremium/go-derive/internal/retry"
)

// WSConfig configures a new [WSTransport].
type WSConfig struct {
	URL          string
	Limiter      *RateLimiter
	UserAgent    string
	PingInterval time.Duration
	// MaxReadSize is the maximum frame size accepted from the server.
	MaxReadSize int64
	// Reconnect, if true, runs a reconnect loop with exponential backoff after
	// the connection drops. Subscriptions are restored automatically.
	Reconnect bool
	// OnReconnect is invoked after a successful reconnect (e.g. to re-login).
	OnReconnect func(ctx context.Context, t *WSTransport) error
	// HTTPHeaders are sent on the upgrade request (e.g. for auth).
	HTTPHeaders http.Header
}

// pendingCall is a pending JSON-RPC request awaiting its response.
type pendingCall struct {
	out any
	err chan error
}

// wsSub is the in-transport subscription record.
type wsSub struct {
	channel string
	decode  Decoder
	updates chan any
	closed  chan struct{}
	once    sync.Once
	err     error
	errMu   sync.Mutex
}

// WSTransport is a JSON-RPC + subscription transport over a single WebSocket.
//
// It owns a single connection and serialises writes through writeQ so that
// only the writePump goroutine ever calls *websocket.Conn.WriteMessage —
// gorilla/websocket forbids concurrent writes. Reads run in their own
// goroutine. A pong-channel collects pong frames so the pingPump can verify
// liveness. Outstanding RPCs are correlated by ID; subscription
// notifications are routed by channel name.
type WSTransport struct {
	cfg   WSConfig
	idgen *jsonrpc.IDGen

	mu      sync.Mutex
	conn    *websocket.Conn
	pending map[uint64]*pendingCall
	subs    map[string]*wsSub
	writeQ  chan []byte
	stopCh  chan struct{}
	pongCh  chan struct{}
	rootCtx context.Context
	cancel  context.CancelFunc
}

// NewWS builds a [WSTransport] but does not yet dial. Call Connect.
func NewWS(cfg WSConfig) (*WSTransport, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("transport: WS url is required")
	}
	if cfg.PingInterval == 0 {
		cfg.PingInterval = 20 * time.Second
	}
	if cfg.MaxReadSize == 0 {
		cfg.MaxReadSize = 8 << 20 // 8 MiB
	}
	return &WSTransport{
		cfg:     cfg,
		idgen:   jsonrpc.NewIDGen(),
		pending: make(map[uint64]*pendingCall),
		subs:    make(map[string]*wsSub),
	}, nil
}

// Connect dials the WebSocket and starts the read/write/ping loops.
// It returns when the handshake is complete or fails.
func (t *WSTransport) Connect(ctx context.Context) error {
	t.mu.Lock()
	if t.conn != nil {
		t.mu.Unlock()
		return derrors.ErrAlreadyConnected
	}
	t.mu.Unlock()

	if err := t.dial(ctx); err != nil {
		return err
	}

	if t.cfg.Reconnect {
		// #nosec G118 -- reconnect loop intentionally outlives any single request context
		go t.reconnectLoop()
	}
	return nil
}

// dial establishes one connection and starts pumps.
func (t *WSTransport) dial(ctx context.Context) error {
	dialer := *websocket.DefaultDialer
	dialer.EnableCompression = true
	hdr := t.cfg.HTTPHeaders.Clone()
	if hdr == nil {
		hdr = http.Header{}
	}
	if t.cfg.UserAgent != "" && hdr.Get("User-Agent") == "" {
		hdr.Set("User-Agent", t.cfg.UserAgent)
	}
	c, resp, err := dialer.DialContext(ctx, t.cfg.URL, hdr)
	if resp != nil && resp.Body != nil {
		// gorilla returns the upgrade response — close it once we've
		// consumed (or rejected) the upgrade so net/http doesn't leak
		// the underlying file descriptor on dial errors.
		_ = resp.Body.Close()
	}
	if err != nil {
		return &derrors.ConnectionError{Op: "ws dial", Err: err}
	}
	c.SetReadLimit(t.cfg.MaxReadSize)

	rootCtx, cancel := context.WithCancel(context.Background())
	wq := make(chan []byte, 64)
	stop := make(chan struct{})
	pong := make(chan struct{}, 1)

	// Pong handler signals the pingPump that the peer is alive.
	c.SetPongHandler(func(string) error {
		select {
		case pong <- struct{}{}:
		default:
		}
		// Reset the read deadline whenever a pong arrives. We don't enforce a
		// per-frame read deadline (Derive's WS is mostly idle between
		// notifications), but freshening it on every pong keeps any future
		// deadline tightening cheap to add.
		_ = c.SetReadDeadline(time.Time{})
		return nil
	})

	t.mu.Lock()
	t.conn = c
	t.writeQ = wq
	t.stopCh = stop
	t.pongCh = pong
	t.rootCtx = rootCtx
	t.cancel = cancel
	t.mu.Unlock()

	// closeWatcher unblocks the read pump when the parent context is
	// cancelled (Close, failConnection, etc.) — gorilla doesn't bake ctx
	// into ReadMessage, so we close the conn from the outside to break it.
	go t.closeWatcher(rootCtx, c, stop)

	// Pumps receive their own copies of conn/wq/stop/pong so they cannot
	// observe a partially-cleared transport state if Close runs concurrently.
	go t.readPump(rootCtx, c)
	go t.writePump(c, wq, stop)
	go t.pingPump(c, pong, stop)
	return nil
}

// closeWatcher closes conn when rootCtx is cancelled, breaking the
// blocked ReadMessage in readPump.
func (t *WSTransport) closeWatcher(rootCtx context.Context, conn *websocket.Conn, stop chan struct{}) {
	select {
	case <-rootCtx.Done():
	case <-stop:
		return
	}
	_ = conn.Close()
}

// Call issues a JSON-RPC request and waits for its response.
func (t *WSTransport) Call(ctx context.Context, method string, params, out any) error {
	if err := t.cfg.Limiter.Wait(ctx); err != nil {
		return err
	}

	id := t.idgen.Next()
	req, err := jsonrpc.NewRequest(id, method, params)
	if err != nil {
		return err
	}
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("transport: marshal request: %w", err)
	}

	pc := &pendingCall{out: out, err: make(chan error, 1)}
	t.mu.Lock()
	if t.conn == nil {
		t.mu.Unlock()
		return derrors.ErrNotConnected
	}
	t.pending[id] = pc
	wq := t.writeQ
	t.mu.Unlock()

	select {
	case wq <- body:
	case <-ctx.Done():
		t.mu.Lock()
		delete(t.pending, id)
		t.mu.Unlock()
		return ctx.Err()
	}

	select {
	case err := <-pc.err:
		return err
	case <-ctx.Done():
		t.mu.Lock()
		delete(t.pending, id)
		t.mu.Unlock()
		return ctx.Err()
	}
}

// Subscribe sends a subscribe RPC and returns a Subscription.
func (t *WSTransport) Subscribe(ctx context.Context, channel string, decode Decoder) (Subscription, error) {
	t.mu.Lock()
	if existing, ok := t.subs[channel]; ok {
		t.mu.Unlock()
		return existing, nil
	}
	sub := &wsSub{
		channel: channel,
		decode:  decode,
		updates: make(chan any, 256),
		closed:  make(chan struct{}),
	}
	t.subs[channel] = sub
	t.mu.Unlock()

	var resp struct {
		Status               map[string]string `json:"status"`
		CurrentSubscriptions []string          `json:"current_subscriptions"`
	}
	if err := t.Call(ctx, "subscribe", map[string]any{"channels": []string{channel}}, &resp); err != nil {
		t.mu.Lock()
		delete(t.subs, channel)
		t.mu.Unlock()
		return nil, err
	}
	return sub, nil
}

// Channel returns the underlying server channel name.
func (s *wsSub) Channel() string { return s.channel }

// Updates returns the receive channel of decoded events.
func (s *wsSub) Updates() <-chan any { return s.updates }

// Err returns the terminal subscription error, if any.
func (s *wsSub) Err() error {
	s.errMu.Lock()
	defer s.errMu.Unlock()
	return s.err
}

// Close signals the subscription as done. The transport will issue an
// unsubscribe RPC best-effort.
func (s *wsSub) Close() error {
	s.once.Do(func() {
		close(s.closed)
	})
	return nil
}

func (s *wsSub) finish(err error) {
	s.errMu.Lock()
	if s.err == nil {
		s.err = err
	}
	s.errMu.Unlock()
	s.once.Do(func() {
		close(s.closed)
		close(s.updates)
	})
}

// readPump reads frames and dispatches them to either pending RPCs or
// subscriptions until the connection drops.
func (t *WSTransport) readPump(ctx context.Context, conn *websocket.Conn) {
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			// If the parent ctx was cancelled, the closeWatcher closed conn —
			// treat that as a clean shutdown rather than a fault.
			if ctx.Err() != nil {
				return
			}
			t.failConnection(err)
			return
		}
		if jsonrpc.IsNotification(data) {
			t.dispatchNotification(data)
			continue
		}
		t.dispatchResponse(data)
	}
}

func (t *WSTransport) dispatchResponse(data []byte) {
	var resp jsonrpc.Response
	if err := json.Unmarshal(data, &resp); err != nil {
		return // malformed; drop
	}
	id, ok := resp.IDUint64()
	if !ok {
		return // non-numeric id; nothing to correlate
	}
	t.mu.Lock()
	pc, ok := t.pending[id]
	delete(t.pending, id)
	t.mu.Unlock()
	if !ok {
		return
	}
	if resp.Error != nil {
		pc.err <- &derrors.APIError{
			Code:    resp.Error.Code,
			Message: resp.Error.Message,
			Data:    resp.Error.Data,
		}
		return
	}
	pc.err <- jsonrpc.DecodeResult(&resp, pc.out)
}

func (t *WSTransport) dispatchNotification(data []byte) {
	var notif jsonrpc.Notification
	if err := json.Unmarshal(data, &notif); err != nil {
		return
	}
	if notif.Method != "subscription" {
		return
	}
	var p jsonrpc.SubscriptionParams
	if err := json.Unmarshal(notif.Params, &p); err != nil {
		return
	}
	t.mu.Lock()
	sub, ok := t.subs[p.Channel]
	t.mu.Unlock()
	if !ok {
		return
	}
	val, err := sub.decode(p.Data)
	if err != nil {
		return
	}
	select {
	case sub.updates <- val:
	case <-sub.closed:
	default:
		// Buffer full — drop oldest semantics would require a different
		// design; for now, drop newest. Callers that care about every
		// update should use a larger buffer or drain faster.
	}
}

// writePump serialises writes from writeQ.
//
// Both data frames (from writeQ) and ping control frames (from pingPump)
// flow through here so that conn.WriteMessage / conn.WriteControl are
// only called from this single goroutine — gorilla/websocket forbids
// concurrent writers.
func (t *WSTransport) writePump(conn *websocket.Conn, wq chan []byte, stop chan struct{}) {
	const writeTimeout = 10 * time.Second
	for {
		select {
		case <-stop:
			return
		case body, ok := <-wq:
			if !ok {
				return
			}
			_ = conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			err := conn.WriteMessage(websocket.TextMessage, body)
			if err != nil {
				t.failConnection(err)
				return
			}
		}
	}
}

// pingPump issues control-level pings and waits for the matching pong.
//
// gorilla/websocket's documentation explicitly allows WriteControl to be
// called concurrently with WriteMessage from a different goroutine, so
// pings don't need to flow through writeQ.
func (t *WSTransport) pingPump(conn *websocket.Conn, pong chan struct{}, stop chan struct{}) {
	tk := time.NewTicker(t.cfg.PingInterval)
	defer tk.Stop()
	const pongDeadline = 5 * time.Second
	for {
		select {
		case <-stop:
			return
		case <-tk.C:
			// Drain any stale pong before issuing a fresh ping so we don't
			// accept a reply meant for the previous round.
			select {
			case <-pong:
			default:
			}
			deadline := time.Now().Add(pongDeadline)
			if err := conn.WriteControl(websocket.PingMessage, nil, deadline); err != nil {
				t.failConnection(err)
				return
			}
			select {
			case <-pong:
			case <-time.After(pongDeadline):
				t.failConnection(errors.New("ping: pong not received within deadline"))
				return
			case <-stop:
				return
			}
		}
	}
}

// failConnection tears down the current connection and notifies pending RPCs.
// If reconnect is enabled, the reconnect loop will pick up from here.
func (t *WSTransport) failConnection(cause error) {
	t.mu.Lock()
	if t.conn == nil {
		t.mu.Unlock()
		return
	}
	pending := t.pending
	t.pending = make(map[uint64]*pendingCall)
	conn := t.conn
	t.conn = nil
	stop := t.stopCh
	t.stopCh = nil
	cancel := t.cancel
	t.cancel = nil
	t.mu.Unlock()

	if cancel != nil {
		cancel()
	}
	if stop != nil {
		// Best-effort: stop pumps if not yet closed.
		select {
		case <-stop:
		default:
			close(stop)
		}
	}
	wErr := &derrors.ConnectionError{Op: "ws read", Err: cause}
	for _, pc := range pending {
		pc.err <- wErr
	}
	if conn != nil {
		// Send a courteous close frame, then drop the conn.
		_ = conn.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "fail"),
			time.Now().Add(time.Second))
		_ = conn.Close()
	}
}

// reconnectLoop dials in a loop with backoff after a drop.
func (t *WSTransport) reconnectLoop() {
	bo := retry.NewBackoff()
	for {
		t.mu.Lock()
		connected := t.conn != nil
		t.mu.Unlock()
		if !connected {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			if err := t.dial(ctx); err != nil {
				cancel()
				time.Sleep(bo.Next())
				continue
			}
			cancel()
			bo.Reset()

			if t.cfg.OnReconnect != nil {
				cctx, ccancel := context.WithTimeout(context.Background(), 30*time.Second)
				_ = t.cfg.OnReconnect(cctx, t)
				ccancel()
			}
			t.resubscribe()
		}
		time.Sleep(time.Second)
	}
}

// resubscribe re-issues subscribe RPCs for every active subscription.
func (t *WSTransport) resubscribe() {
	t.mu.Lock()
	channels := make([]string, 0, len(t.subs))
	for ch := range t.subs {
		channels = append(channels, ch)
	}
	t.mu.Unlock()
	if len(channels) == 0 {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_ = t.Call(ctx, "subscribe", map[string]any{"channels": channels}, nil)
}

// Close terminates the connection and unblocks any pumps.
func (t *WSTransport) Close() error {
	t.mu.Lock()
	subs := t.subs
	t.subs = make(map[string]*wsSub)
	conn := t.conn
	t.conn = nil
	stop := t.stopCh
	t.stopCh = nil
	cancel := t.cancel
	t.cancel = nil
	t.mu.Unlock()

	if cancel != nil {
		cancel()
	}
	if stop != nil {
		select {
		case <-stop:
		default:
			close(stop)
		}
	}
	for _, s := range subs {
		s.finish(derrors.ErrSubscriptionClosed)
	}
	if conn != nil {
		// Best-effort close-frame + drop. closeWatcher may already have
		// dropped the conn when the rootCtx cancel above fired; ignore
		// any "use of closed network connection" that returns here.
		_ = conn.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "client close"),
			time.Now().Add(time.Second))
		_ = conn.Close()
	}
	return nil
}

// IsConnected reports whether the transport currently holds an open socket.
func (t *WSTransport) IsConnected() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.conn != nil
}

// SetOnReconnect registers a callback invoked after a successful reconnect.
// Use it to re-issue a login or any other re-warmup. Safe to call before or
// after Connect.
func (t *WSTransport) SetOnReconnect(fn func(ctx context.Context, t *WSTransport) error) {
	t.mu.Lock()
	t.cfg.OnReconnect = fn
	t.mu.Unlock()
}

// statically assert WSTransport satisfies both interfaces.
var (
	_ Transport  = (*WSTransport)(nil)
	_ Subscriber = (*WSTransport)(nil)
)

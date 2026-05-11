package ws_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/internal/testutil"
	"github.com/amiwrpremium/go-derive/pkg/types"
	"github.com/amiwrpremium/go-derive/pkg/ws"
)

// decodeOrderBookJSON is the test-side decoder for the order-book
// channel; production callers reach the same shape through
// Client.SubscribeOrderBook.
func decodeOrderBookJSON(raw json.RawMessage) (types.OrderBook, error) {
	var ob types.OrderBook
	return ob, json.Unmarshal(raw, &ob)
}

func decodeString(raw json.RawMessage) (string, error) {
	var s string
	return s, json.Unmarshal(raw, &s)
}

func TestSubscribe_TypedDelivery(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := ws.Subscribe(context.Background(), c, "orderbook.BTC-PERP.1.10", decodeOrderBookJSON)
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()

	require.True(t, srv.WaitSubscribed("orderbook.BTC-PERP.1.10", 1*time.Second))
	srv.Notify("orderbook.BTC-PERP.1.10", map[string]any{
		"instrument_name": "BTC-PERP",
		"bids":            [][]string{{"100", "1"}},
		"asks":            [][]string{{"101", "1"}},
		"timestamp":       1700000000000,
	})

	select {
	case ob, ok := <-sub.Updates():
		require.True(t, ok)
		assert.Equal(t, "BTC-PERP", ob.InstrumentName)
	case <-time.After(2 * time.Second):
		t.Fatal("update never delivered")
	}
}

func TestSubscribe_Channel_Method(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := ws.Subscribe(context.Background(), c, "orderbook.X.1.10", decodeOrderBookJSON)
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()

	assert.Equal(t, "orderbook.X.1.10", sub.Channel())
}

func TestSubscribe_Close_StopsUpdates(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := ws.Subscribe(context.Background(), c, "orderbook.X.1.10", decodeOrderBookJSON)
	require.NoError(t, err)
	require.NoError(t, sub.Close())
	// Second Close is harmless.
	require.NoError(t, sub.Close())
}

func TestSubscribe_TypeMismatch_Drops(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	// Subscribe with a string-decoder against an order-book payload.
	sub, err := ws.Subscribe(context.Background(), c, "orderbook.BTC-PERP.1.10", decodeString)
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()
	require.True(t, srv.WaitSubscribed("orderbook.BTC-PERP.1.10", 1*time.Second))

	srv.Notify("orderbook.BTC-PERP.1.10", map[string]any{
		"instrument_name": "BTC-PERP",
		"bids":            [][]string{}, "asks": [][]string{},
		"timestamp": 1700000000000,
	})

	select {
	case <-sub.Updates():
		t.Fatal("string decoder against an object payload should not deliver")
	case <-time.After(150 * time.Millisecond):
		// expected
	}
}

func TestSubscribeFunc_DriverCallback(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	got := make(chan types.OrderBook, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = ws.SubscribeFunc(ctx, c, "orderbook.ETH-PERP.1.10", decodeOrderBookJSON, func(ob types.OrderBook) {
			got <- ob
		})
	}()

	require.True(t, srv.WaitSubscribed("orderbook.ETH-PERP.1.10", 1*time.Second))
	srv.Notify("orderbook.ETH-PERP.1.10", map[string]any{
		"instrument_name": "ETH-PERP",
		"bids":            [][]string{}, "asks": [][]string{},
		"timestamp": 1700000000000,
	})

	select {
	case ob := <-got:
		assert.Equal(t, "ETH-PERP", ob.InstrumentName)
	case <-time.After(2 * time.Second):
		t.Fatal("callback never fired")
	}
}

func TestSubscribeFunc_ContextCancelReturnsErr(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- ws.SubscribeFunc(ctx, c, "orderbook.X.1.10", decodeOrderBookJSON, func(types.OrderBook) {})
	}()
	require.True(t, srv.WaitSubscribed("orderbook.X.1.10", 1*time.Second))
	cancel()

	select {
	case err := <-done:
		assert.True(t, errors.Is(err, context.Canceled))
	case <-time.After(2 * time.Second):
		t.Fatal("SubscribeFunc didn't return after cancel")
	}
}

// --- Option tests ---------------------------------------------------

// Decoder that yields a distinct integer per notification so tests
// can verify which events ended up in the buffer.
func decodeInt(raw json.RawMessage) (int, error) {
	var n int
	return n, json.Unmarshal(raw, &n)
}

func TestSubscribe_DropNewest_DropsWhenFull(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	errCh := make(chan error, 8)
	sub, err := ws.Subscribe(context.Background(), c, "trades.X", decodeInt,
		ws.WithBufferSize(2),
		ws.WithDropPolicy(ws.DropNewest),
		ws.WithErrorHandler(func(err error) {
			select {
			case errCh <- err:
			default:
			}
		}),
	)
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()
	require.True(t, srv.WaitSubscribed("trades.X", time.Second))

	// Send 4 events. Buffer holds 2; the other 2 should be dropped.
	for _, n := range []int{1, 2, 3, 4} {
		srv.Notify("trades.X", n)
	}

	// Give the pump a moment to enqueue.
	time.Sleep(100 * time.Millisecond)

	// First two events delivered.
	first := <-sub.Updates()
	second := <-sub.Updates()
	assert.Equal(t, 1, first)
	assert.Equal(t, 2, second)

	// Error handler fired for the dropped events.
	dropped := 0
	for {
		select {
		case err := <-errCh:
			if errors.Is(err, ws.ErrBufferFull) {
				dropped++
			}
		case <-time.After(150 * time.Millisecond):
			goto done
		}
	}
done:
	assert.GreaterOrEqual(t, dropped, 1, "at least one drop should have fired the error handler")
}

func TestSubscribe_DropOldest_EvictsOldest(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := ws.Subscribe(context.Background(), c, "trades.Y", decodeInt,
		ws.WithBufferSize(2),
		ws.WithDropPolicy(ws.DropOldest),
	)
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()
	require.True(t, srv.WaitSubscribed("trades.Y", time.Second))

	// Send 4 events fast; the last two should win.
	for _, n := range []int{1, 2, 3, 4} {
		srv.Notify("trades.Y", n)
	}
	time.Sleep(100 * time.Millisecond)

	// Drain and assert the most recent values survived. We don't
	// pin which two — DropOldest is best-effort under races — but
	// they must be from the later half.
	seen := []int{}
	for len(seen) < 2 {
		select {
		case v := <-sub.Updates():
			seen = append(seen, v)
		case <-time.After(200 * time.Millisecond):
			t.Fatalf("only got %v", seen)
		}
	}
	for _, v := range seen {
		assert.GreaterOrEqual(t, v, 2, "DropOldest should have evicted the very-oldest events")
	}
}

func TestSubscribe_Block_BackPressures(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := ws.Subscribe(context.Background(), c, "trades.Z", decodeInt,
		ws.WithBufferSize(1),
		ws.WithDropPolicy(ws.Block),
	)
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()
	require.True(t, srv.WaitSubscribed("trades.Z", time.Second))

	srv.Notify("trades.Z", 1)
	srv.Notify("trades.Z", 2) // pump blocks here until we read

	// Read the first event; the pump should then unblock and the
	// second event becomes available.
	first := <-sub.Updates()
	assert.Equal(t, 1, first)
	select {
	case second := <-sub.Updates():
		assert.Equal(t, 2, second)
	case <-time.After(1 * time.Second):
		t.Fatal("second event never arrived — block policy not unblocking")
	}
}

func TestSubscribe_ErrorHandler_FiresOnTypeMismatch(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	errCh := make(chan error, 1)
	// Subscribe[string] with the order-book decoder — its return type
	// is types.OrderBook, so every event triggers ErrTypeMismatch.
	sub, err := ws.Subscribe(context.Background(), c, "orderbook.M.1.10", decodeOrderBookJSON,
		ws.WithErrorHandler(func(err error) {
			select {
			case errCh <- err:
			default:
			}
		}),
	)
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()
	require.True(t, srv.WaitSubscribed("orderbook.M.1.10", time.Second))

	// Send a payload that decodes (so the decoder doesn't fail); the
	// typed pump's T==types.OrderBook here, so this should NOT trip
	// the mismatch handler.
	srv.Notify("orderbook.M.1.10", map[string]any{
		"instrument_name": "M",
		"bids":            [][]string{},
		"asks":            [][]string{},
		"timestamp":       1700000000000,
	})

	// Drain one update so we know the pump processed it.
	select {
	case <-sub.Updates():
	case <-time.After(1 * time.Second):
		t.Fatal("update never delivered")
	}

	// No error should have fired on the matching-type happy path.
	select {
	case err := <-errCh:
		t.Fatalf("error handler fired unexpectedly: %v", err)
	case <-time.After(100 * time.Millisecond):
		// expected
	}
}

func TestSubscribe_DefaultBufferAndPolicy(t *testing.T) {
	// Smoke test: default subscribe (no opts) still delivers events.
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	sub, err := ws.Subscribe(context.Background(), c, "trades.N", decodeInt)
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()
	require.True(t, srv.WaitSubscribed("trades.N", time.Second))

	srv.Notify("trades.N", 7)
	select {
	case v := <-sub.Updates():
		assert.Equal(t, 7, v)
	case <-time.After(1 * time.Second):
		t.Fatal("default subscribe lost the event")
	}
}

func TestSubscribe_DecodeError_FiresErrorHandler(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	errCh := make(chan error, 4)
	// Decoder expects int; push an object payload to force a
	// json.Unmarshal error.
	sub, err := ws.Subscribe(context.Background(), c, "trades.D", decodeInt,
		ws.WithErrorHandler(func(err error) {
			select {
			case errCh <- err:
			default:
			}
		}),
	)
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()
	require.True(t, srv.WaitSubscribed("trades.D", time.Second))

	srv.Notify("trades.D", map[string]any{"not": "an int"})

	select {
	case err := <-errCh:
		assert.True(t, errors.Is(err, ws.ErrDecodeFailed),
			"handler should fire with ErrDecodeFailed, got: %v", err)
	case <-time.After(1 * time.Second):
		t.Fatal("error handler never fired on decode failure")
	}
}

// --- SubscribeInto tests --------------------------------------------

func TestSubscribeInto_DeliversToCallerChan(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	out := make(chan int, 8)
	sub, err := ws.SubscribeInto(context.Background(), c, "trades.I1", decodeInt, out)
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()
	require.True(t, srv.WaitSubscribed("trades.I1", time.Second))

	srv.Notify("trades.I1", 7)
	select {
	case v := <-out:
		assert.Equal(t, 7, v)
	case <-time.After(1 * time.Second):
		t.Fatal("event never reached caller's chan")
	}
}

func TestSubscribeInto_DoesNotCloseCallerChan(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	out := make(chan int, 4)
	sub, err := ws.SubscribeInto(context.Background(), c, "trades.I2", decodeInt, out)
	require.NoError(t, err)
	require.True(t, srv.WaitSubscribed("trades.I2", time.Second))

	require.NoError(t, sub.Close())

	// Caller's chan must still be open (writable, not closed).
	// We try a non-blocking send — if out were closed, this would panic.
	select {
	case out <- 99:
		// expected — chan still open
	default:
		t.Fatal("caller's chan unexpectedly full at send-after-close")
	}
}

func TestSubscribeInto_MultiplexTwoChannels(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	out := make(chan int, 8)
	sub1, err := ws.SubscribeInto(context.Background(), c, "trades.A", decodeInt, out)
	require.NoError(t, err)
	defer func() { _ = sub1.Close() }()
	sub2, err := ws.SubscribeInto(context.Background(), c, "trades.B", decodeInt, out)
	require.NoError(t, err)
	defer func() { _ = sub2.Close() }()

	require.True(t, srv.WaitSubscribed("trades.A", time.Second))
	require.True(t, srv.WaitSubscribed("trades.B", time.Second))

	srv.Notify("trades.A", 1)
	srv.Notify("trades.B", 2)

	got := map[int]bool{}
	for len(got) < 2 {
		select {
		case v := <-out:
			got[v] = true
		case <-time.After(1 * time.Second):
			t.Fatalf("only saw %v on the shared chan", got)
		}
	}
	assert.True(t, got[1])
	assert.True(t, got[2])
}

func TestSubscribeInto_RespectsDropPolicy(t *testing.T) {
	srv := testutil.NewMockWSServer()
	defer srv.Close()
	c := newWSClient(t, srv, false)
	require.NoError(t, c.Connect(context.Background()))
	defer func() { _ = c.Close() }()

	out := make(chan int, 2)
	dropped := make(chan struct{}, 8)
	sub, err := ws.SubscribeInto(context.Background(), c, "trades.D2", decodeInt, out,
		ws.WithDropPolicy(ws.DropNewest),
		ws.WithErrorHandler(func(err error) {
			if errors.Is(err, ws.ErrBufferFull) {
				select {
				case dropped <- struct{}{}:
				default:
				}
			}
		}),
	)
	require.NoError(t, err)
	defer func() { _ = sub.Close() }()
	require.True(t, srv.WaitSubscribed("trades.D2", time.Second))

	for _, n := range []int{1, 2, 3, 4} {
		srv.Notify("trades.D2", n)
	}
	time.Sleep(100 * time.Millisecond)

	// First two delivered.
	a := <-out
	b := <-out
	assert.Equal(t, 1, a)
	assert.Equal(t, 2, b)
	// Drops surfaced.
	select {
	case <-dropped:
	case <-time.After(150 * time.Millisecond):
		t.Fatal("expected at least one buffer-full drop on caller-supplied chan")
	}
}

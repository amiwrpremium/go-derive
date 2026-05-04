package errors_test

import (
	stderrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"

	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
)

func TestConnectionError_Error(t *testing.T) {
	e := &derrors.ConnectionError{Op: "dial", Err: stderrors.New("refused")}
	got := e.Error()
	assert.Contains(t, got, "dial")
	assert.Contains(t, got, "refused")
}

func TestConnectionError_Unwrap(t *testing.T) {
	inner := stderrors.New("dial timeout")
	e := &derrors.ConnectionError{Op: "dial", Err: inner}
	assert.True(t, stderrors.Is(e, inner))
	assert.Same(t, inner, stderrors.Unwrap(e))
}

func TestConnectionError_NilInner(t *testing.T) {
	// Op alone is enough — Unwrap returns nil.
	e := &derrors.ConnectionError{Op: "dial"}
	assert.Contains(t, e.Error(), "dial")
	assert.Nil(t, stderrors.Unwrap(e))
}

func TestTimeoutError_Error(t *testing.T) {
	e := &derrors.TimeoutError{Method: "private/order"}
	got := e.Error()
	assert.Contains(t, got, "private/order")
	assert.Contains(t, got, "timeout")
}

func TestTimeoutError_Implements_Error(_ *testing.T) {
	var _ error = &derrors.TimeoutError{}
}

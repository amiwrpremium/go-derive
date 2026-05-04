package errors_test

import (
	stderrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"

	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
)

func TestSigningError_Error(t *testing.T) {
	e := &derrors.SigningError{Op: "parse", Err: stderrors.New("bad key")}
	got := e.Error()
	assert.Contains(t, got, "parse")
	assert.Contains(t, got, "bad key")
}

func TestSigningError_Unwrap(t *testing.T) {
	inner := stderrors.New("bad key")
	e := &derrors.SigningError{Op: "parse", Err: inner}
	assert.True(t, stderrors.Is(e, inner))
	assert.Same(t, inner, stderrors.Unwrap(e))
}

func TestSigningError_NilInner(t *testing.T) {
	e := &derrors.SigningError{Op: "noop"}
	assert.Contains(t, e.Error(), "noop")
	assert.Nil(t, stderrors.Unwrap(e))
}

func TestExpiredSignatureError_Error(t *testing.T) {
	e := &derrors.ExpiredSignatureError{ExpiryUnixSec: 100, NowUnixSec: 200}
	got := e.Error()
	assert.Contains(t, got, "100")
	assert.Contains(t, got, "200")
	assert.Contains(t, got, "expired")
}

func TestExpiredSignatureError_FutureExpiry(t *testing.T) {
	// Even when ExpiryUnixSec > NowUnixSec, the error format is unchanged
	// (the type is a plain data carrier; semantics are the caller's).
	e := &derrors.ExpiredSignatureError{ExpiryUnixSec: 1000, NowUnixSec: 500}
	assert.Contains(t, e.Error(), "1000")
}

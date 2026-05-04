package auth_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/auth"
)

func TestHTTPHeaders_NilSignerYieldsNoHeaders(t *testing.T) {
	h, err := auth.HTTPHeaders(context.Background(), nil, time.Now())
	require.NoError(t, err)
	assert.Nil(t, h)
}

func TestHTTPHeaders_PopulatesAllThreeFields(t *testing.T) {
	s, err := auth.NewLocalSigner(testKey)
	require.NoError(t, err)
	now := time.Unix(1700000000, 123_000_000) // millisecond-precise
	h, err := auth.HTTPHeaders(context.Background(), s, now)
	require.NoError(t, err)
	assert.Equal(t, s.Owner().Hex(), h.Get("X-LyraWallet"))
	assert.Equal(t, "1700000000123", h.Get("X-LyraTimestamp"))
	assert.True(t, strings.HasPrefix(h.Get("X-LyraSignature"), "0x"))
	assert.Len(t, h.Get("X-LyraSignature"), 2+65*2)
}

// failSigner forces the SignAuthHeader path to error so HTTPHeaders'
// error-propagation branch is exercised.
type failSigner struct{ auth.Signer }

func (failSigner) SignAuthHeader(context.Context, time.Time) (auth.Signature, error) {
	return auth.Signature{}, errBoom
}

var errBoom = newErr("boom")

type sentinelErr struct{ msg string }

func (e sentinelErr) Error() string { return e.msg }

func newErr(s string) sentinelErr { return sentinelErr{msg: s} }

func TestHTTPHeaders_PropagatesSignerError(t *testing.T) {
	// Wrap a real signer so the embedded interface's Address/Owner work,
	// then the test override returns an error from SignAuthHeader.
	real, err := auth.NewLocalSigner(testKey)
	require.NoError(t, err)
	s := failSigner{Signer: real}
	_, err = auth.HTTPHeaders(context.Background(), s, time.Now())
	assert.ErrorContains(t, err, "boom")
}

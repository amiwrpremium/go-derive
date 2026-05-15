package auth_test

// EIP-191 (personal_sign) is exercised via SignAuthHeader. Cover:
//   - Same timestamp produces the same signature (deterministic over key+digest).
//   - Different timestamps produce different signatures (length prefix changes).
//   - Recovered address matches signer.

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/auth"
)

func TestEIP191_DeterministicForSameTimestamp(t *testing.T) {
	s, err := auth.NewLocalSigner(testKey)
	require.NoError(t, err)
	ts := time.Unix(1700000000, 0)

	sig1, err := s.SignAuthHeader(context.Background(), ts)
	require.NoError(t, err)
	sig2, err := s.SignAuthHeader(context.Background(), ts)
	require.NoError(t, err)
	assert.Equal(t, sig1, sig2)
}

func TestEIP191_DifferentTimestampsDifferSignature(t *testing.T) {
	s, err := auth.NewLocalSigner(testKey)
	require.NoError(t, err)
	a, err := s.SignAuthHeader(context.Background(), time.Unix(1700000000, 0))
	require.NoError(t, err)
	b, err := s.SignAuthHeader(context.Background(), time.Unix(1700000001, 0))
	require.NoError(t, err)
	assert.NotEqual(t, a, b)
}

func TestEIP191_RecoverableViaPersonalSignDigest(t *testing.T) {
	s, err := auth.NewLocalSigner(testKey)
	require.NoError(t, err)
	ts := time.Unix(1700000000, 0)
	sig, err := s.SignAuthHeader(context.Background(), ts)
	require.NoError(t, err)

	msg := []byte(strconv.FormatInt(ts.UnixMilli(), 10))
	digest := personalHash(msg)
	pub, err := crypto.SigToPub(digest, normaliseV(sig[:]))
	require.NoError(t, err)
	assert.Equal(t, s.SessionAddress(), crypto.PubkeyToAddress(*pub))
}

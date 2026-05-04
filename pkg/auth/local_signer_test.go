package auth_test

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/internal/netconf"
	"github.com/amiwrpremium/go-derive/pkg/auth"
)

func TestNewLocalSigner_Hex0xPrefixed(t *testing.T) {
	s, err := auth.NewLocalSigner(testKey)
	require.NoError(t, err)
	require.NotNil(t, s)
}

func TestNewLocalSigner_HexNoPrefix(t *testing.T) {
	s, err := auth.NewLocalSigner(strings.TrimPrefix(testKey, "0x"))
	require.NoError(t, err)
	require.NotNil(t, s)
}

func TestNewLocalSigner_RejectsBadHex(t *testing.T) {
	for _, in := range []string{"", "not-hex", "0xZZ", "0x12"} {
		t.Run(in, func(t *testing.T) {
			_, err := auth.NewLocalSigner(in)
			assert.Error(t, err)
		})
	}
}

func TestLocalSigner_AddressMatchesPublicKey(t *testing.T) {
	s, err := auth.NewLocalSigner(testKey)
	require.NoError(t, err)
	k, err := crypto.HexToECDSA(strings.TrimPrefix(testKey, "0x"))
	require.NoError(t, err)
	assert.Equal(t, crypto.PubkeyToAddress(k.PublicKey), s.Address())
}

func TestLocalSigner_OwnerEqualsAddress(t *testing.T) {
	s, err := auth.NewLocalSigner(testKey)
	require.NoError(t, err)
	assert.Equal(t, s.Address(), s.Owner(), "LocalSigner has no separate owner")
}

func TestLocalSigner_SignAuthHeader_Recover(t *testing.T) {
	s, err := auth.NewLocalSigner(testKey)
	require.NoError(t, err)
	ts := time.Now()

	sig, err := s.SignAuthHeader(context.Background(), ts)
	require.NoError(t, err)

	msg := []byte(strconv.FormatInt(ts.UnixMilli(), 10))
	digest := personalHash(msg)
	pub, err := crypto.SigToPub(digest, normaliseV(sig[:]))
	require.NoError(t, err)
	assert.Equal(t, s.Address(), crypto.PubkeyToAddress(*pub))
}

func TestLocalSigner_SignAction_Determinism(t *testing.T) {
	s, err := auth.NewLocalSigner(testKey)
	require.NoError(t, err)
	domain := netconf.Mainnet().EIP712Domain()
	action := auth.ActionData{SubaccountID: 1, Nonce: 12345, Expiry: 1700000000}
	a, err := s.SignAction(context.Background(), domain, action)
	require.NoError(t, err)
	b, err := s.SignAction(context.Background(), domain, action)
	require.NoError(t, err)
	assert.Equal(t, a, b)
}

func TestLocalSigner_SignAction_PopulatesOwnerAndSignerWhenZero(t *testing.T) {
	// When Owner/Signer fields on the action are unset, the signer fills
	// them in with its own address; mutating one of those fields changes
	// the signature.
	s, err := auth.NewLocalSigner(testKey)
	require.NoError(t, err)
	domain := netconf.Mainnet().EIP712Domain()

	a, err := s.SignAction(context.Background(), domain, auth.ActionData{Nonce: 1})
	require.NoError(t, err)
	b, err := s.SignAction(context.Background(), domain, auth.ActionData{Nonce: 2})
	require.NoError(t, err)
	assert.NotEqual(t, a, b)
}

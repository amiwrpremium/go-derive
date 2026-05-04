package auth_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/internal/netconf"
	"github.com/amiwrpremium/go-derive/pkg/auth"
)

func TestNewSessionKeySigner_Happy(t *testing.T) {
	owner := common.HexToAddress("0x1111111111111111111111111111111111111111")
	s, err := auth.NewSessionKeySigner(testKey, owner)
	require.NoError(t, err)
	assert.NotNil(t, s)
}

func TestNewSessionKeySigner_RejectsBadKey(t *testing.T) {
	_, err := auth.NewSessionKeySigner("not-hex", common.Address{})
	assert.Error(t, err)
}

func TestSessionKeySigner_OwnerSeparateFromAddress(t *testing.T) {
	owner := common.HexToAddress("0x1111111111111111111111111111111111111111")
	s, err := auth.NewSessionKeySigner(testKey, owner)
	require.NoError(t, err)
	assert.NotEqual(t, s.Address(), s.Owner())
	assert.Equal(t, owner, s.Owner())
}

func TestSessionKeySigner_SignAuthHeaderDelegates(t *testing.T) {
	owner := common.HexToAddress("0x1111111111111111111111111111111111111111")
	s, err := auth.NewSessionKeySigner(testKey, owner)
	require.NoError(t, err)
	// Should sign with the session key (Address), not the owner.
	sig, err := s.SignAuthHeader(context.Background(), timeNowDeterministic())
	require.NoError(t, err)
	assert.NotEqual(t, [65]byte{}, sig)
}

func TestSessionKeySigner_SignActionStampsOwnerAndSigner(t *testing.T) {
	// SignAction should override Owner with the configured owner address
	// and Signer with the session-key address. Verifying that means signing
	// twice with different external Owner/Signer values still yields the
	// same signature.
	owner := common.HexToAddress("0x2222222222222222222222222222222222222222")
	s, err := auth.NewSessionKeySigner(testKey, owner)
	require.NoError(t, err)
	domain := netconf.Mainnet().EIP712Domain()
	a := auth.ActionData{Nonce: 1, Owner: common.Address{}, Signer: common.Address{}}
	b := auth.ActionData{Nonce: 1, Owner: common.HexToAddress("0xdeadbeef00000000000000000000000000000000"), Signer: common.HexToAddress("0xfeedface00000000000000000000000000000000")}

	sigA, err := s.SignAction(context.Background(), domain, a)
	require.NoError(t, err)
	sigB, err := s.SignAction(context.Background(), domain, b)
	require.NoError(t, err)
	assert.Equal(t, sigA, sigB, "SignAction must overwrite Owner/Signer fields with the configured ones")
}

// timeNowDeterministic is a tiny indirection so changing this in one place
// propagates to every test that relies on a fixed timestamp.
func timeNowDeterministic() (t timeT) { return timeT{} }

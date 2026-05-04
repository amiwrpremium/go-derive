package auth_test

// EIP-712 internals are unexported, so we can only test them transitively
// through SignAction. Two assertions:
//
//   1. domainSeparator depends on ChainID — same action on mainnet vs testnet
//      produces different signatures.
//   2. hashTypedData applies the \x19\x01 envelope, which means changing the
//      verifying-contract address changes the signature too.

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/internal/netconf"
	"github.com/amiwrpremium/go-derive/pkg/auth"
)

func TestDomain_ChainIDAffectsSignature(t *testing.T) {
	s, err := auth.NewLocalSigner(testKey)
	require.NoError(t, err)
	a := auth.ActionData{Nonce: 1, Expiry: 1}

	mainnetSig, err := s.SignAction(context.Background(), netconf.Mainnet().EIP712Domain(), a)
	require.NoError(t, err)
	testnetSig, err := s.SignAction(context.Background(), netconf.Testnet().EIP712Domain(), a)
	require.NoError(t, err)

	assert.NotEqual(t, mainnetSig, testnetSig,
		"different chain IDs must produce different signatures")
}

func TestDomain_VerifyingContractAffectsSignature(t *testing.T) {
	s, err := auth.NewLocalSigner(testKey)
	require.NoError(t, err)
	a := auth.ActionData{Nonce: 1}

	cfg := netconf.Mainnet()
	d1 := cfg.EIP712Domain()
	d2 := d1
	d2.VerifyingContract = common.HexToAddress("0xdeadbeef00000000000000000000000000000000").Hex()

	sig1, err := s.SignAction(context.Background(), d1, a)
	require.NoError(t, err)
	sig2, err := s.SignAction(context.Background(), d2, a)
	require.NoError(t, err)
	assert.NotEqual(t, sig1, sig2)
}

func TestDomain_SameInputsSameSignature(t *testing.T) {
	// Determinism: same key + same action + same domain → same signature
	// (already covered in local_signer_test, repeated here as the EIP-712
	// envelope is the load-bearing piece).
	s, err := auth.NewLocalSigner(testKey)
	require.NoError(t, err)
	d := netconf.Mainnet().EIP712Domain()
	a := auth.ActionData{Nonce: 99, Expiry: 1700000000}
	x, err := s.SignAction(context.Background(), d, a)
	require.NoError(t, err)
	y, err := s.SignAction(context.Background(), d, a)
	require.NoError(t, err)
	assert.Equal(t, x, y)
}

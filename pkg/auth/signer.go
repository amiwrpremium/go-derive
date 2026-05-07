package auth

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/amiwrpremium/go-derive/internal/netconf"
)

// Signature is a 65-byte ECDSA signature in `r || s || v` byte order,
// where `v` follows Ethereum's 27/28 convention (not the raw 0/1 form
// go-ethereum produces internally — Derive's on-chain ecrecover path
// expects 27/28).
type Signature [65]byte

// Hex returns the canonical 0x-prefixed lowercase-hex representation.
// Length is always 132 characters (2 prefix + 65 bytes × 2).
func (s Signature) Hex() string {
	const hexChars = "0123456789abcdef"
	out := make([]byte, 2+len(s)*2)
	out[0] = '0'
	out[1] = 'x'
	for i, b := range s {
		out[2+i*2] = hexChars[b>>4]
		out[2+i*2+1] = hexChars[b&0x0f]
	}
	return string(out)
}

// Signer abstracts over the source of cryptographic signatures. The SDK
// uses it for both per-request auth-header signing (EIP-191) and
// per-action EIP-712 signing.
//
// Concrete implementations in this package:
//
//   - [LocalSigner]        — secp256k1 key held in process; owner == address.
//   - [SessionKeySigner]   — session key signs but reports a separate owner.
//
// External implementations are welcome: a hardware wallet, KMS-backed
// key, or HSM-backed key all fit cleanly behind this interface.
type Signer interface {
	// Address returns the public address whose signatures the
	// implementation produces. For session keys this is the session
	// key's address, not the owner's.
	Address() common.Address

	// Owner returns the owner (smart-account) address. For [LocalSigner]
	// this equals [Signer.Address]; for [SessionKeySigner] it is the
	// distinct registered owner.
	Owner() common.Address

	// SignAction produces an EIP-712 signature over the action struct
	// hash with Derive's per-network domain. The implementation is
	// responsible for filling Action.Owner and Action.Signer if they
	// are zero.
	SignAction(ctx context.Context, domain netconf.Domain, action ActionData) (Signature, error)

	// SignAuthHeader produces an EIP-191 personal-sign signature over
	// the millisecond-timestamp string. The result is used as the
	// X-LyraSignature header on REST and as the `signature` field on
	// the WS `public/login` RPC.
	SignAuthHeader(ctx context.Context, ts time.Time) (Signature, error)
}

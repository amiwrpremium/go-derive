// Package derive — implements the two cryptographic signing flows Derive's
// API requires.
//
// # Two flows, one Signer
//
// Every authenticated Derive request involves cryptography in one of two
// places:
//
//  1. Per-request authentication of the caller. Sent as REST headers
//     (X-LyraWallet, X-LyraTimestamp, X-LyraSignature) or as a one-shot
//     `public/login` RPC over WebSocket. The signature is an EIP-191
//     personal-sign over the millisecond timestamp.
//
//  2. Per-action authorisation of order placement, cancels, transfers and
//     RFQ flows. The signature is an EIP-712 typed-data hash over an
//     `Action` struct whose `data` field is the keccak256 of an
//     ABI-encoded module-specific payload.
//
// Both flows go through the same [Signer] interface; concrete
// implementations include [LocalSigner] (owner key in process) and
// [SessionKeySigner] (session key delegating from a separate owner
// address).
//
// # Production setup
//
// Derive deployments use session keys. The owner is a smart-account on
// Derive Chain; the session key is a hot key registered on-chain as
// authorised to sign on its behalf. Use [NewSessionKeySigner] for
// production trading so the long-lived owner key never sits in the
// trading process's memory.
//
// # Test fixtures
//
// All signing test vectors live in pkg/auth/*_test.go. The tests verify
// that Derive's expected signature bytes can be reproduced from a known
// secp256k1 key — they're the canary for any future change to EIP-712
// hashing here or upstream in go-ethereum.
package derive

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/amiwrpremium/go-derive/internal/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// actionTypeHash is keccak256 of the canonical EIP-712 type string for
// Derive's `Action` struct:
//
//	Action(uint256 subaccountId, uint256 nonce, address module,
//	       bytes32 data, uint256 expiry, address owner, address signer)
var actionTypeHash = keccak([]byte(
	"Action(uint256 subaccountId,uint256 nonce,address module,bytes32 data,uint256 expiry,address owner,address signer)",
))

// ActionData is the input to Derive's order/cancel/transfer signing flow.
//
// It mirrors Solidity's `Action` struct field-for-field. The Data field is
// the keccak256 of the ABI-encoded module-specific payload — for trades
// that's [TradeModuleData.Hash], for transfers it's
// [TransferModuleData.Hash], and so on.
//
// Use [ActionData.Hash] to compute the EIP-712 struct hash; in normal use
// [Signer.SignAction] does that for you and returns the [Signature].
type ActionData struct {
	// SubaccountID is the placing subaccount id.
	SubaccountID int64
	// Nonce is a strictly-increasing per-subaccount nonce.
	// Use [NonceGen] to source these.
	Nonce uint64
	// Module is the on-chain Derive module contract this action targets
	// (e.g. the TradeModule for orders, TransferModule for transfers).
	Module common.Address
	// Data is keccak256 of the module-specific ABI-encoded payload.
	Data [32]byte
	// Expiry is the Unix timestamp (seconds) after which the signature is
	// no longer valid.
	Expiry int64
	// Owner is the smart-account owner address.
	Owner common.Address
	// Signer is the session-key (or owner) address that signed.
	Signer common.Address
}

// Hash returns the EIP-712 hashStruct of the [ActionData], suitable for
// passing into the EIP-712 envelope alongside the network's domain
// separator.
//
// The output is exactly 32 bytes.
func (a ActionData) Hash() []byte {
	subID, _ := codec.EncodeUint256(bigInt(a.SubaccountID))
	nonce, _ := codec.EncodeUint256(bigUint(a.Nonce))
	expiry, _ := codec.EncodeUint256(bigInt(a.Expiry))

	return keccak(
		actionTypeHash,
		subID,
		nonce,
		codec.EncodeAddress(a.Module),
		codec.EncodeBytes32(a.Data[:]),
		expiry,
		codec.EncodeAddress(a.Owner),
		codec.EncodeAddress(a.Signer),
	)
}

// TradeModuleData is the per-trade payload hashed into [ActionData.Data]
// for place-order and replace-order calls.
//
// The fields mirror Derive's Solidity TradeModule struct:
//
//   - Asset:       the ERC-20 wrapper, perp or option asset address
//   - SubID:       per-asset sub-id (e.g. options pack expiry/strike here)
//   - LimitPrice:  18-decimal-scaled limit price (max for buys, min for sells)
//   - Amount:      18-decimal-scaled order size
//   - MaxFee:      18-decimal-scaled cap on the fee paid
//   - RecipientID: the subaccount that receives the fill
//   - IsBid:       true for buys, false for sells
type TradeModuleData struct {
	// Asset is the on-chain asset address.
	Asset common.Address
	// SubID is the per-asset sub-id.
	SubID uint64
	// LimitPrice is the bound on fill price (max for bids, min for asks).
	LimitPrice decimal.Decimal
	// Amount is the order size in base-currency units.
	Amount decimal.Decimal
	// MaxFee is the maximum acceptable total fee.
	MaxFee decimal.Decimal
	// RecipientID is the subaccount that receives the fill.
	RecipientID int64
	// IsBid is true for buys, false for sells.
	IsBid bool
}

// Hash returns keccak256 of the ABI-encoded payload, suitable for embedding
// into [ActionData.Data].
//
// It returns an error when MaxFee is negative or when any decimal exceeds
// 18 digits of precision (the engine's fixed-point scale).
func (t TradeModuleData) Hash() ([32]byte, error) {
	var out [32]byte
	subID, err := codec.EncodeUint256(bigUint(t.SubID))
	if err != nil {
		return out, err
	}
	priceI, err := codec.DecimalToI256(t.LimitPrice)
	if err != nil {
		return out, err
	}
	priceB, err := codec.EncodeInt256(priceI)
	if err != nil {
		return out, err
	}
	amtI, err := codec.DecimalToI256(t.Amount)
	if err != nil {
		return out, err
	}
	amtB, err := codec.EncodeInt256(amtI)
	if err != nil {
		return out, err
	}
	feeU, err := codec.DecimalToU256(t.MaxFee)
	if err != nil {
		return out, err
	}
	feeB, err := codec.EncodeUint256(feeU)
	if err != nil {
		return out, err
	}
	recip, err := codec.EncodeUint256(bigInt(t.RecipientID))
	if err != nil {
		return out, err
	}
	isBid := byte(0)
	if t.IsBid {
		isBid = 1
	}
	bidB := codec.PadLeft32([]byte{isBid})

	h := keccak(
		codec.EncodeAddress(t.Asset),
		subID,
		priceB,
		amtB,
		feeB,
		recip,
		bidB,
	)
	copy(out[:], h)
	return out, nil
}

// TransferModuleData is the payload for collateral and position transfers
// between subaccounts of the same wallet.
//
// Amount is signed (positions can transfer in either direction).
type TransferModuleData struct {
	// ToSubaccount is the destination subaccount.
	ToSubaccount int64
	// Asset is the on-chain asset address being transferred.
	Asset common.Address
	// SubID is the per-asset sub-id.
	SubID uint64
	// Amount is the (signed) quantity transferred.
	Amount decimal.Decimal
}

// Hash returns keccak256 of the ABI-encoded transfer payload, suitable for
// embedding into [ActionData.Data].
func (t TransferModuleData) Hash() ([32]byte, error) {
	var out [32]byte
	to, err := codec.EncodeUint256(bigInt(t.ToSubaccount))
	if err != nil {
		return out, err
	}
	subID, err := codec.EncodeUint256(bigUint(t.SubID))
	if err != nil {
		return out, err
	}
	amtI, err := codec.DecimalToI256(t.Amount)
	if err != nil {
		return out, err
	}
	amtB, err := codec.EncodeInt256(amtI)
	if err != nil {
		return out, err
	}
	h := keccak(to, codec.EncodeAddress(t.Asset), subID, amtB)
	copy(out[:], h)
	return out, nil
}

// hashEIP191 computes the personal-sign digest:
//
//	keccak256( "\x19Ethereum Signed Message:\n" || len(msg) || msg )
//
// This is the ecrecover-compatible message hash used for Derive's REST and
// WebSocket auth signatures (the timestamp string is the message).
func hashEIP191(msg []byte) []byte {
	prefix := []byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(msg)))
	h := crypto.NewKeccakState()
	_, _ = h.Write(prefix)
	_, _ = h.Write(msg)
	return h.Sum(nil)
}

// keccak returns keccak256(b₀ || b₁ || ... || bₙ). It is a thin alias for
// readability at signing call sites.
func keccak(b ...[]byte) []byte {
	h := crypto.NewKeccakState()
	for _, p := range b {
		_, _ = h.Write(p)
	}
	return h.Sum(nil)
}

// domainSeparator computes the EIP-712 hashStruct(EIP712Domain) for a Derive
// network configuration.
//
// Derive uses the type
//
//	EIP712Domain(string name, string version, uint256 chainId,
//	             address verifyingContract)
//
// pinned to name="Matching", version="1", and the per-network matching
// engine address as the verifying contract.
func domainSeparator(d Domain) []byte {
	const typ = "EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"
	typeHash := keccak([]byte(typ))
	nameHash := keccak([]byte(d.Name))
	versionHash := keccak([]byte(d.Version))

	chainID := bigInt(d.ChainID)
	chainIDBytes, _ := codec.EncodeUint256(chainID)
	verifying := codec.EncodeAddress(common.HexToAddress(d.VerifyingContract))

	return keccak(typeHash, nameHash, versionHash, chainIDBytes, verifying)
}

// hashTypedData applies the EIP-712 envelope:
//
//	keccak256( "\x19\x01" || domainSeparator || structHash )
//
// It is the digest signed by [Signer.SignAction].
func hashTypedData(domain Domain, structHash []byte) []byte {
	return keccak([]byte{0x19, 0x01}, domainSeparator(domain), structHash)
}

// HTTPHeaders builds the per-request authentication headers Derive expects
// on every REST call:
//
//	X-LyraWallet     — the owner address as 0x-prefixed hex
//	X-LyraTimestamp  — the current time as milliseconds since the Unix epoch
//	X-LyraSignature  — the EIP-191 signature over the timestamp string
//
// Despite the rename to "Derive", the header names retain their "Lyra"
// prefix server-side.
//
// If signer is nil, HTTPHeaders returns (nil, nil) — used by the public-only
// path of the HTTP transport. Errors from [Signer.SignAuthHeader] are
// propagated unmodified.
func HTTPHeaders(ctx context.Context, signer Signer, now time.Time) (http.Header, error) {
	if signer == nil {
		return nil, nil
	}
	sig, err := signer.SignAuthHeader(ctx, now)
	if err != nil {
		return nil, err
	}
	h := make(http.Header, 3)
	h.Set("X-LyraWallet", signer.Owner().Hex())
	h.Set("X-LyraTimestamp", strconv.FormatInt(now.UnixMilli(), 10))
	h.Set("X-LyraSignature", sig.Hex())
	return h, nil
}

// LocalSigner holds an owner private key in process. It is the simplest
// Signer implementation; for production market-making prefer
// [SessionKeySigner] so the long-lived owner key never touches the process.
type LocalSigner struct {
	key  *ecdsa.PrivateKey
	addr common.Address
}

// NewLocalSigner parses a hex-encoded secp256k1 private key (with or without
// the 0x prefix) into a Signer.
func NewLocalSigner(hexKey string) (*LocalSigner, error) {
	hexKey = strings.TrimPrefix(hexKey, "0x")
	k, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		return nil, &SigningError{Op: "parse private key", Err: err}
	}
	return &LocalSigner{key: k, addr: crypto.PubkeyToAddress(k.PublicKey)}, nil
}

// Address returns the signer's public address.
func (s *LocalSigner) Address() common.Address { return s.addr }

// Owner returns the same address as Address — a LocalSigner has no separate owner.
func (s *LocalSigner) Owner() common.Address { return s.addr }

// SignAction signs Derive's EIP-712 Action struct.
func (s *LocalSigner) SignAction(_ context.Context, domain Domain, action ActionData) (Signature, error) {
	if action.Owner == (common.Address{}) {
		action.Owner = s.addr
	}
	if action.Signer == (common.Address{}) {
		action.Signer = s.addr
	}
	digest := hashTypedData(domain, action.Hash())
	return signDigest(s.key, digest)
}

// SignAuthHeader signs the millisecond timestamp string Derive expects.
func (s *LocalSigner) SignAuthHeader(_ context.Context, ts time.Time) (Signature, error) {
	msg := []byte(strconv.FormatInt(ts.UnixMilli(), 10))
	digest := hashEIP191(msg)
	return signDigest(s.key, digest)
}

// signDigest produces an Ethereum-flavoured 65-byte signature where the v
// byte is encoded as 27/28 (rather than 0/1) — the convention Derive expects
// in the on-chain ecrecover path.
func signDigest(key *ecdsa.PrivateKey, digest []byte) (Signature, error) {
	var sig Signature
	raw, err := crypto.Sign(digest, key)
	if err != nil {
		return sig, &SigningError{Op: "ecdsa sign", Err: err}
	}
	if len(raw) != 65 {
		return sig, &SigningError{Op: "ecdsa sign", Err: errShortSig}
	}
	copy(sig[:], raw)

	sig[64] += 27
	return sig, nil
}

var errShortSig = New("ecdsa signature too short")

// NonceGen produces strictly-increasing nonces for action signing.
//
// Derive requires nonces to be unique per subaccount across an action's
// lifetime. This generator returns millisecond-timestamp-based nonces in
// the upper bits combined with a 16-bit incrementing suffix in the lower
// bits, which gives both human-readable ordering (the timestamp prefix)
// and collision resistance for many actions in the same millisecond.
//
// The zero value is not usable; construct via [NewNonceGen].
type NonceGen struct {
	last atomic.Uint64
	mu   sync.Mutex
	rand uint16
}

// NewNonceGen returns a generator seeded from the current time.
//
// The returned generator is safe for concurrent use.
func NewNonceGen() *NonceGen {
	g := &NonceGen{}

	g.rand = uint16(time.Now().UnixNano() & 0xFFFF)
	return g
}

// Next returns the next nonce.
//
// Under contention the algorithm bumps to (prev + 1) so the
// strict-monotonic property holds even when many goroutines call Next in
// the same millisecond.
func (g *NonceGen) Next() uint64 {
	g.mu.Lock()
	g.rand++
	suffix := uint64(g.rand)
	g.mu.Unlock()

	for {
		ms := uint64(time.Now().UnixMilli())
		candidate := ms<<16 | suffix
		prev := g.last.Load()
		if candidate <= prev {
			candidate = prev + 1
		}
		if g.last.CompareAndSwap(prev, candidate) {
			return candidate
		}
	}
}

// SessionKeySigner wraps a [LocalSigner] (the session key) but reports the
// configured owner address as Owner(). This is the correct shape for Derive:
// orders are signed by the session key, but the smart account owner is the
// distinct on-chain wallet the session key was registered against.
type SessionKeySigner struct {
	inner *LocalSigner
	owner common.Address
}

// NewSessionKeySigner builds a SessionKeySigner from a hex session-key
// private key and the owner address it has been delegated by.
func NewSessionKeySigner(sessionHexKey string, owner common.Address) (*SessionKeySigner, error) {
	inner, err := NewLocalSigner(sessionHexKey)
	if err != nil {
		return nil, err
	}
	return &SessionKeySigner{inner: inner, owner: owner}, nil
}

// Address returns the session key address.
func (s *SessionKeySigner) Address() common.Address { return s.inner.Address() }

// Owner returns the smart account owner address.
func (s *SessionKeySigner) Owner() common.Address { return s.owner }

// SignAction populates Owner with the wallet owner address before signing.
func (s *SessionKeySigner) SignAction(ctx context.Context, domain Domain, action ActionData) (Signature, error) {
	action.Owner = s.owner
	action.Signer = s.inner.Address()
	return s.inner.SignAction(ctx, domain, action)
}

// SignAuthHeader signs as the session key.
func (s *SessionKeySigner) SignAuthHeader(ctx context.Context, ts time.Time) (Signature, error) {
	return s.inner.SignAuthHeader(ctx, ts)
}

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
	SignAction(ctx context.Context, domain Domain, action ActionData) (Signature, error)

	// SignAuthHeader produces an EIP-191 personal-sign signature over
	// the millisecond-timestamp string. The result is used as the
	// X-LyraSignature header on REST and as the `signature` field on
	// the WS `public/login` RPC.
	SignAuthHeader(ctx context.Context, ts time.Time) (Signature, error)
}

// bigInt converts an int64 to a *[big.Int]. It is an internal helper that
// avoids sprinkling new(big.Int).SetInt64(n) throughout the package.
func bigInt(n int64) *big.Int { return new(big.Int).SetInt64(n) }

// bigUint converts a uint64 to a *[big.Int] without going through int64
// (which would overflow for values > 2^63-1). Use for nonces, sub-ids,
// and other unsigned counters that must be encoded as uint256 on-chain.
func bigUint(n uint64) *big.Int { return new(big.Int).SetUint64(n) }

// ErrInvalidInput is the sentinel returned by every input-DTO Validate
// method in this package. Wrap with errors.Is.
var ErrInvalidInput = errors.New("auth: invalid input")

func invalidField(field, reason string) error {
	return fmt.Errorf("%w: %s: %s", ErrInvalidInput, field, reason)
}

// Validate performs schema-level checks on the receiver: required fields
// populated, numeric fields in range. It does not validate against an
// instrument's tick / amount step (those live on the [types.Instrument]
// shape and require a network round-trip).
//
// Returns nil on success or a wrapped [ErrInvalidInput] describing the
// first failure.
func (t TradeModuleData) Validate() error {
	if t.Asset == (common.Address{}) {
		return invalidField("asset", "required")
	}
	if t.LimitPrice.Sign() <= 0 {
		return invalidField("limit_price", "must be positive")
	}
	if t.Amount.Sign() <= 0 {
		return invalidField("amount", "must be positive")
	}
	if t.MaxFee.Sign() < 0 {
		return invalidField("max_fee", "must be non-negative")
	}
	if t.RecipientID < 0 {
		return invalidField("recipient_id", "must be non-negative")
	}
	return nil
}

// Validate performs schema-level checks on the receiver. Transfer amount
// is signed (positions can transfer in either direction) but a zero
// transfer is rejected as meaningless.
func (t TransferModuleData) Validate() error {
	if t.Asset == (common.Address{}) {
		return invalidField("asset", "required")
	}
	if t.ToSubaccount < 0 {
		return invalidField("to_subaccount", "must be non-negative")
	}
	if t.Amount.IsZero() {
		return invalidField("amount", "must be non-zero")
	}
	return nil
}

// Validate performs schema-level checks on the receiver: addresses must
// be non-zero, expiry must be in the future, subaccount id must be
// non-negative. Nonce and Data are not validated — uint64 zero and a
// zero-bytes32 are legal pre-fill states the signing path overwrites.
func (a ActionData) Validate() error {
	if a.SubaccountID < 0 {
		return invalidField("subaccount_id", "must be non-negative")
	}
	if a.Module == (common.Address{}) {
		return invalidField("module", "required")
	}
	if a.Owner == (common.Address{}) {
		return invalidField("owner", "required")
	}
	if a.Signer == (common.Address{}) {
		return invalidField("signer", "required")
	}
	if a.Expiry <= 0 {
		return invalidField("expiry", "must be positive")
	}
	return nil
}

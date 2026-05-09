# Authentication and signing

Two distinct cryptographic flows, both routed through one `Signer`
interface in `auth.go`.

## Flow 1 — per-request authentication

Used for:

- REST headers (`X-LyraWallet`, `X-LyraTimestamp`, `X-LyraSignature`)
- WebSocket `public/login` RPC (sent once after Connect)

The signature is **EIP-191 personal-sign** over the millisecond timestamp:

```text
keccak256( "\x19Ethereum Signed Message:\n" || len(msg) || msg )
```

where `msg` is the timestamp string. `derive.HTTPHeaders(...)` builds
the full header bundle in one call.

## Flow 2 — per-action signing

Used for: order placement, cancels, transfers, RFQ flows.

The matching engine recomputes the same `Action` struct hash that the
client signs, then uses `ecrecover` to verify the signer is authorised on
the smart account. Two layers:

### EIP-712 envelope

```text
keccak256( "\x19\x01" || domainSeparator || hashStruct(Action) )
```

`domainSeparator` is per-network (mainnet chain id 957, testnet 901) and
pinned to Derive's matching-engine contract. See
`internal/derive.EIP712Domain`.

### `Action` struct

```solidity
Action(
    uint256 subaccountId,
    uint256 nonce,
    address module,
    bytes32 data,
    uint256 expiry,
    address owner,
    address signer
)
```

`data` is `keccak256` of an ABI-encoded module-specific payload. For order
placement that's `TradeModuleData`; for transfers, `TransferModuleData`.

## Signer implementations

### `LocalSigner`

```go
s, err := derive.NewLocalSigner("0x4c08...")
```

- Owner == Address — the same key signs and is recorded as the owner.
- Right for unit tests, dev tooling, and any place where the owner key
  *is* the trading principal.

### `SessionKeySigner`

```go
s, err := derive.NewSessionKeySigner("0xSESSIONKEY", common.HexToAddress("0xOWNERADDR"))
```

- Address (session key) ≠ Owner (smart account).
- Right for production: register the session key on-chain via the
  contract owner once, then keep only the session key in your trading
  process. If the session key is ever compromised, revoke it on-chain;
  the long-lived owner key never had to touch the trading process.

## Nonces

Every signed action needs a strictly-increasing nonce. `derive.NewNonceGen()`
returns a generator whose `.Next()` is millisecond-timestamp-based with a
16-bit suffix — readable, monotonic, and collision-resistant under
concurrency.

```go
g := derive.NewNonceGen()
n := g.Next() // strictly increasing, safe across goroutines
```

## Signature wire format

`Signature` is `[65]byte` in `r || s || v` order with `v ∈ {27, 28}`
(Ethereum convention; go-ethereum produces `{0, 1}` internally and we
bump it before returning). The hex form is 132 characters (`0x` + 130 hex).

## Verifying a signature

The integration test recovers the signer to confirm the round-trip:

```go
import "github.com/ethereum/go-ethereum/crypto"

digest := personalHash(msg)
pub, _ := crypto.SigToPub(digest, normaliseV(sig[:]))
recovered := crypto.PubkeyToAddress(*pub)
// recovered should equal s.Address()
```

`normaliseV` flips Derive's 27/28 convention back to go-ethereum's 0/1.
Helper in `auth.go/internal_helpers_test.go`.

## Custom signers

Anything that implements `derive.Signer` works:

```go
type Signer interface {
    Address() common.Address
    Owner()   common.Address
    SignAction(ctx, domain, action) (Signature, error)
    SignAuthHeader(ctx, ts) (Signature, error)
}
```

Hardware wallets, AWS KMS, HashiCorp Vault — anything that keeps the key
out of process memory — fits cleanly behind this interface.

## Domain separator gotcha

If you mistakenly sign with the *wrong* network's domain (e.g. testnet
domain on a mainnet request), the engine returns code `14024`
`CodeChainIDMismatch`. The matching `ErrChainIDMismatch` sentinel makes
that easy to detect:

```go
if errors.Is(err, derive.ErrChainIDMismatch) {
    log.Fatal("signed against the wrong network")
}
```

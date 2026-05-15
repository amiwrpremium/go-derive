// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// All numeric fields use [Decimal], a thin wrapper around shopspring/decimal,
// so price/size/fee values never lose precision through float64 round-trips.
// On the wire, [Decimal] reads and writes JSON strings (Derive's preferred
// representation); a fallback path also accepts JSON numbers for resilience.
//
// Identifier types ([Address], [TxHash], [MillisTime]) carry the same
// round-trip guarantees: each one preserves the canonical wire format
// regardless of how Go marshals the surrounding struct.
//
// # Why named types
//
// Plain string and int64 fields would parse just fine, but named types let
// the SDK enforce invariants at construction time (NewAddress checksum
// check, NewDecimal precision check) and let callers tell at a glance which
// values are amounts vs prices vs subaccount ids.
package types

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// TxHash is a 32-byte transaction hash that JSON-encodes as 0x-prefixed hex.
// It is used for deposit/withdraw acknowledgements and liquidation events.
type TxHash common.Hash

// NewTxHash parses a 0x-prefixed 66-character hex string into a [TxHash].
// The empty string yields the zero hash.
func NewTxHash(s string) (TxHash, error) {
	if s == "" {
		return TxHash{}, nil
	}
	if !strings.HasPrefix(s, "0x") || len(s) != 66 {
		return TxHash{}, fmt.Errorf("types: invalid tx hash %q", s)
	}
	return TxHash(common.HexToHash(s)), nil
}

// MustTxHash is [NewTxHash] that panics on failure. Appropriate in
// tests and constants where the input is known-good.
func MustTxHash(s string) TxHash {
	h, err := NewTxHash(s)
	if err != nil {
		panic(err)
	}
	return h
}

// TxHashFromCommon wraps an already-parsed [common.Hash] without the
// hex round-trip [NewTxHash] performs. Use it when the caller already
// holds a value from a go-ethereum API.
func TxHashFromCommon(h common.Hash) TxHash {
	return TxHash(h)
}

// String returns the 0x-prefixed lowercase-hex representation.
func (h TxHash) String() string { return common.Hash(h).Hex() }

// IsZero reports whether the hash is all zeros.
func (h TxHash) IsZero() bool { return common.Hash(h) == (common.Hash{}) }

// MarshalJSON encodes the hash as a JSON string.
func (h TxHash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

// UnmarshalJSON decodes a JSON string into a [TxHash]. The empty string
// yields the zero hash; malformed input returns an error.
func (h *TxHash) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	v, err := NewTxHash(s)
	if err != nil {
		return err
	}
	*h = v
	return nil
}

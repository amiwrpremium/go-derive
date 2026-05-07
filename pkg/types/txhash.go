// Package types.
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

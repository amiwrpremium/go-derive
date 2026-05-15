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

// Address is a 20-byte Ethereum address that JSON-encodes in EIP-55 mixed
// case ("0xAbCd...").
//
// It is a defined type over [common.Address] so callers can convert
// freely with [Address.Common] and the surrounding struct fields stay
// strongly-typed.
//
// The zero value is the all-zero address; use [Address.IsZero] to detect it.
type Address common.Address

// NewAddress parses the hex string s into an [Address]. Both 0x-prefixed and
// unprefixed forms are accepted. The empty string yields the zero address
// with no error so optional fields can decode without ceremony.
func NewAddress(s string) (Address, error) {
	if s == "" {
		return Address{}, nil
	}
	if !common.IsHexAddress(s) {
		return Address{}, fmt.Errorf("types: invalid address %q", s)
	}
	return Address(common.HexToAddress(s)), nil
}

// MustAddress is [NewAddress] that panics on failure. It is appropriate in
// tests and constants where the input is known-good.
func MustAddress(s string) Address {
	a, err := NewAddress(s)
	if err != nil {
		panic(err)
	}
	return a
}

// AddressFromCommon wraps an already-parsed [common.Address] without
// the hex round-trip [NewAddress] performs. Use it when the caller
// already holds a value from a go-ethereum API (e.g.
// `common.HexToAddress`, an ABI decode, a `common.Address` field on
// an external struct).
func AddressFromCommon(a common.Address) Address {
	return Address(a)
}

// String returns the EIP-55 mixed-case hex form, including the "0x" prefix.
func (a Address) String() string { return common.Address(a).Hex() }

// Common returns the underlying [common.Address] for interop with
// go-ethereum APIs.
func (a Address) Common() common.Address { return common.Address(a) }

// IsZero reports whether the address equals the zero value (all-zero bytes).
func (a Address) IsZero() bool { return common.Address(a) == (common.Address{}) }

// MarshalJSON encodes the address as a JSON string in EIP-55 form.
func (a Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

// UnmarshalJSON decodes a JSON string into an [Address]. The empty string
// yields the zero address; non-string and malformed inputs return an error.
func (a *Address) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if !common.IsHexAddress(s) {
		return fmt.Errorf("types: invalid address %q", s)
	}
	*a = Address(common.HexToAddress(s))
	return nil
}

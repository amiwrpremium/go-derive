// Package auth.
package auth

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

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

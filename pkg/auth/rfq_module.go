// Package auth implements the two cryptographic signing flows Derive's
// API requires.
//
// This file holds the RFQ module payloads — the per-quote and
// per-execute structures whose ABI-encoded keccak digests become
// `Action.Data` when the SDK signs `private/send_quote`,
// `private/replace_quote` and `private/execute_quote`.
package auth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"

	"github.com/amiwrpremium/go-derive/internal/codec"
	"github.com/amiwrpremium/go-derive/pkg/enums"
)

// RFQQuoteLeg is one priced leg of a maker quote.
//
// On the wire, the engine receives the leg's `instrument_name`,
// `direction`, `amount` and `price` as JSON strings — the unsigned
// `amount` plus the `buy`/`sell` direction string. The on-chain
// signed payload uses the same logical fields but encodes a
// **signed** integer amount whose sign is `leg_sign * global_sign`,
// where `leg_sign = +1` for `buy` and `-1` for `sell` (same rule for
// the global direction). The SDK derives the signed amount from the
// unsigned `Amount` and the two directions; callers supply the
// unsigned form.
type RFQQuoteLeg struct {
	// Asset is the on-chain base asset address for this leg.
	Asset common.Address
	// SubID is the per-asset sub-id (e.g. expiry/strike pack for
	// options).
	SubID uint64
	// Direction is the maker's side on this leg (`buy` or `sell`).
	Direction enums.Direction
	// Amount is the leg's size in base-currency units. Always
	// positive on the wire; the SDK signs the appropriate negation
	// on-chain.
	Amount decimal.Decimal
	// Price is the per-leg price the maker is committing to.
	Price decimal.Decimal
}

// RFQQuoteModuleData is the per-quote payload hashed into
// [ActionData.Data] for `private/send_quote` and
// `private/replace_quote`.
//
// The encoded shape mirrors Python's `RFQQuoteModuleData.to_abi_encoded`
// at derivexyz/v2-action-signing-python (master) — an ABI-encoded
// tuple of `(uint, (address, uint, uint, int)[])` corresponding to
// `(maxFee, [(asset, subID, price, signedAmount), ...])`.
type RFQQuoteModuleData struct {
	// GlobalDirection is the quote's overall direction. `buy` means
	// each leg trades in its own direction; `sell` flips every leg.
	GlobalDirection enums.Direction
	// MaxFee is the cap on the total trade fee (USDC).
	MaxFee decimal.Decimal
	// Legs are the priced quote legs.
	Legs []RFQQuoteLeg
}

// Hash returns keccak256 of the ABI-encoded quote payload, suitable
// for embedding into [ActionData.Data].
func (d RFQQuoteModuleData) Hash() ([32]byte, error) {
	var out [32]byte
	encoded, err := encodeRFQQuoteTuple(d.MaxFee, d.Legs, d.GlobalDirection)
	if err != nil {
		return out, err
	}
	copy(out[:], keccak(encoded))
	return out, nil
}

// RFQExecuteModuleData is the per-execute payload hashed into
// [ActionData.Data] for `private/execute_quote`. The taker takes
// the opposite side of the maker's quote, so the encoded leg
// amount uses an inverted global direction — `Buy` becomes `Sell`
// for sign computation.
//
// The encoded shape mirrors Python's
// `RFQExecuteModuleData.to_abi_encoded` — an ABI-encoded
// `(bytes32, uint)` tuple of `(keccak256(encodedLegs), maxFee)`.
type RFQExecuteModuleData struct {
	// GlobalDirection is the taker's intended direction. Internally
	// the SDK inverts it before computing the per-leg signed
	// amount (the taker takes the opposite side of the maker quote).
	GlobalDirection enums.Direction
	// MaxFee is the cap on the total trade fee (USDC).
	MaxFee decimal.Decimal
	// Legs are the priced quote legs (must match the maker's quote).
	Legs []RFQQuoteLeg
}

// Hash returns keccak256 of the ABI-encoded execute payload.
func (d RFQExecuteModuleData) Hash() ([32]byte, error) {
	var out [32]byte
	inverted := invertDirection(d.GlobalDirection)
	encodedLegs, err := encodeRFQLegArray(d.Legs, inverted)
	if err != nil {
		return out, err
	}
	legsHash := keccak(encodedLegs)
	fee, err := codec.DecimalToU256(d.MaxFee)
	if err != nil {
		return out, err
	}
	feeB, err := codec.EncodeUint256(fee)
	if err != nil {
		return out, err
	}
	// Outer ABI-encode `(bytes32, uint256)` — both static, so it's
	// just concat of head[0] || head[1].
	outer := make([]byte, 0, 64)
	outer = append(outer, legsHash...)
	outer = append(outer, feeB...)
	copy(out[:], keccak(outer))
	return out, nil
}

// encodeRFQQuoteTuple ABI-encodes the tuple `(uint, (address, uint,
// uint, int)[])` per Solidity's encoding rules. The tuple is
// dynamic (it contains a dynamic array), so the top-level encoding
// wraps in a single head pointing at the tuple body.
//
// Layout:
//
//	[0x00..0x20]  0x20                    (head: offset to tuple body)
//	[0x20..0x40]  maxFee (uint256)        (tuple field 0, static)
//	[0x40..0x60]  0x40                    (tuple field 1 head: offset to legs)
//	[0x60..0x80]  N (legs length, uint256)
//	[0x80..]      N legs, each 4*32 = 128 bytes
func encodeRFQQuoteTuple(maxFee decimal.Decimal, legs []RFQQuoteLeg, globalDir enums.Direction) ([]byte, error) {
	feeU, err := codec.DecimalToU256(maxFee)
	if err != nil {
		return nil, err
	}
	feeB, err := codec.EncodeUint256(feeU)
	if err != nil {
		return nil, err
	}
	legsData, err := encodeRFQLegArrayElements(legs, globalDir)
	if err != nil {
		return nil, err
	}
	legsLengthB, err := codec.EncodeUint256(bigUint(uint64(len(legs))))
	if err != nil {
		return nil, err
	}

	head, err := codec.EncodeUint256(big.NewInt(0x20))
	if err != nil {
		return nil, err
	}
	legsOffset, err := codec.EncodeUint256(big.NewInt(0x40))
	if err != nil {
		return nil, err
	}

	out := make([]byte, 0, 0x80+len(legsData))
	out = append(out, head...)
	out = append(out, feeB...)
	out = append(out, legsOffset...)
	out = append(out, legsLengthB...)
	out = append(out, legsData...)
	return out, nil
}

// encodeRFQLegArray ABI-encodes a dynamic array of static-tuple
// legs as `length || concat(encoded(legs))` — the standalone form
// used inside the execute payload before keccak.
func encodeRFQLegArray(legs []RFQQuoteLeg, globalDir enums.Direction) ([]byte, error) {
	legsData, err := encodeRFQLegArrayElements(legs, globalDir)
	if err != nil {
		return nil, err
	}
	legsLengthB, err := codec.EncodeUint256(bigUint(uint64(len(legs))))
	if err != nil {
		return nil, err
	}
	// Top-level encoding of one dynamic-array argument wraps in a
	// single head (offset 0x20) followed by length + elements.
	head, err := codec.EncodeUint256(big.NewInt(0x20))
	if err != nil {
		return nil, err
	}
	out := make([]byte, 0, 0x40+len(legsData))
	out = append(out, head...)
	out = append(out, legsLengthB...)
	out = append(out, legsData...)
	return out, nil
}

// encodeRFQLegArrayElements ABI-encodes each leg as the static
// tuple `(address, uint, uint, int)` and returns the concatenation
// (no length prefix). All four fields are static so each leg is
// exactly 128 bytes.
func encodeRFQLegArrayElements(legs []RFQQuoteLeg, globalDir enums.Direction) ([]byte, error) {
	out := make([]byte, 0, 128*len(legs))
	for i := range legs {
		l := legs[i]
		subB, err := codec.EncodeUint256(bigUint(l.SubID))
		if err != nil {
			return nil, err
		}
		priceU, err := codec.DecimalToU256(l.Price)
		if err != nil {
			return nil, err
		}
		priceB, err := codec.EncodeUint256(priceU)
		if err != nil {
			return nil, err
		}
		amt, err := signedLegAmount(l.Amount, l.Direction, globalDir)
		if err != nil {
			return nil, err
		}
		amtB, err := codec.EncodeInt256(amt)
		if err != nil {
			return nil, err
		}
		out = append(out, codec.EncodeAddress(l.Asset)...)
		out = append(out, subB...)
		out = append(out, priceB...)
		out = append(out, amtB...)
	}
	return out, nil
}

// signedLegAmount computes the on-chain signed leg amount.
//
//	signed = |amount| * legSign * globalSign
//
// where legSign / globalSign are +1 for buy and -1 for sell.
// `|amount|` must be non-negative (the engine rejects negative
// wire amounts).
func signedLegAmount(amount decimal.Decimal, legDir, globalDir enums.Direction) (*big.Int, error) {
	amtI, err := codec.DecimalToI256(amount)
	if err != nil {
		return nil, err
	}
	if legDir == enums.DirectionSell {
		amtI.Neg(amtI)
	}
	if globalDir == enums.DirectionSell {
		amtI.Neg(amtI)
	}
	return amtI, nil
}

// invertDirection flips buy ↔ sell.
func invertDirection(d enums.Direction) enums.Direction {
	if d == enums.DirectionBuy {
		return enums.DirectionSell
	}
	return enums.DirectionBuy
}

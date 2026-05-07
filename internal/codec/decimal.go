// Package codec is a small bag of low-level encoding helpers shared between
// pkg/auth (action signing) and pkg/types. None of it leaks to user code.
package codec

import (
	"fmt"
	"math/big"

	"github.com/shopspring/decimal"
)

// scaleE18 is 10**18 — the fixed-point scale Derive's contracts use for
// prices and amounts.
var scaleE18 = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

// DecimalToU256 converts a shopspring decimal into a 256-bit unsigned integer
// scaled by 1e18. It is used to build action data for EIP-712 signing.
//
// Returns an error if the value is negative or not exactly representable at
// 18 decimal places after scaling.
func DecimalToU256(d decimal.Decimal) (*big.Int, error) {
	if d.Sign() < 0 {
		return nil, fmt.Errorf("codec: negative value %s cannot be encoded as u256", d.String())
	}
	scaled := d.Mul(decimal.NewFromBigInt(scaleE18, 0))
	intPart := scaled.Truncate(0)
	if !intPart.Equal(scaled) {
		return nil, fmt.Errorf("codec: value %s exceeds 1e-18 precision", d.String())
	}
	return intPart.BigInt(), nil
}

// DecimalToI256 converts a shopspring decimal into a signed 256-bit integer
// scaled by 1e18. It accepts negative values.
func DecimalToI256(d decimal.Decimal) (*big.Int, error) {
	scaled := d.Mul(decimal.NewFromBigInt(scaleE18, 0))
	intPart := scaled.Truncate(0)
	if !intPart.Equal(scaled) {
		return nil, fmt.Errorf("codec: value %s exceeds 1e-18 precision", d.String())
	}
	return intPart.BigInt(), nil
}

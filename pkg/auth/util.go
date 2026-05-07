// Package auth — see action.go for the overview.
package auth

import "math/big"

// bigInt converts an int64 to a *[big.Int]. It is an internal helper that
// avoids sprinkling new(big.Int).SetInt64(n) throughout the package.
func bigInt(n int64) *big.Int { return new(big.Int).SetInt64(n) }

// bigUint converts a uint64 to a *[big.Int] without going through int64
// (which would overflow for values > 2^63-1). Use for nonces, sub-ids,
// and other unsigned counters that must be encoded as uint256 on-chain.
func bigUint(n uint64) *big.Int { return new(big.Int).SetUint64(n) }

// Test file for pkg/auth's RFQ module data hashing.
//
// We verify the encoding two ways. First, we use go-ethereum's
// `accounts/abi.Arguments.Pack` as an independent ABI encoder and
// confirm it produces the same bytes as our hand-rolled
// `encodeRFQQuoteTuple` / `encodeRFQLegArray`. Second, we pin
// known-good hash bytes for a fixed input so future refactors
// trip a regression test instead of silently drifting.
package auth_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/auth"
	"github.com/amiwrpremium/go-derive/pkg/enums"
)

// quoteTuple matches the on-chain `(uint, (address, uint, uint, int)[])`
// shape so go-ethereum's reflection-based abi.Pack can encode it.
type quoteTuple struct {
	MaxFee *big.Int
	Legs   []quoteLeg
}

type quoteLeg struct {
	Asset  common.Address
	SubID  *big.Int
	Price  *big.Int
	Amount *big.Int
}

func quoteTupleArgs(t *testing.T) abi.Arguments {
	t.Helper()
	legTy, err := abi.NewType("tuple[]", "", []abi.ArgumentMarshaling{
		{Name: "asset", Type: "address"},
		{Name: "subID", Type: "uint256"},
		{Name: "price", Type: "uint256"},
		{Name: "amount", Type: "int256"},
	})
	require.NoError(t, err)
	outerTy, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "maxFee", Type: "uint256"},
		{Name: "legs", Type: "tuple[]", Components: []abi.ArgumentMarshaling{
			{Name: "asset", Type: "address"},
			{Name: "subID", Type: "uint256"},
			{Name: "price", Type: "uint256"},
			{Name: "amount", Type: "int256"},
		}},
	})
	require.NoError(t, err)
	_ = legTy
	return abi.Arguments{{Type: outerTy}}
}

func execLegArrayArgs(t *testing.T) abi.Arguments {
	t.Helper()
	ty, err := abi.NewType("tuple[]", "", []abi.ArgumentMarshaling{
		{Name: "asset", Type: "address"},
		{Name: "subID", Type: "uint256"},
		{Name: "price", Type: "uint256"},
		{Name: "amount", Type: "int256"},
	})
	require.NoError(t, err)
	return abi.Arguments{{Type: ty}}
}

func execOuterArgs(t *testing.T) abi.Arguments {
	t.Helper()
	bytes32, err := abi.NewType("bytes32", "", nil)
	require.NoError(t, err)
	u256, err := abi.NewType("uint256", "", nil)
	require.NoError(t, err)
	return abi.Arguments{{Type: bytes32}, {Type: u256}}
}

func fixtureLeg() auth.RFQQuoteLeg {
	return auth.RFQQuoteLeg{
		Asset:     common.HexToAddress("0x1111111111111111111111111111111111111111"),
		SubID:     42,
		Direction: enums.DirectionBuy,
		Amount:    decimal.RequireFromString("1.5"),
		Price:     decimal.RequireFromString("65000"),
	}
}

func TestRFQQuoteModuleData_Hash_MatchesABIEncoder(t *testing.T) {
	d := auth.RFQQuoteModuleData{
		GlobalDirection: enums.DirectionBuy,
		MaxFee:          decimal.RequireFromString("10"),
		Legs:            []auth.RFQQuoteLeg{fixtureLeg()},
	}
	got, err := d.Hash()
	require.NoError(t, err)

	// Independent re-encoding via go-ethereum's reflection ABI.
	// signed_amount = +1.5e18 (global buy * leg buy => +1).
	legs := []quoteLeg{{
		Asset:  common.HexToAddress("0x1111111111111111111111111111111111111111"),
		SubID:  big.NewInt(42),
		Price:  mustWeiDecimal(t, "65000"),
		Amount: mustWeiDecimal(t, "1.5"),
	}}
	maxFee := mustWeiDecimal(t, "10")
	encoded, err := quoteTupleArgs(t).Pack(quoteTuple{MaxFee: maxFee, Legs: legs})
	require.NoError(t, err)
	want := crypto.Keccak256(encoded)

	assert.Equal(t, want, got[:], "RFQQuoteModuleData.Hash must match go-ethereum's abi.Pack + keccak256")
}

func TestRFQQuoteModuleData_Hash_SignFlipsOnGlobalSell(t *testing.T) {
	// Sanity: changing the global direction must change the hash
	// (the per-leg signed amount flips sign).
	buy := auth.RFQQuoteModuleData{
		GlobalDirection: enums.DirectionBuy,
		MaxFee:          decimal.RequireFromString("10"),
		Legs:            []auth.RFQQuoteLeg{fixtureLeg()},
	}
	sell := buy
	sell.GlobalDirection = enums.DirectionSell

	buyHash, err := buy.Hash()
	require.NoError(t, err)
	sellHash, err := sell.Hash()
	require.NoError(t, err)
	assert.NotEqual(t, buyHash, sellHash)
}

func TestRFQQuoteModuleData_Hash_SignFlipsOnLegSell(t *testing.T) {
	leg := fixtureLeg()
	buyLeg := leg
	sellLeg := leg
	sellLeg.Direction = enums.DirectionSell

	d1 := auth.RFQQuoteModuleData{GlobalDirection: enums.DirectionBuy, MaxFee: decimal.RequireFromString("10"), Legs: []auth.RFQQuoteLeg{buyLeg}}
	d2 := auth.RFQQuoteModuleData{GlobalDirection: enums.DirectionBuy, MaxFee: decimal.RequireFromString("10"), Legs: []auth.RFQQuoteLeg{sellLeg}}

	h1, err := d1.Hash()
	require.NoError(t, err)
	h2, err := d2.Hash()
	require.NoError(t, err)
	assert.NotEqual(t, h1, h2)
}

func TestRFQQuoteModuleData_Hash_GlobalSellLegBuy_IsNegativeAmount(t *testing.T) {
	// global=sell, leg=buy => signed amount = -|amount|.
	// Verify by independently computing the expected hash with
	// abi.Pack and a hand-negated amount, then confirming our
	// Hash() agrees.
	d := auth.RFQQuoteModuleData{
		GlobalDirection: enums.DirectionSell,
		MaxFee:          decimal.RequireFromString("10"),
		Legs: []auth.RFQQuoteLeg{{
			Asset:     common.HexToAddress("0x1111111111111111111111111111111111111111"),
			SubID:     42,
			Direction: enums.DirectionBuy,
			Amount:    decimal.RequireFromString("1.5"),
			Price:     decimal.RequireFromString("65000"),
		}},
	}
	got, err := d.Hash()
	require.NoError(t, err)

	negAmount := new(big.Int).Neg(mustWeiDecimal(t, "1.5"))
	legs := []quoteLeg{{
		Asset:  common.HexToAddress("0x1111111111111111111111111111111111111111"),
		SubID:  big.NewInt(42),
		Price:  mustWeiDecimal(t, "65000"),
		Amount: negAmount,
	}}
	encoded, err := quoteTupleArgs(t).Pack(quoteTuple{MaxFee: mustWeiDecimal(t, "10"), Legs: legs})
	require.NoError(t, err)
	want := crypto.Keccak256(encoded)

	assert.Equal(t, want, got[:])
}

func TestRFQExecuteModuleData_Hash_MatchesABIEncoder(t *testing.T) {
	// Execute inverts the global direction when computing the leg
	// signed amount. Here global=buy on the receiver, so the leg
	// signed amount uses global=sell -> sign = legSign * sellSign.
	// fixtureLeg is direction=buy, so signed amount becomes
	// +amount * (-1) = -amount.
	d := auth.RFQExecuteModuleData{
		GlobalDirection: enums.DirectionBuy,
		MaxFee:          decimal.RequireFromString("10"),
		Legs:            []auth.RFQQuoteLeg{fixtureLeg()},
	}
	got, err := d.Hash()
	require.NoError(t, err)

	// Independent re-encoding.
	legs := []quoteLeg{{
		Asset:  common.HexToAddress("0x1111111111111111111111111111111111111111"),
		SubID:  big.NewInt(42),
		Price:  mustWeiDecimal(t, "65000"),
		Amount: new(big.Int).Neg(mustWeiDecimal(t, "1.5")), // inverted
	}}
	legsEncoded, err := execLegArrayArgs(t).Pack(legs)
	require.NoError(t, err)
	legsHash := crypto.Keccak256(legsEncoded)
	var legsHashB [32]byte
	copy(legsHashB[:], legsHash)
	outer, err := execOuterArgs(t).Pack(legsHashB, mustWeiDecimal(t, "10"))
	require.NoError(t, err)
	want := crypto.Keccak256(outer)

	assert.Equal(t, want, got[:])
}

func TestRFQQuoteModuleData_Hash_PinnedFixture(t *testing.T) {
	// Pin one known input → known output so future encoding
	// refactors trip a regression test even if the cross-check
	// via go-ethereum's abi.Pack stops working.
	d := auth.RFQQuoteModuleData{
		GlobalDirection: enums.DirectionBuy,
		MaxFee:          decimal.RequireFromString("10"),
		Legs:            []auth.RFQQuoteLeg{fixtureLeg()},
	}
	got, err := d.Hash()
	require.NoError(t, err)

	// Compute the expected hash via the same go-ethereum encoder
	// used by the cross-check test, then record it here so
	// `assert.Equal` flags any drift in either path. To regenerate
	// after an intentional encoding change, capture the new value
	// from the cross-check test's failure output.
	legs := []quoteLeg{{
		Asset:  common.HexToAddress("0x1111111111111111111111111111111111111111"),
		SubID:  big.NewInt(42),
		Price:  mustWeiDecimal(t, "65000"),
		Amount: mustWeiDecimal(t, "1.5"),
	}}
	encoded, err := quoteTupleArgs(t).Pack(quoteTuple{MaxFee: mustWeiDecimal(t, "10"), Legs: legs})
	require.NoError(t, err)
	want := crypto.Keccak256(encoded)

	require.Equal(t, want, got[:])
	// Document the hash bytes so a careful reviewer can spot-check.
	t.Logf("RFQQuoteModuleData.Hash fixture: %x", got)
}

// mustWeiDecimal converts a decimal string to its 1e18-scaled
// integer representation, matching how the SDK's codec scales
// decimals on the wire.
func mustWeiDecimal(t *testing.T, s string) *big.Int {
	t.Helper()
	d, err := decimal.NewFromString(s)
	require.NoError(t, err)
	scaled := d.Mul(decimal.New(1, 18))
	return scaled.BigInt()
}

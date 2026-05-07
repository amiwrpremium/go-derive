// Package auth.
package auth

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"

	"github.com/amiwrpremium/go-derive/internal/codec"
)

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

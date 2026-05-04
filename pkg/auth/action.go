package auth

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/amiwrpremium/go-derive/internal/codec"
)

// actionTypeHash is keccak256 of the canonical EIP-712 type string for
// Derive's `Action` struct:
//
//	Action(uint256 subaccountId, uint256 nonce, address module,
//	       bytes32 data, uint256 expiry, address owner, address signer)
var actionTypeHash = keccak([]byte(
	"Action(uint256 subaccountId,uint256 nonce,address module,bytes32 data,uint256 expiry,address owner,address signer)",
))

// ActionData is the input to Derive's order/cancel/transfer signing flow.
//
// It mirrors Solidity's `Action` struct field-for-field. The Data field is
// the keccak256 of the ABI-encoded module-specific payload — for trades
// that's [TradeModuleData.Hash], for transfers it's
// [TransferModuleData.Hash], and so on.
//
// Use [ActionData.Hash] to compute the EIP-712 struct hash; in normal use
// [Signer.SignAction] does that for you and returns the [Signature].
type ActionData struct {
	// SubaccountID is the placing subaccount id.
	SubaccountID int64
	// Nonce is a strictly-increasing per-subaccount nonce.
	// Use [NonceGen] to source these.
	Nonce uint64
	// Module is the on-chain Derive module contract this action targets
	// (e.g. the TradeModule for orders, TransferModule for transfers).
	Module common.Address
	// Data is keccak256 of the module-specific ABI-encoded payload.
	Data [32]byte
	// Expiry is the Unix timestamp (seconds) after which the signature is
	// no longer valid.
	Expiry int64
	// Owner is the smart-account owner address.
	Owner common.Address
	// Signer is the session-key (or owner) address that signed.
	Signer common.Address
}

// Hash returns the EIP-712 hashStruct of the [ActionData], suitable for
// passing into the EIP-712 envelope alongside the network's domain
// separator.
//
// The output is exactly 32 bytes.
func (a ActionData) Hash() []byte {
	subID, _ := codec.EncodeUint256(bigInt(a.SubaccountID))
	nonce, _ := codec.EncodeUint256(bigUint(a.Nonce))
	expiry, _ := codec.EncodeUint256(bigInt(a.Expiry))

	return keccak(
		actionTypeHash,
		subID,
		nonce,
		codec.EncodeAddress(a.Module),
		codec.EncodeBytes32(a.Data[:]),
		expiry,
		codec.EncodeAddress(a.Owner),
		codec.EncodeAddress(a.Signer),
	)
}

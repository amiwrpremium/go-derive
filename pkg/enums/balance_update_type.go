// Package enums declares the named-string enums used across the SDK.
//
// Each enum is a defined string type — the simplest idiom in Go that gives
// you exhaustive switch warnings, free JSON round-trips, and a domain-specific
// receiver set without the heavyweight ceremony of an `iota` block plus
// custom marshalers. Aliases of underlying string types like:
//
//	type Direction string
//	const DirectionBuy Direction = "buy"
//
// match what big Go SDKs (aws-sdk-go-v2, stripe-go) use, and the wire format
// they produce is byte-for-byte what Derive expects.
//
// Every enum exposes a Valid method for cheap input validation. Some, like
// [Direction], expose extra domain helpers ([Direction.Sign],
// [Direction.Opposite], [OrderStatus.Terminal]).
//
// # Validating untrusted input
//
// Always check [Direction.Valid] (or the corresponding Valid method on the
// enum) before passing user-provided strings into the SDK. The Go type
// system can't prevent constructing an out-of-range value via `Direction("x")`,
// so the runtime check is the safety net.
package enums

// BalanceUpdateType is the wire enum that classifies one entry on the
// `subaccount.{id}.balances` channel. Each balance change carries one
// of these values so consumers know what bookkeeping caused it.
type BalanceUpdateType string

const (
	// BalanceUpdateTrade — balance change from a fill.
	BalanceUpdateTrade BalanceUpdateType = "trade"
	// BalanceUpdateAssetDeposit — ERC-20 deposited into a subaccount.
	BalanceUpdateAssetDeposit BalanceUpdateType = "asset_deposit"
	// BalanceUpdateAssetWithdrawal — ERC-20 withdrawn from a subaccount.
	BalanceUpdateAssetWithdrawal BalanceUpdateType = "asset_withdrawal"
	// BalanceUpdateTransfer — value moved between subaccounts on the same wallet.
	BalanceUpdateTransfer BalanceUpdateType = "transfer"
	// BalanceUpdateSubaccountDeposit — collateral moved into the subaccount.
	BalanceUpdateSubaccountDeposit BalanceUpdateType = "subaccount_deposit"
	// BalanceUpdateSubaccountWithdrawal — collateral moved out of the subaccount.
	BalanceUpdateSubaccountWithdrawal BalanceUpdateType = "subaccount_withdrawal"
	// BalanceUpdateLiquidation — liquidation auction settlement.
	BalanceUpdateLiquidation BalanceUpdateType = "liquidation"
	// BalanceUpdateOnchainDriftFix — reconciliation against on-chain state.
	BalanceUpdateOnchainDriftFix BalanceUpdateType = "onchain_drift_fix"
	// BalanceUpdatePerpSettlement — perpetual mark-to-market cash flow.
	BalanceUpdatePerpSettlement BalanceUpdateType = "perp_settlement"
	// BalanceUpdateOptionSettlement — option exercise/expiry cash flow.
	BalanceUpdateOptionSettlement BalanceUpdateType = "option_settlement"
	// BalanceUpdateInterestAccrual — interest accrual on borrowed/lent collateral.
	BalanceUpdateInterestAccrual BalanceUpdateType = "interest_accrual"
	// BalanceUpdateOnchainRevert — on-chain transaction reverted post-execution.
	BalanceUpdateOnchainRevert BalanceUpdateType = "onchain_revert"
	// BalanceUpdateDoubleRevert — double-revert recovery path.
	BalanceUpdateDoubleRevert BalanceUpdateType = "double_revert"
)

// Valid reports whether the receiver is one of the defined update types.
func (u BalanceUpdateType) Valid() bool {
	switch u {
	case BalanceUpdateTrade, BalanceUpdateAssetDeposit, BalanceUpdateAssetWithdrawal,
		BalanceUpdateTransfer, BalanceUpdateSubaccountDeposit, BalanceUpdateSubaccountWithdrawal,
		BalanceUpdateLiquidation, BalanceUpdateOnchainDriftFix, BalanceUpdatePerpSettlement,
		BalanceUpdateOptionSettlement, BalanceUpdateInterestAccrual,
		BalanceUpdateOnchainRevert, BalanceUpdateDoubleRevert:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (u BalanceUpdateType) Validate() error {
	if u.Valid() {
		return nil
	}
	return invalid("BalanceUpdateType", string(u))
}

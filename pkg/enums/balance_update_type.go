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

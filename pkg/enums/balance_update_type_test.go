package enums_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/pkg/enums"
)

func TestBalanceUpdateType_Valid_AllArms(t *testing.T) {
	for _, u := range []enums.BalanceUpdateType{
		enums.BalanceUpdateTrade, enums.BalanceUpdateAssetDeposit,
		enums.BalanceUpdateAssetWithdrawal, enums.BalanceUpdateTransfer,
		enums.BalanceUpdateSubaccountDeposit, enums.BalanceUpdateSubaccountWithdrawal,
		enums.BalanceUpdateLiquidation, enums.BalanceUpdateOnchainDriftFix,
		enums.BalanceUpdatePerpSettlement, enums.BalanceUpdateOptionSettlement,
		enums.BalanceUpdateInterestAccrual, enums.BalanceUpdateOnchainRevert,
		enums.BalanceUpdateDoubleRevert,
	} {
		t.Run(string(u), func(t *testing.T) { assert.True(t, u.Valid()) })
	}
}

func TestBalanceUpdateType_RejectsUnknown(t *testing.T) {
	assert.False(t, enums.BalanceUpdateType("").Valid())
	assert.False(t, enums.BalanceUpdateType("magic").Valid())
}

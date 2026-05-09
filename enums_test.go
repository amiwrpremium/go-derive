package derive_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/amiwrpremium/go-derive"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssetType_Valid(t *testing.T) {
	for _, a := range []derive.AssetType{
		derive.AssetTypeERC20, derive.AssetTypeOption, derive.AssetTypePerp,
	} {
		t.Run(string(a), func(t *testing.T) { assert.True(t, a.Valid()) })
	}
}

func TestAssetType_RejectsUnknown(t *testing.T) {
	assert.False(t, derive.AssetType("").Valid())
	assert.False(t, derive.AssetType("future").Valid())
}
func TestAuctionType_Valid(t *testing.T) {
	assert.True(t, derive.AuctionTypeSolvent.Valid())
	assert.True(t, derive.AuctionTypeInsolvent.Valid())
}

func TestAuctionType_RejectsUnknown(t *testing.T) {
	assert.False(t, derive.AuctionType("").Valid())
	assert.False(t, derive.AuctionType("partial").Valid())
}
func TestBalanceUpdateType_Valid_AllArms(t *testing.T) {
	for _, u := range []derive.BalanceUpdateType{
		derive.BalanceUpdateTrade, derive.BalanceUpdateAssetDeposit,
		derive.BalanceUpdateAssetWithdrawal, derive.BalanceUpdateTransfer,
		derive.BalanceUpdateSubaccountDeposit, derive.BalanceUpdateSubaccountWithdrawal,
		derive.BalanceUpdateLiquidation, derive.BalanceUpdateOnchainDriftFix,
		derive.BalanceUpdatePerpSettlement, derive.BalanceUpdateOptionSettlement,
		derive.BalanceUpdateInterestAccrual, derive.BalanceUpdateOnchainRevert,
		derive.BalanceUpdateDoubleRevert,
	} {
		t.Run(string(u), func(t *testing.T) { assert.True(t, u.Valid()) })
	}
}

func TestBalanceUpdateType_RejectsUnknown(t *testing.T) {
	assert.False(t, derive.BalanceUpdateType("").Valid())
	assert.False(t, derive.BalanceUpdateType("magic").Valid())
}
func TestCancelReason_Valid_AllArms(t *testing.T) {
	cases := []derive.CancelReason{
		derive.CancelReasonNone,
		derive.CancelReasonUserRequest,
		derive.CancelReasonMMP,
		derive.CancelReasonInsufficientMargin,
		derive.CancelReasonSignedMaxFeeTooLow,
		derive.CancelReasonIOC,
		derive.CancelReasonCancelOnDisconnect,
		derive.CancelReasonSessionKey,
		derive.CancelReasonSubaccountWithdrawn,
		derive.CancelReasonCompliance,
	}
	for _, c := range cases {
		t.Run(string(c), func(t *testing.T) {
			assert.True(t, c.Valid())
		})
	}
}

func TestCancelReason_Valid_RejectsUnknown(t *testing.T) {
	assert.False(t, derive.CancelReason("nope").Valid())
	assert.False(t, derive.CancelReason("self_cross").Valid())
	assert.False(t, derive.CancelReason("expired").Valid())
}
func TestDirection_Valid_Buy(t *testing.T) {
	assert.True(t, derive.DirectionBuy.Valid())
}

func TestDirection_Valid_Sell(t *testing.T) {
	assert.True(t, derive.DirectionSell.Valid())
}

func TestDirection_Valid_RejectsUnknown(t *testing.T) {
	assert.False(t, derive.Direction("unknown").Valid())
	assert.False(t, derive.Direction("").Valid())
	assert.False(t, derive.Direction("BUY").Valid(), "case-sensitive")
}

func TestDirection_Sign_Buy(t *testing.T) {
	assert.Equal(t, 1, derive.DirectionBuy.Sign())
}

func TestDirection_Sign_Sell(t *testing.T) {
	assert.Equal(t, -1, derive.DirectionSell.Sign())
}

func TestDirection_Sign_PanicsOnInvalid(t *testing.T) {
	assert.Panics(t, func() { _ = derive.Direction("nope").Sign() })
	assert.Panics(t, func() { _ = derive.Direction("").Sign() })
}

func TestDirection_Opposite_BuyToSell(t *testing.T) {
	assert.Equal(t, derive.DirectionSell, derive.DirectionBuy.Opposite())
}

func TestDirection_Opposite_SellToBuy(t *testing.T) {
	assert.Equal(t, derive.DirectionBuy, derive.DirectionSell.Opposite())
}

// Opposite returns Buy for any non-Buy value (the else arm), so an
// unknown value also flips to Buy. We document and verify that contract.
func TestDirection_Opposite_UnknownTreatedAsSell(t *testing.T) {
	assert.Equal(t, derive.DirectionBuy, derive.Direction("unknown").Opposite())
}

func TestDirection_JSONMarshal(t *testing.T) {
	type wrap struct {
		D derive.Direction `json:"d"`
	}
	b, err := json.Marshal(wrap{D: derive.DirectionBuy})
	require.NoError(t, err)
	assert.JSONEq(t, `{"d":"buy"}`, string(b))
}

func TestDirection_JSONUnmarshal(t *testing.T) {
	type wrap struct {
		D derive.Direction `json:"d"`
	}
	var got wrap
	require.NoError(t, json.Unmarshal([]byte(`{"d":"sell"}`), &got))
	assert.Equal(t, derive.DirectionSell, got.D)
}
func TestEnvironment_Valid_Mainnet(t *testing.T) {
	assert.True(t, derive.EnvironmentMainnet.Valid())
}
func TestEnvironment_Valid_Testnet(t *testing.T) {
	assert.True(t, derive.EnvironmentTestnet.Valid())
}

func TestEnvironment_Valid_RejectsUnknown(t *testing.T) {
	for _, v := range []string{"", "prod", "MAINNET", "staging"} {
		assert.False(t, derive.Environment(v).Valid(), "value %q", v)
	}
}

func TestEnvironment_Validate(t *testing.T) {
	assert.NoError(t, derive.EnvironmentMainnet.Validate())
	assert.NoError(t, derive.EnvironmentTestnet.Validate())

	err := derive.Environment("staging").Validate()
	assert.Error(t, err)
	assert.ErrorIs(t, err, derive.ErrInvalidEnum)
	assert.Contains(t, err.Error(), "Environment")
	assert.Contains(t, err.Error(), "staging")
}
func TestInstrumentType_Valid_Perp(t *testing.T) { assert.True(t, derive.InstrumentTypePerp.Valid()) }
func TestInstrumentType_Valid_Option(t *testing.T) {
	assert.True(t, derive.InstrumentTypeOption.Valid())
}
func TestInstrumentType_Valid_ERC20(t *testing.T) { assert.True(t, derive.InstrumentTypeERC20.Valid()) }

func TestInstrumentType_Valid_RejectsUnknown(t *testing.T) {
	for _, v := range []string{"", "future", "spot", "PERP"} {
		assert.False(t, derive.InstrumentType(v).Valid(), "value %q", v)
	}
}
func TestLiquidityRole_Valid_Maker(t *testing.T) { assert.True(t, derive.LiquidityRoleMaker.Valid()) }
func TestLiquidityRole_Valid_Taker(t *testing.T) { assert.True(t, derive.LiquidityRoleTaker.Valid()) }

func TestLiquidityRole_Valid_RejectsUnknown(t *testing.T) {
	for _, v := range []string{"", "MAKER", "neither"} {
		assert.False(t, derive.LiquidityRole(v).Valid(), "value %q", v)
	}
}
func TestMarginType_Valid(t *testing.T) {
	for _, m := range []derive.MarginType{
		derive.MarginTypeSM, derive.MarginTypePM, derive.MarginTypePM2,
	} {
		t.Run(string(m), func(t *testing.T) { assert.True(t, m.Valid()) })
	}
}

func TestMarginType_RejectsUnknown(t *testing.T) {
	assert.False(t, derive.MarginType("").Valid())
	assert.False(t, derive.MarginType("pm3").Valid())
	assert.False(t, derive.MarginType("sm").Valid())
}
func TestOptionType_Valid_Call(t *testing.T) { assert.True(t, derive.OptionTypeCall.Valid()) }
func TestOptionType_Valid_Put(t *testing.T)  { assert.True(t, derive.OptionTypePut.Valid()) }

func TestOptionType_Valid_RejectsUnknown(t *testing.T) {
	for _, v := range []string{"", "straddle", "CALL"} {
		assert.False(t, derive.OptionType(v).Valid(), "value %q", v)
	}
}
func TestOrderStatus_Valid_AllArms(t *testing.T) {
	cases := []derive.OrderStatus{
		derive.OrderStatusOpen,
		derive.OrderStatusFilled,
		derive.OrderStatusCancelled,
		derive.OrderStatusExpired,
		derive.OrderStatusRejected,
	}
	for _, c := range cases {
		t.Run(string(c), func(t *testing.T) {
			assert.True(t, c.Valid())
		})
	}
}

func TestOrderStatus_Valid_RejectsUnknown(t *testing.T) {
	assert.False(t, derive.OrderStatus("").Valid())
	assert.False(t, derive.OrderStatus("haunted").Valid())

	assert.False(t, derive.OrderStatus("untriggered").Valid())
	assert.False(t, derive.OrderStatus("insufficient_margin").Valid())
}

func TestOrderStatus_Terminal_TerminalArms(t *testing.T) {
	terminal := []derive.OrderStatus{
		derive.OrderStatusFilled,
		derive.OrderStatusCancelled,
		derive.OrderStatusExpired,
		derive.OrderStatusRejected,
	}
	for _, s := range terminal {
		t.Run(string(s), func(t *testing.T) {
			assert.True(t, s.Terminal())
		})
	}
}

func TestOrderStatus_Terminal_OpenArm(t *testing.T) {
	assert.False(t, derive.OrderStatusOpen.Terminal())
}

func TestOrderStatus_Terminal_DefaultArm(t *testing.T) {

	assert.False(t, derive.OrderStatus("???").Terminal())
}

func TestOrderStatus_JSONRoundTrip(t *testing.T) {
	type wrap struct {
		S derive.OrderStatus `json:"s"`
	}
	in := wrap{S: derive.OrderStatusFilled}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	assert.JSONEq(t, `{"s":"filled"}`, string(b))

	var out wrap
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in, out)
}
func TestOrderType_Valid_Limit(t *testing.T) {
	assert.True(t, derive.OrderTypeLimit.Valid())
}

func TestOrderType_Valid_Market(t *testing.T) {
	assert.True(t, derive.OrderTypeMarket.Valid())
}

func TestOrderType_Valid_RejectsUnknown(t *testing.T) {
	for _, v := range []string{"", "stop", "LIMIT"} {
		assert.False(t, derive.OrderType(v).Valid(), "value %q", v)
	}
}

func TestOrderType_JSONRoundTrip(t *testing.T) {
	type wrap struct {
		T derive.OrderType `json:"t"`
	}
	in := wrap{T: derive.OrderTypeMarket}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	assert.JSONEq(t, `{"t":"market"}`, string(b))

	var out wrap
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in, out)
}
func TestQuoteStatus_Valid(t *testing.T) {
	for _, q := range []derive.QuoteStatus{
		derive.QuoteStatusOpen, derive.QuoteStatusFilled,
		derive.QuoteStatusCancelled, derive.QuoteStatusExpired,
	} {
		t.Run(string(q), func(t *testing.T) { assert.True(t, q.Valid()) })
	}
}

func TestQuoteStatus_Terminal(t *testing.T) {
	assert.False(t, derive.QuoteStatusOpen.Terminal())
	assert.True(t, derive.QuoteStatusFilled.Terminal())
	assert.True(t, derive.QuoteStatusCancelled.Terminal())
	assert.True(t, derive.QuoteStatusExpired.Terminal())
}

func TestQuoteStatus_RejectsUnknown(t *testing.T) {
	assert.False(t, derive.QuoteStatus("").Valid())
	assert.False(t, derive.QuoteStatus("pending").Valid())
}
func TestTimeInForce_Valid_AllArms(t *testing.T) {
	cases := []derive.TimeInForce{
		derive.TimeInForceGTC,
		derive.TimeInForcePostOnly,
		derive.TimeInForceFOK,
		derive.TimeInForceIOC,
	}
	for _, c := range cases {
		t.Run(string(c), func(t *testing.T) {
			assert.True(t, c.Valid())
		})
	}
}

func TestTimeInForce_Valid_RejectsUnknown(t *testing.T) {
	for _, v := range []string{"", "FOK", "forever", "day"} {
		assert.False(t, derive.TimeInForce(v).Valid(), "value %q", v)
	}
}

func TestTimeInForce_JSONRoundTrip(t *testing.T) {
	type wrap struct {
		T derive.TimeInForce `json:"t"`
	}
	in := wrap{T: derive.TimeInForcePostOnly}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	assert.JSONEq(t, `{"t":"post_only"}`, string(b))

	var out wrap
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in, out)
}
func TestTxStatus_Valid(t *testing.T) {
	for _, s := range []derive.TxStatus{
		derive.TxStatusRequested, derive.TxStatusPending, derive.TxStatusSettled,
		derive.TxStatusReverted, derive.TxStatusIgnored,
	} {
		t.Run(string(s), func(t *testing.T) { assert.True(t, s.Valid()) })
	}
}

func TestTxStatus_Terminal(t *testing.T) {
	assert.False(t, derive.TxStatusRequested.Terminal())
	assert.False(t, derive.TxStatusPending.Terminal())
	assert.True(t, derive.TxStatusSettled.Terminal())
	assert.True(t, derive.TxStatusReverted.Terminal())
	assert.True(t, derive.TxStatusIgnored.Terminal())
}

func TestTxStatus_RejectsUnknown(t *testing.T) {
	assert.False(t, derive.TxStatus("").Valid())
	assert.False(t, derive.TxStatus("done").Valid())
}

// validator is the common interface every enum's Validate method satisfies.
type validator interface{ Validate() error }

// TestEnums_Validate_AcceptsValid runs Validate on a known-good value
// for every enum. All must return nil.
func TestEnums_Validate_AcceptsValid(t *testing.T) {
	cases := []struct {
		name string
		v    validator
	}{
		{"AssetType", derive.AssetTypeERC20},
		{"AuctionType", derive.AuctionTypeSolvent},
		{"BalanceUpdateType", derive.BalanceUpdateTrade},
		{"CancelReason", derive.CancelReasonUserRequest},
		{"Direction", derive.DirectionBuy},
		{"InstrumentType", derive.InstrumentTypePerp},
		{"LiquidityRole", derive.LiquidityRoleMaker},
		{"MarginType", derive.MarginTypePM},
		{"OptionType", derive.OptionTypeCall},
		{"OrderStatus", derive.OrderStatusOpen},
		{"OrderType", derive.OrderTypeLimit},
		{"QuoteStatus", derive.QuoteStatusOpen},
		{"TimeInForce", derive.TimeInForceGTC},
		{"TxStatus", derive.TxStatusPending},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.NoError(t, c.v.Validate())
		})
	}
}

// TestEnums_Validate_RejectsUnknown runs Validate on a value the enum
// would not accept and confirms the result wraps ErrInvalidEnum.
func TestEnums_Validate_RejectsUnknown(t *testing.T) {
	cases := []struct {
		name string
		v    validator
	}{
		{"AssetType", derive.AssetType("future")},
		{"AuctionType", derive.AuctionType("partial")},
		{"BalanceUpdateType", derive.BalanceUpdateType("magic")},
		{"CancelReason", derive.CancelReason("self_cross")},
		{"Direction", derive.Direction("up")},
		{"InstrumentType", derive.InstrumentType("future")},
		{"LiquidityRole", derive.LiquidityRole("middleman")},
		{"MarginType", derive.MarginType("pm3")},
		{"OptionType", derive.OptionType("call")},
		{"OrderStatus", derive.OrderStatus("untriggered")},
		{"OrderType", derive.OrderType("stop")},
		{"QuoteStatus", derive.QuoteStatus("pending")},
		{"TimeInForce", derive.TimeInForce("immediate")},
		{"TxStatus", derive.TxStatus("done")},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.v.Validate()
			assert.Error(t, err)
			assert.True(t, errors.Is(err, derive.ErrInvalidEnum),
				"expected wrap of ErrInvalidEnum, got %v", err)
		})
	}
}

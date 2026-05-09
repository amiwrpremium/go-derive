package enums_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/amiwrpremium/go-derive/pkg/enums"
)

// validator is the common interface every enum's Validate method satisfies.
type validator interface{ Validate() error }

// TestEnums_Validate_AcceptsValid runs Validate on a known-good value
// for every enum. All must return nil.
func TestEnums_Validate_AcceptsValid(t *testing.T) {
	cases := []struct {
		name string
		v    validator
	}{
		{"AssetType", enums.AssetTypeERC20},
		{"AuctionType", enums.AuctionTypeSolvent},
		{"BalanceUpdateType", enums.BalanceUpdateTrade},
		{"CancelReason", enums.CancelReasonUserRequest},
		{"Direction", enums.DirectionBuy},
		{"InstrumentType", enums.InstrumentTypePerp},
		{"LiquidityRole", enums.LiquidityRoleMaker},
		{"MarginType", enums.MarginTypePM},
		{"OptionType", enums.OptionTypeCall},
		{"OrderStatus", enums.OrderStatusOpen},
		{"OrderType", enums.OrderTypeLimit},
		{"QuoteStatus", enums.QuoteStatusOpen},
		{"TimeInForce", enums.TimeInForceGTC},
		{"TxStatus", enums.TxStatusPending},
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
		{"AssetType", enums.AssetType("future")},
		{"AuctionType", enums.AuctionType("partial")},
		{"BalanceUpdateType", enums.BalanceUpdateType("magic")},
		{"CancelReason", enums.CancelReason("self_cross")},
		{"Direction", enums.Direction("up")},
		{"InstrumentType", enums.InstrumentType("future")},
		{"LiquidityRole", enums.LiquidityRole("middleman")},
		{"MarginType", enums.MarginType("pm3")},
		{"OptionType", enums.OptionType("call")}, // legacy wire value
		{"OrderStatus", enums.OrderStatus("haunted")},
		{"OrderType", enums.OrderType("stop")},
		{"QuoteStatus", enums.QuoteStatus("pending")},
		{"TimeInForce", enums.TimeInForce("immediate")},
		{"TxStatus", enums.TxStatus("done")},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.v.Validate()
			assert.Error(t, err)
			assert.True(t, errors.Is(err, enums.ErrInvalidEnum),
				"expected wrap of ErrInvalidEnum, got %v", err)
		})
	}
}

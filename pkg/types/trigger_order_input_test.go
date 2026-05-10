package types_test

import (
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func validTriggerOrderInput() types.TriggerOrderInput {
	return types.TriggerOrderInput{
		PlaceOrderInput: types.PlaceOrderInput{
			InstrumentName: "BTC-PERP",
			Asset:          types.Address(common.HexToAddress("0x1111111111111111111111111111111111111111")),
			Direction:      enums.DirectionSell,
			OrderType:      enums.OrderTypeLimit,
			Amount:         types.MustDecimal("1"),
			LimitPrice:     types.MustDecimal("60000"),
			MaxFee:         types.MustDecimal("1"),
		},
		TriggerType:      enums.TriggerTypeStopLoss,
		TriggerPriceType: enums.TriggerPriceTypeMark,
		TriggerPrice:     types.MustDecimal("59000"),
	}
}

func TestTriggerOrderInput_Validate_Happy(t *testing.T) {
	require.NoError(t, validTriggerOrderInput().Validate())
}

func TestTriggerOrderInput_Validate_Rejects(t *testing.T) {
	cases := []struct {
		name string
		mut  func(*types.TriggerOrderInput)
		want string
	}{
		{"bad trigger_type", func(in *types.TriggerOrderInput) { in.TriggerType = enums.TriggerType("trail") }, "trigger_type"},
		{"bad price_type", func(in *types.TriggerOrderInput) { in.TriggerPriceType = enums.TriggerPriceType("last") }, "trigger_price_type"},
		{"zero trigger_price", func(in *types.TriggerOrderInput) { in.TriggerPrice = types.MustDecimal("0") }, "trigger_price"},
		{"underlying invalid", func(in *types.TriggerOrderInput) { in.InstrumentName = "" }, "instrument_name"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			in := validTriggerOrderInput()
			c.mut(&in)
			err := in.Validate()
			require.Error(t, err)
			assert.True(t, errors.Is(err, types.ErrInvalidParams))
			assert.Contains(t, err.Error(), c.want)
		})
	}
}

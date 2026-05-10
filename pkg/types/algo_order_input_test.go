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

func validAlgoOrderInput() types.AlgoOrderInput {
	return types.AlgoOrderInput{
		PlaceOrderInput: types.PlaceOrderInput{
			InstrumentName: "BTC-PERP",
			Asset:          types.Address(common.HexToAddress("0x1111111111111111111111111111111111111111")),
			Direction:      enums.DirectionBuy,
			OrderType:      enums.OrderTypeLimit,
			Amount:         types.MustDecimal("1"),
			LimitPrice:     types.MustDecimal("100"),
			MaxFee:         types.MustDecimal("1"),
		},
		AlgoType:        enums.AlgoTypeTWAP,
		AlgoDurationSec: 3600,
		AlgoNumSlices:   60,
	}
}

func TestAlgoOrderInput_Validate_Happy(t *testing.T) {
	require.NoError(t, validAlgoOrderInput().Validate())
}

func TestAlgoOrderInput_Validate_Rejects(t *testing.T) {
	cases := []struct {
		name string
		mut  func(*types.AlgoOrderInput)
		want string
	}{
		{"bad algo_type", func(in *types.AlgoOrderInput) { in.AlgoType = enums.AlgoType("vwap") }, "algo_type"},
		{"zero duration", func(in *types.AlgoOrderInput) { in.AlgoDurationSec = 0 }, "algo_duration_sec"},
		{"zero slices", func(in *types.AlgoOrderInput) { in.AlgoNumSlices = 0 }, "algo_num_slices"},
		{"underlying invalid", func(in *types.AlgoOrderInput) { in.InstrumentName = "" }, "instrument_name"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			in := validAlgoOrderInput()
			c.mut(&in)
			err := in.Validate()
			require.Error(t, err)
			assert.True(t, errors.Is(err, types.ErrInvalidParams))
			assert.Contains(t, err.Error(), c.want)
		})
	}
}

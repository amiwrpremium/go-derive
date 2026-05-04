package types_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestOrderParams_Builder_Chains(t *testing.T) {
	got := types.NewOrderParams("BTC-PERP", enums.DirectionBuy, enums.OrderTypeLimit,
		types.MustDecimal("0.5"), types.MustDecimal("65000")).
		WithMaxFee(types.MustDecimal("10")).
		WithLabel("scalp-1").
		WithSubaccount(123).
		WithTimeInForce(enums.TimeInForcePostOnly).
		WithMMP().
		WithReduceOnly()

	assert.Equal(t, "BTC-PERP", got.InstrumentName)
	assert.Equal(t, enums.DirectionBuy, got.Direction)
	assert.Equal(t, enums.OrderTypeLimit, got.OrderType)
	assert.Equal(t, "0.5", got.Amount.String())
	assert.Equal(t, "65000", got.LimitPrice.String())
	assert.Equal(t, "10", got.MaxFee.String())
	assert.Equal(t, "scalp-1", got.Label)
	assert.Equal(t, int64(123), got.SubaccountID)
	assert.Equal(t, enums.TimeInForcePostOnly, got.TimeInForce)
	assert.True(t, got.MMP)
	assert.True(t, got.ReduceOnly)
}

func TestOrderParams_Builder_DoesNotMutateOriginal(t *testing.T) {
	base := types.NewOrderParams("BTC-PERP", enums.DirectionBuy, enums.OrderTypeLimit,
		types.MustDecimal("1"), types.MustDecimal("100"))
	_ = base.WithLabel("a")
	assert.Empty(t, base.Label, "WithLabel must return a copy, not mutate base")
}

func TestOrderParams_Validate_OK(t *testing.T) {
	p := types.NewOrderParams("BTC-PERP", enums.DirectionBuy, enums.OrderTypeLimit,
		types.MustDecimal("1"), types.MustDecimal("100"))
	require.NoError(t, p.Validate())
}

func TestOrderParams_Validate_RejectsBlankInstrument(t *testing.T) {
	p := types.NewOrderParams("", enums.DirectionBuy, enums.OrderTypeLimit,
		types.MustDecimal("1"), types.MustDecimal("100"))
	err := p.Validate()
	assert.ErrorIs(t, err, types.ErrInvalidParams)
	assert.Contains(t, err.Error(), "instrument_name")
}

func TestOrderParams_Validate_RejectsBadEnum(t *testing.T) {
	p := types.NewOrderParams("BTC-PERP", enums.Direction("up"), enums.OrderTypeLimit,
		types.MustDecimal("1"), types.MustDecimal("100"))
	err := p.Validate()
	assert.ErrorIs(t, err, types.ErrInvalidParams)
	assert.Contains(t, err.Error(), "direction")
}

func TestOrderParams_Validate_RejectsZeroAmount(t *testing.T) {
	p := types.NewOrderParams("BTC-PERP", enums.DirectionBuy, enums.OrderTypeLimit,
		types.MustDecimal("0"), types.MustDecimal("100"))
	err := p.Validate()
	assert.ErrorIs(t, err, types.ErrInvalidParams)
	assert.Contains(t, err.Error(), "amount")
}

func TestOrderParams_Validate_RejectsNegativePrice(t *testing.T) {
	p := types.NewOrderParams("BTC-PERP", enums.DirectionBuy, enums.OrderTypeLimit,
		types.MustDecimal("1"), types.MustDecimal("-1"))
	err := p.Validate()
	assert.ErrorIs(t, err, types.ErrInvalidParams)
	assert.Contains(t, err.Error(), "limit_price")
}

func TestOrderParams_Validate_RejectsNegativeMaxFee(t *testing.T) {
	p := types.NewOrderParams("BTC-PERP", enums.DirectionBuy, enums.OrderTypeLimit,
		types.MustDecimal("1"), types.MustDecimal("100")).
		WithMaxFee(types.MustDecimal("-0.1"))
	err := p.Validate()
	assert.ErrorIs(t, err, types.ErrInvalidParams)
	assert.Contains(t, err.Error(), "max_fee")
}

func TestCancelOrderParams_Builder_AndValidate(t *testing.T) {
	got := types.NewCancelOrderParams(7).WithOrderID("O1").WithInstrument("BTC-PERP")
	require.NoError(t, got.Validate())
	assert.Equal(t, "O1", got.OrderID)
	assert.Equal(t, "BTC-PERP", got.InstrumentName)

	got2 := types.NewCancelOrderParams(7).WithLabel("scalp")
	require.NoError(t, got2.Validate())

	noTarget := types.NewCancelOrderParams(7)
	err := noTarget.Validate()
	assert.ErrorIs(t, err, types.ErrInvalidParams)

	bad := types.NewCancelOrderParams(-1).WithOrderID("O1")
	assert.Error(t, bad.Validate())
}

func TestReplaceOrderParams_Builder_AndValidate(t *testing.T) {
	o := types.NewOrderParams("BTC-PERP", enums.DirectionBuy, enums.OrderTypeLimit,
		types.MustDecimal("1"), types.MustDecimal("100"))
	r := types.NewReplaceOrderParams("OLD", o)
	require.NoError(t, r.Validate())
	assert.Equal(t, "OLD", r.OrderIDToCancel)

	bad := types.NewReplaceOrderParams("", o)
	assert.ErrorIs(t, bad.Validate(), types.ErrInvalidParams)

	badInner := types.NewReplaceOrderParams("OLD",
		types.NewOrderParams("", enums.DirectionBuy, enums.OrderTypeLimit,
			types.MustDecimal("1"), types.MustDecimal("100")))
	err := badInner.Validate()
	assert.True(t, errors.Is(err, types.ErrInvalidParams))
}

func TestPageRequest_Builder_AndValidate(t *testing.T) {
	p := types.NewPageRequest().WithPage(2).WithPageSize(50)
	require.NoError(t, p.Validate())
	assert.Equal(t, 2, p.Page)
	assert.Equal(t, 50, p.PageSize)

	require.NoError(t, types.NewPageRequest().Validate()) // zeros are server defaults
	assert.ErrorIs(t, types.NewPageRequest().WithPage(-1).Validate(), types.ErrInvalidParams)
	assert.ErrorIs(t, types.NewPageRequest().WithPageSize(-1).Validate(), types.ErrInvalidParams)
}

// TestOrderParams_Builder_AllSetters exercises every With* setter at least
// once, verifying the value is applied and the original is not mutated.
func TestOrderParams_Builder_AllSetters(t *testing.T) {
	base := types.NewOrderParams("X", enums.DirectionBuy, enums.OrderTypeLimit,
		types.MustDecimal("1"), types.MustDecimal("100"))

	a := base.WithInstrument("BTC-PERP")
	assert.Equal(t, "BTC-PERP", a.InstrumentName)

	b := base.WithDirection(enums.DirectionSell)
	assert.Equal(t, enums.DirectionSell, b.Direction)

	c := base.WithOrderType(enums.OrderTypeMarket)
	assert.Equal(t, enums.OrderTypeMarket, c.OrderType)

	d := base.WithAmount(types.MustDecimal("2"))
	assert.Equal(t, "2", d.Amount.String())

	e := base.WithLimitPrice(types.MustDecimal("200"))
	assert.Equal(t, "200", e.LimitPrice.String())

	addr := types.MustAddress("0x1111111111111111111111111111111111111111")
	f := base.WithSignature(addr, "0xsig", 99, 1700000000)
	assert.Equal(t, types.MustAddress("0x1111111111111111111111111111111111111111"), f.Signer)
	assert.Equal(t, "0xsig", f.Signature)
	assert.Equal(t, uint64(99), f.Nonce)
	assert.Equal(t, int64(1700000000), f.SignatureExpiry)

	assert.Equal(t, "X", base.InstrumentName) // original unmutated
}

func TestCancelOrderParams_WithSignature(t *testing.T) {
	addr := types.MustAddress("0x1111111111111111111111111111111111111111")
	p := types.NewCancelOrderParams(7).WithOrderID("O1").
		WithSignature(addr, "0xsig", 11, 1700000000)
	assert.Equal(t, types.MustAddress("0x1111111111111111111111111111111111111111"), p.Signer)
	assert.Equal(t, "0xsig", p.Signature)
	assert.Equal(t, uint64(11), p.Nonce)
	assert.Equal(t, int64(1700000000), p.SignatureExpiry)
}

func TestReplaceOrderParams_AllSetters(t *testing.T) {
	o := types.NewOrderParams("BTC-PERP", enums.DirectionBuy, enums.OrderTypeLimit,
		types.MustDecimal("1"), types.MustDecimal("100"))
	r := types.NewReplaceOrderParams("X", o)

	r2 := r.WithOrderIDToCancel("OLD")
	assert.Equal(t, "OLD", r2.OrderIDToCancel)
	assert.Equal(t, "X", r.OrderIDToCancel)

	o2 := types.NewOrderParams("ETH-PERP", enums.DirectionSell, enums.OrderTypeMarket,
		types.MustDecimal("2"), types.MustDecimal("3500"))
	r3 := r.WithNewOrder(o2)
	assert.Equal(t, "ETH-PERP", r3.NewOrder.InstrumentName)
}

// TestOrderParams_Validate_RejectsRemainingBranches covers branches the
// existing happy-path tests don't reach: bad order_type, bad time_in_force,
// negative subaccount, negative signature_expiry.
func TestOrderParams_Validate_RejectsRemainingBranches(t *testing.T) {
	cases := []struct {
		name string
		mut  func(*types.OrderParams)
		want string
	}{
		{"bad order_type", func(p *types.OrderParams) { p.OrderType = enums.OrderType("nope") }, "order_type"},
		{"bad tif", func(p *types.OrderParams) { p.TimeInForce = enums.TimeInForce("forever") }, "time_in_force"},
		{"negative subaccount", func(p *types.OrderParams) { p.SubaccountID = -1 }, "subaccount_id"},
		{"negative expiry", func(p *types.OrderParams) { p.SignatureExpiry = -1 }, "signature_expiry"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := types.NewOrderParams("BTC-PERP", enums.DirectionBuy, enums.OrderTypeLimit,
				types.MustDecimal("1"), types.MustDecimal("100"))
			c.mut(&p)
			err := p.Validate()
			require.Error(t, err)
			assert.True(t, errors.Is(err, types.ErrInvalidParams))
			assert.Contains(t, err.Error(), c.want)
		})
	}
}

func TestCancelOrderParams_Validate_RejectsBadExpiry(t *testing.T) {
	p := types.NewCancelOrderParams(0).WithOrderID("O1")
	p.SignatureExpiry = -1
	assert.ErrorIs(t, p.Validate(), types.ErrInvalidParams)
}

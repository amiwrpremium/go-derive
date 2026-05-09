package derive_test

import (
	"encoding/json"
	"errors"
	"github.com/amiwrpremium/go-derive"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// FuzzNewAddress verifies the parser is panic-free for arbitrary input.
func FuzzNewAddress(f *testing.F) {
	f.Add("0x1111111111111111111111111111111111111111")
	f.Add("0xZZZZ")
	f.Add("not-an-address")
	f.Add("")
	f.Add("0x")
	f.Add(string(make([]byte, 1024)))

	f.Fuzz(func(t *testing.T, s string) {
		a, err := derive.NewAddress(s)
		if err != nil && !a.IsZero() {
			t.Fatalf("error path leaked non-zero address for %q", s)
		}
	})
}

// FuzzNewTxHash verifies the parser is panic-free for arbitrary input.
func FuzzNewTxHash(f *testing.F) {
	f.Add("0x1111111111111111111111111111111111111111111111111111111111111111")
	f.Add("0xabc")
	f.Add("0x")
	f.Add("")
	f.Add("not-a-hash")

	f.Fuzz(func(t *testing.T, s string) {
		h, err := derive.NewTxHash(s)
		if err != nil && !h.IsZero() {
			t.Fatalf("error path leaked non-zero hash for %q", s)
		}
	})
}

const (
	addrUpper = "0x1234567890ABCDEF1234567890ABCDEF12345678"
	addrLower = "0x1234567890abcdef1234567890abcdef12345678"
)

func TestNewAddress_Valid(t *testing.T) {
	a, err := derive.NewAddress(addrLower)
	require.NoError(t, err)
	assert.False(t, a.IsZero())
	assert.Equal(t, common.HexToAddress(addrLower), a.Common())
}

func TestNewAddress_Empty(t *testing.T) {
	a, err := derive.NewAddress("")
	require.NoError(t, err)
	assert.True(t, a.IsZero())
}

func TestNewAddress_Invalid(t *testing.T) {
	_, err := derive.NewAddress("not-an-address")
	assert.Error(t, err)
}

func TestMustAddress_Panics(t *testing.T) {
	assert.Panics(t, func() { derive.MustAddress("nope") })
	assert.NotPanics(t, func() { derive.MustAddress(addrLower) })
}

func TestAddress_StringIsChecksummed(t *testing.T) {
	a := derive.MustAddress(addrLower)

	assert.Equal(t, common.HexToAddress(addrLower).Hex(), a.String())
}

func TestAddress_JSONRoundTrip(t *testing.T) {
	type wrap struct {
		A derive.Address `json:"a"`
	}
	in := wrap{A: derive.MustAddress(addrLower)}
	b, err := json.Marshal(in)
	require.NoError(t, err)

	var out wrap
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.A, out.A)
}

func TestAddress_UnmarshalEmpty(t *testing.T) {
	var a derive.Address
	require.NoError(t, json.Unmarshal([]byte(`""`), &a))
	assert.True(t, a.IsZero())
}

func TestAddress_UnmarshalInvalid(t *testing.T) {
	var a derive.Address
	err := json.Unmarshal([]byte(`"not-an-address"`), &a)
	assert.Error(t, err)
}

func TestAddress_UnmarshalNonString(t *testing.T) {
	var a derive.Address
	err := json.Unmarshal([]byte(`123`), &a)
	assert.Error(t, err)
}
func TestCollateral_Decode(t *testing.T) {
	payload := `{
		"asset_name": "USDC",
		"asset_type": "erc20",
		"amount": "10000",
		"mark_price": "1",
		"mark_value": "10000",
		"cumulative_interest": "5",
		"pending_interest": "0.1",
		"initial_margin": "100",
		"maintenance_margin": "50"
	}`
	var c derive.Collateral
	require.NoError(t, json.Unmarshal([]byte(payload), &c))
	assert.Equal(t, "USDC", c.AssetName)
	assert.Equal(t, derive.AssetTypeERC20, c.AssetType)
	assert.Equal(t, "10000", c.Amount.String())
}

func TestCollateral_RoundTrip(t *testing.T) {
	in := derive.Collateral{
		AssetName: "USDC", AssetType: derive.AssetTypeERC20,
		Amount:    derive.MustDecimal("100"),
		MarkValue: derive.MustDecimal("100"),
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	var out derive.Collateral
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.AssetName, out.AssetName)
	assert.Equal(t, in.Amount.String(), out.Amount.String())
}

func TestBalance_Decode(t *testing.T) {
	payload := `{
		"subaccount_id": 123,
		"subaccount_value": "10000",
		"initial_margin": "5000",
		"maintenance_margin": "3000",
		"collaterals": [{"asset_name": "USDC", "asset_type": "erc20", "amount": "10000", "mark_value": "10000"}],
		"positions": []
	}`
	var b derive.Balance
	require.NoError(t, json.Unmarshal([]byte(payload), &b))
	assert.Equal(t, int64(123), b.SubaccountID)
	require.Len(t, b.Collaterals, 1)
	assert.Equal(t, "USDC", b.Collaterals[0].AssetName)
	assert.Empty(t, b.Positions)
}

func TestBalance_OmitsEmptyPositionsOnMarshal(t *testing.T) {
	in := derive.Balance{
		SubaccountID:      1,
		SubaccountValue:   derive.MustDecimal("0"),
		InitialMargin:     derive.MustDecimal("0"),
		MaintenanceMargin: derive.MustDecimal("0"),
		Collaterals:       []derive.Collateral{},
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	assert.NotContains(t, string(b), "positions")
}
func TestCandle_Decode(t *testing.T) {
	payload := `{
		"timestamp": 1700000000000,
		"open": "100",
		"high": "110",
		"low": "95",
		"close": "105",
		"volume": "10000"
	}`
	var c derive.Candle
	require.NoError(t, json.Unmarshal([]byte(payload), &c))
	assert.Equal(t, "100", c.Open.String())
	assert.Equal(t, "110", c.High.String())
	assert.Equal(t, "95", c.Low.String())
	assert.Equal(t, "105", c.Close.String())
}

// Decimal's custom MarshalJSON defeats `omitempty` — a zero decimal still
// serialises as the string "0". Document and pin the actual behaviour.
func TestCandle_ZeroVolumeStillSerialized(t *testing.T) {
	in := derive.Candle{
		Open:  derive.MustDecimal("1"),
		High:  derive.MustDecimal("1"),
		Low:   derive.MustDecimal("1"),
		Close: derive.MustDecimal("1"),
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	assert.Contains(t, string(b), `"volume":"0"`)
}

// FuzzNewDecimal verifies the parser never panics on arbitrary input. Bad
// input must produce a (zero, error); good input must round-trip.
func FuzzNewDecimal(f *testing.F) {
	f.Add("0")
	f.Add("1.5")
	f.Add("-2.5")
	f.Add("0.000000000000000001")
	f.Add("100000000000000000000")
	f.Add("1e10")
	f.Add("not-a-number")
	f.Add("")
	f.Add("0x10")
	f.Add(string([]byte{0xff, 0xfe}))

	f.Fuzz(func(t *testing.T, s string) {
		d, err := derive.NewDecimal(s)
		if err != nil {

			if !d.IsZero() {
				t.Fatalf("error path returned non-zero decimal: %q -> %s", s, d)
			}
			return
		}

		b, err := json.Marshal(d)
		if err != nil {
			t.Fatalf("marshal succeeded-parse decimal: %v", err)
		}
		var back derive.Decimal
		if err := json.Unmarshal(b, &back); err != nil {
			t.Fatalf("round-trip unmarshal: %v", err)
		}
	})
}

// FuzzDecimal_UnmarshalJSON checks the JSON unmarshaler doesn't panic.
func FuzzDecimal_UnmarshalJSON(f *testing.F) {
	f.Add([]byte(`"1.5"`))
	f.Add([]byte(`1.5`))
	f.Add([]byte(`null`))
	f.Add([]byte(`""`))
	f.Add([]byte(`"-0.001"`))
	f.Add([]byte(`{`))
	f.Add([]byte(``))
	f.Add([]byte{0xff, 0xfe})

	f.Fuzz(func(t *testing.T, raw []byte) {
		var d derive.Decimal
		_ = d.UnmarshalJSON(raw)
	})
}
func TestNewDecimal_Valid(t *testing.T) {
	d, err := derive.NewDecimal("1.234")
	require.NoError(t, err)
	assert.Equal(t, "1.234", d.String())
}

func TestNewDecimal_Invalid(t *testing.T) {
	_, err := derive.NewDecimal("not-a-number")
	assert.Error(t, err)
}

func TestMustDecimal_Panics(t *testing.T) {
	assert.Panics(t, func() { derive.MustDecimal("not-a-number") })
	assert.NotPanics(t, func() { derive.MustDecimal("0") })
}

func TestDecimalFromInt(t *testing.T) {
	d := derive.DecimalFromInt(42)
	assert.Equal(t, "42", d.String())
	assert.False(t, d.IsZero())
	assert.Equal(t, 1, d.Sign())

	z := derive.DecimalFromInt(0)
	assert.True(t, z.IsZero())
	assert.Equal(t, 0, z.Sign())

	n := derive.DecimalFromInt(-5)
	assert.Equal(t, -1, n.Sign())
}

func TestDecimal_JSONRoundTrip(t *testing.T) {
	type wrap struct {
		Price derive.Decimal `json:"price"`
	}
	in := wrap{Price: derive.MustDecimal("65000.5")}

	b, err := json.Marshal(in)
	require.NoError(t, err)
	assert.JSONEq(t, `{"price":"65000.5"}`, string(b))

	var out wrap
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.Price.String(), out.Price.String())
}

func TestDecimal_UnmarshalAcceptsNumber(t *testing.T) {
	var d derive.Decimal
	require.NoError(t, json.Unmarshal([]byte(`123.45`), &d))
	assert.Equal(t, "123.45", d.String())
}

func TestDecimal_UnmarshalNullAndEmpty(t *testing.T) {
	var d derive.Decimal
	require.NoError(t, json.Unmarshal([]byte(`null`), &d))
	assert.True(t, d.IsZero())

	d2 := derive.MustDecimal("1")
	require.NoError(t, json.Unmarshal([]byte(`""`), &d2))

	assert.Equal(t, "1", d2.String())
}

func TestDecimal_UnmarshalMalformedString(t *testing.T) {
	var d derive.Decimal
	err := json.Unmarshal([]byte(`"abc"`), &d)
	assert.Error(t, err)
}

func TestDecimal_UnmarshalMalformedNumber(t *testing.T) {
	var d derive.Decimal
	err := json.Unmarshal([]byte(`{`), &d)
	assert.Error(t, err)
}

func TestDecimal_Inner(t *testing.T) {
	d := derive.MustDecimal("5.5")
	inner := d.Inner()
	assert.Equal(t, "5.5", inner.String())
}
func TestInstrument_DecodePerp(t *testing.T) {
	payload := `{
		"instrument_name": "BTC-PERP",
		"base_currency": "BTC",
		"quote_currency": "USDC",
		"instrument_type": "perp",
		"is_active": true,
		"tick_size": "0.5",
		"minimum_amount": "0.001",
		"maximum_amount": "1000",
		"amount_step": "0.001",
		"mark_price": "65000.5",
		"index_price": "65000",
		"perp_details": {"index": "BTC", "max_leverage": "50"}
	}`
	var inst derive.Instrument
	require.NoError(t, json.Unmarshal([]byte(payload), &inst))
	assert.Equal(t, derive.InstrumentTypePerp, inst.Type)
	require.NotNil(t, inst.Perp)
	assert.Equal(t, "BTC", inst.Perp.IndexName)
	assert.Equal(t, "50", inst.Perp.MaxLeverage.String())
	assert.Nil(t, inst.Option)
	assert.Nil(t, inst.ERC20)
}

func TestInstrument_DecodeOption(t *testing.T) {
	payload := `{
		"instrument_name": "BTC-25DEC25-65000-C",
		"instrument_type": "option",
		"is_active": true,
		"tick_size": "0.5",
		"minimum_amount": "0.01",
		"maximum_amount": "100",
		"amount_step": "0.01",
		"option_details": {"option_type": "C", "strike": "65000", "expiry": 1735689600000, "index": "BTC"}
	}`
	var inst derive.Instrument
	require.NoError(t, json.Unmarshal([]byte(payload), &inst))
	require.NotNil(t, inst.Option)
	assert.Equal(t, derive.OptionTypeCall, inst.Option.OptionType)
	assert.Equal(t, "65000", inst.Option.Strike.String())
}

func TestInstrument_DecodeERC20(t *testing.T) {
	payload := `{
		"instrument_name": "USDC",
		"instrument_type": "erc20",
		"is_active": true,
		"tick_size": "0.01",
		"minimum_amount": "1",
		"maximum_amount": "1000000",
		"amount_step": "1",
		"erc20_details": {
			"underlying_erc20_address": "0x1111111111111111111111111111111111111111",
			"borrow_index": "1.0",
			"supply_index": "1.0"
		}
	}`
	var inst derive.Instrument
	require.NoError(t, json.Unmarshal([]byte(payload), &inst))
	require.NotNil(t, inst.ERC20)

	assert.Equal(t, "1", inst.ERC20.BorrowIndex.String())
}

func TestInstrument_RoundTrip(t *testing.T) {
	in := derive.Instrument{
		Name:          "BTC-PERP",
		BaseCurrency:  "BTC",
		QuoteCurrency: "USDC",
		Type:          derive.InstrumentTypePerp,
		IsActive:      true,
		TickSize:      derive.MustDecimal("0.5"),
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	var out derive.Instrument
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.Name, out.Name)
	assert.Equal(t, in.Type, out.Type)
}
func TestLiquidation_Decode(t *testing.T) {
	payload := `{
		"subaccount_id": 9,
		"timestamp": 1700000000000,
		"tx_hash": "0x1111111111111111111111111111111111111111111111111111111111111111"
	}`
	var l derive.Liquidation
	require.NoError(t, json.Unmarshal([]byte(payload), &l))
	assert.Equal(t, int64(9), l.SubaccountID)
	assert.False(t, l.TxHash.IsZero())
}

// Custom MarshalJSON on TxHash defeats omitempty (Go's json package only
// omits the natural zero value for built-in types). The wire format always
// includes the hash field even when zero — document and pin that.
func TestLiquidation_ZeroTxHashSerializedAsZeroString(t *testing.T) {
	in := derive.Liquidation{SubaccountID: 1}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	assert.Contains(t, string(b), `"tx_hash":"0x0000000000000000000000000000000000000000000000000000000000000000"`)
}

// FuzzOrderBookLevel_UnmarshalJSON guards the array-shaped level decoder
// against panics on adversarial input.
func FuzzOrderBookLevel_UnmarshalJSON(f *testing.F) {
	f.Add([]byte(`["100","1"]`))
	f.Add([]byte(`["",""]`))
	f.Add([]byte(`[]`))
	f.Add([]byte(`[1,2,3]`))
	f.Add([]byte(`{"price":"1","amount":"2"}`))
	f.Add([]byte(`null`))
	f.Add([]byte(``))

	f.Fuzz(func(t *testing.T, raw []byte) {
		var l derive.OrderBookLevel
		_ = l.UnmarshalJSON(raw)
	})
}
func TestOrderBookLevel_JSONRoundTrip(t *testing.T) {
	in := derive.OrderBookLevel{
		Price:  derive.MustDecimal("100.5"),
		Amount: derive.MustDecimal("0.25"),
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	assert.JSONEq(t, `["100.5","0.25"]`, string(b))

	var out derive.OrderBookLevel
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.Price.String(), out.Price.String())
	assert.Equal(t, in.Amount.String(), out.Amount.String())
}

func TestOrderBookLevel_UnmarshalRejectsObject(t *testing.T) {
	var l derive.OrderBookLevel
	err := json.Unmarshal([]byte(`{"price":"1","amount":"2"}`), &l)
	assert.Error(t, err)
}

func TestOrderBook_DecodeFullPayload(t *testing.T) {
	payload := `{
		"instrument_name": "BTC-PERP",
		"bids": [["65000","0.1"],["64999","0.2"]],
		"asks": [["65001","0.3"]],
		"timestamp": 1700000000000,
		"publish_time": 1700000000005
	}`
	var ob derive.OrderBook
	require.NoError(t, json.Unmarshal([]byte(payload), &ob))
	assert.Equal(t, "BTC-PERP", ob.InstrumentName)
	require.Len(t, ob.Bids, 2)
	require.Len(t, ob.Asks, 1)
	assert.Equal(t, "65000", ob.Bids[0].Price.String())
	assert.Equal(t, "0.1", ob.Bids[0].Amount.String())
	assert.Equal(t, "65001", ob.Asks[0].Price.String())
	assert.Equal(t, int64(1700000000000), ob.Timestamp.Millis())
}
func TestOrderParams_Builder_Chains(t *testing.T) {
	got := derive.NewOrderParams("BTC-PERP", derive.DirectionBuy, derive.OrderTypeLimit,
		derive.MustDecimal("0.5"), derive.MustDecimal("65000")).
		WithMaxFee(derive.MustDecimal("10")).
		WithLabel("scalp-1").
		WithSubaccount(123).
		WithTimeInForce(derive.TimeInForcePostOnly).
		WithMMP().
		WithReduceOnly()

	assert.Equal(t, "BTC-PERP", got.InstrumentName)
	assert.Equal(t, derive.DirectionBuy, got.Direction)
	assert.Equal(t, derive.OrderTypeLimit, got.OrderType)
	assert.Equal(t, "0.5", got.Amount.String())
	assert.Equal(t, "65000", got.LimitPrice.String())
	assert.Equal(t, "10", got.MaxFee.String())
	assert.Equal(t, "scalp-1", got.Label)
	assert.Equal(t, int64(123), got.SubaccountID)
	assert.Equal(t, derive.TimeInForcePostOnly, got.TimeInForce)
	assert.True(t, got.MMP)
	assert.True(t, got.ReduceOnly)
}

func TestOrderParams_Builder_DoesNotMutateOriginal(t *testing.T) {
	base := derive.NewOrderParams("BTC-PERP", derive.DirectionBuy, derive.OrderTypeLimit,
		derive.MustDecimal("1"), derive.MustDecimal("100"))
	_ = base.WithLabel("a")
	assert.Empty(t, base.Label, "WithLabel must return a copy, not mutate base")
}

func TestOrderParams_Validate_OK(t *testing.T) {
	p := derive.NewOrderParams("BTC-PERP", derive.DirectionBuy, derive.OrderTypeLimit,
		derive.MustDecimal("1"), derive.MustDecimal("100"))
	require.NoError(t, p.Validate())
}

func TestOrderParams_Validate_RejectsBlankInstrument(t *testing.T) {
	p := derive.NewOrderParams("", derive.DirectionBuy, derive.OrderTypeLimit,
		derive.MustDecimal("1"), derive.MustDecimal("100"))
	err := p.Validate()
	assert.ErrorIs(t, err, derive.ErrInvalidParams)
	assert.Contains(t, err.Error(), "instrument_name")
}

func TestOrderParams_Validate_RejectsBadEnum(t *testing.T) {
	p := derive.NewOrderParams("BTC-PERP", derive.Direction("up"), derive.OrderTypeLimit,
		derive.MustDecimal("1"), derive.MustDecimal("100"))
	err := p.Validate()
	assert.ErrorIs(t, err, derive.ErrInvalidParams)
	assert.Contains(t, err.Error(), "direction")
}

func TestOrderParams_Validate_RejectsZeroAmount(t *testing.T) {
	p := derive.NewOrderParams("BTC-PERP", derive.DirectionBuy, derive.OrderTypeLimit,
		derive.MustDecimal("0"), derive.MustDecimal("100"))
	err := p.Validate()
	assert.ErrorIs(t, err, derive.ErrInvalidParams)
	assert.Contains(t, err.Error(), "amount")
}

func TestOrderParams_Validate_RejectsNegativePrice(t *testing.T) {
	p := derive.NewOrderParams("BTC-PERP", derive.DirectionBuy, derive.OrderTypeLimit,
		derive.MustDecimal("1"), derive.MustDecimal("-1"))
	err := p.Validate()
	assert.ErrorIs(t, err, derive.ErrInvalidParams)
	assert.Contains(t, err.Error(), "limit_price")
}

func TestOrderParams_Validate_RejectsNegativeMaxFee(t *testing.T) {
	p := derive.NewOrderParams("BTC-PERP", derive.DirectionBuy, derive.OrderTypeLimit,
		derive.MustDecimal("1"), derive.MustDecimal("100")).
		WithMaxFee(derive.MustDecimal("-0.1"))
	err := p.Validate()
	assert.ErrorIs(t, err, derive.ErrInvalidParams)
	assert.Contains(t, err.Error(), "max_fee")
}

func TestCancelOrderParams_Builder_AndValidate(t *testing.T) {
	got := derive.NewCancelOrderParams(7).WithOrderID("O1").WithInstrument("BTC-PERP")
	require.NoError(t, got.Validate())
	assert.Equal(t, "O1", got.OrderID)
	assert.Equal(t, "BTC-PERP", got.InstrumentName)

	got2 := derive.NewCancelOrderParams(7).WithLabel("scalp")
	require.NoError(t, got2.Validate())

	noTarget := derive.NewCancelOrderParams(7)
	err := noTarget.Validate()
	assert.ErrorIs(t, err, derive.ErrInvalidParams)

	bad := derive.NewCancelOrderParams(-1).WithOrderID("O1")
	assert.Error(t, bad.Validate())
}

func TestReplaceOrderParams_Builder_AndValidate(t *testing.T) {
	o := derive.NewOrderParams("BTC-PERP", derive.DirectionBuy, derive.OrderTypeLimit,
		derive.MustDecimal("1"), derive.MustDecimal("100"))
	r := derive.NewReplaceOrderParams("OLD", o)
	require.NoError(t, r.Validate())
	assert.Equal(t, "OLD", r.OrderIDToCancel)

	bad := derive.NewReplaceOrderParams("", o)
	assert.ErrorIs(t, bad.Validate(), derive.ErrInvalidParams)

	badInner := derive.NewReplaceOrderParams("OLD",
		derive.NewOrderParams("", derive.DirectionBuy, derive.OrderTypeLimit,
			derive.MustDecimal("1"), derive.MustDecimal("100")))
	err := badInner.Validate()
	assert.True(t, errors.Is(err, derive.ErrInvalidParams))
}

func TestPageRequest_Builder_AndValidate(t *testing.T) {
	p := derive.NewPageRequest().WithPage(2).WithPageSize(50)
	require.NoError(t, p.Validate())
	assert.Equal(t, 2, p.Page)
	assert.Equal(t, 50, p.PageSize)

	require.NoError(t, derive.NewPageRequest().Validate())
	assert.ErrorIs(t, derive.NewPageRequest().WithPage(-1).Validate(), derive.ErrInvalidParams)
	assert.ErrorIs(t, derive.NewPageRequest().WithPageSize(-1).Validate(), derive.ErrInvalidParams)
}

// TestOrderParams_Builder_AllSetters exercises every With* setter at least
// once, verifying the value is applied and the original is not mutated.
func TestOrderParams_Builder_AllSetters(t *testing.T) {
	base := derive.NewOrderParams("X", derive.DirectionBuy, derive.OrderTypeLimit,
		derive.MustDecimal("1"), derive.MustDecimal("100"))

	a := base.WithInstrument("BTC-PERP")
	assert.Equal(t, "BTC-PERP", a.InstrumentName)

	b := base.WithDirection(derive.DirectionSell)
	assert.Equal(t, derive.DirectionSell, b.Direction)

	c := base.WithOrderType(derive.OrderTypeMarket)
	assert.Equal(t, derive.OrderTypeMarket, c.OrderType)

	d := base.WithAmount(derive.MustDecimal("2"))
	assert.Equal(t, "2", d.Amount.String())

	e := base.WithLimitPrice(derive.MustDecimal("200"))
	assert.Equal(t, "200", e.LimitPrice.String())

	addr := derive.MustAddress("0x1111111111111111111111111111111111111111")
	f := base.WithSignature(addr, "0xsig", 99, 1700000000)
	assert.Equal(t, derive.MustAddress("0x1111111111111111111111111111111111111111"), f.Signer)
	assert.Equal(t, "0xsig", f.Signature)
	assert.Equal(t, uint64(99), f.Nonce)
	assert.Equal(t, int64(1700000000), f.SignatureExpiry)

	assert.Equal(t, "X", base.InstrumentName)
}

func TestCancelOrderParams_WithSignature(t *testing.T) {
	addr := derive.MustAddress("0x1111111111111111111111111111111111111111")
	p := derive.NewCancelOrderParams(7).WithOrderID("O1").
		WithSignature(addr, "0xsig", 11, 1700000000)
	assert.Equal(t, derive.MustAddress("0x1111111111111111111111111111111111111111"), p.Signer)
	assert.Equal(t, "0xsig", p.Signature)
	assert.Equal(t, uint64(11), p.Nonce)
	assert.Equal(t, int64(1700000000), p.SignatureExpiry)
}

func TestReplaceOrderParams_AllSetters(t *testing.T) {
	o := derive.NewOrderParams("BTC-PERP", derive.DirectionBuy, derive.OrderTypeLimit,
		derive.MustDecimal("1"), derive.MustDecimal("100"))
	r := derive.NewReplaceOrderParams("X", o)

	r2 := r.WithOrderIDToCancel("OLD")
	assert.Equal(t, "OLD", r2.OrderIDToCancel)
	assert.Equal(t, "X", r.OrderIDToCancel)

	o2 := derive.NewOrderParams("ETH-PERP", derive.DirectionSell, derive.OrderTypeMarket,
		derive.MustDecimal("2"), derive.MustDecimal("3500"))
	r3 := r.WithNewOrder(o2)
	assert.Equal(t, "ETH-PERP", r3.NewOrder.InstrumentName)
}

// TestOrderParams_Validate_RejectsRemainingBranches covers branches the
// existing happy-path tests don't reach: bad order_type, bad time_in_force,
// negative subaccount, negative signature_expiry.
func TestOrderParams_Validate_RejectsRemainingBranches(t *testing.T) {
	cases := []struct {
		name string
		mut  func(*derive.OrderParams)
		want string
	}{
		{"bad order_type", func(p *derive.OrderParams) { p.OrderType = derive.OrderType("nope") }, "order_type"},
		{"bad tif", func(p *derive.OrderParams) { p.TimeInForce = derive.TimeInForce("forever") }, "time_in_force"},
		{"negative subaccount", func(p *derive.OrderParams) { p.SubaccountID = -1 }, "subaccount_id"},
		{"negative expiry", func(p *derive.OrderParams) { p.SignatureExpiry = -1 }, "signature_expiry"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := derive.NewOrderParams("BTC-PERP", derive.DirectionBuy, derive.OrderTypeLimit,
				derive.MustDecimal("1"), derive.MustDecimal("100"))
			c.mut(&p)
			err := p.Validate()
			require.Error(t, err)
			assert.True(t, errors.Is(err, derive.ErrInvalidParams))
			assert.Contains(t, err.Error(), c.want)
		})
	}
}

func TestCancelOrderParams_Validate_RejectsBadExpiry(t *testing.T) {
	p := derive.NewCancelOrderParams(0).WithOrderID("O1")
	p.SignatureExpiry = -1
	assert.ErrorIs(t, p.Validate(), derive.ErrInvalidParams)
}
func TestOrder_DecodeFull(t *testing.T) {
	payload := `{
		"order_id": "O1",
		"subaccount_id": 1,
		"instrument_name": "BTC-PERP",
		"direction": "buy",
		"order_type": "limit",
		"time_in_force": "gtc",
		"order_status": "open",
		"amount": "0.1",
		"filled_amount": "0",
		"limit_price": "65000",
		"average_price": "0",
		"max_fee": "10",
		"nonce": 1,
		"signer": "0x1111111111111111111111111111111111111111",
		"label": "alpha",
		"cancel_reason": "",
		"mmp": false,
		"reduce_only": true,
		"creation_timestamp": 1700000000000,
		"last_update_timestamp": 1700000000005
	}`
	var o derive.Order
	require.NoError(t, json.Unmarshal([]byte(payload), &o))
	assert.Equal(t, "O1", o.OrderID)
	assert.Equal(t, derive.DirectionBuy, o.Direction)
	assert.Equal(t, derive.OrderStatusOpen, o.OrderStatus)
	assert.True(t, o.ReduceOnly)
	assert.Equal(t, "alpha", o.Label)
}

func TestOrderParams_OmitsEmptyOptionalFields(t *testing.T) {
	in := derive.OrderParams{
		InstrumentName:  "BTC-PERP",
		Direction:       derive.DirectionBuy,
		OrderType:       derive.OrderTypeLimit,
		Amount:          derive.MustDecimal("1"),
		LimitPrice:      derive.MustDecimal("100"),
		MaxFee:          derive.MustDecimal("10"),
		SubaccountID:    1,
		Nonce:           1,
		Signer:          derive.MustAddress("0x1111111111111111111111111111111111111111"),
		Signature:       "0xdeadbeef",
		SignatureExpiry: 1700000000,
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	s := string(b)
	assert.NotContains(t, s, `"label"`)
	assert.NotContains(t, s, `"trigger_type"`)
	assert.NotContains(t, s, `"reduce_only":true`)
}

func TestOrderParams_IncludesPopulatedOptionals(t *testing.T) {
	in := derive.OrderParams{
		InstrumentName: "BTC-PERP",
		Direction:      derive.DirectionSell,
		OrderType:      derive.OrderTypeLimit,
		Label:          "lbl",
		ReduceOnly:     true,
		MMP:            true,
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	s := string(b)
	assert.Contains(t, s, `"label":"lbl"`)
	assert.Contains(t, s, `"reduce_only":true`)
	assert.Contains(t, s, `"mmp":true`)
}

func TestCancelOrderParams_RoundTrip(t *testing.T) {
	in := derive.CancelOrderParams{
		SubaccountID:   1,
		InstrumentName: "BTC-PERP",
		OrderID:        "O1",
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	var out derive.CancelOrderParams
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in, out)
}

func TestReplaceOrderParams_Embedded(t *testing.T) {
	in := derive.ReplaceOrderParams{
		OrderIDToCancel: "O1",
		NewOrder: derive.OrderParams{
			InstrumentName: "BTC-PERP",
			Direction:      derive.DirectionBuy,
		},
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	assert.Contains(t, string(b), `"order_id_to_cancel":"O1"`)
	assert.Contains(t, string(b), `"new_order"`)
}
func TestPage_RoundTrip(t *testing.T) {
	in := derive.Page{NumPages: 5, Count: 100}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	var out derive.Page
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in, out)
}

// Derive may add fields like `current_page` later — the SDK should accept
// them silently rather than fail the decode.
func TestPage_IgnoresUnknownFields(t *testing.T) {
	raw := []byte(`{"num_pages":3,"count":50,"current_page":2,"page_size":20}`)
	var p derive.Page
	require.NoError(t, json.Unmarshal(raw, &p))
	assert.Equal(t, 3, p.NumPages)
	assert.Equal(t, 50, p.Count)
}

func TestPageRequest_OmitsZeroFields(t *testing.T) {
	in := derive.PageRequest{}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	assert.JSONEq(t, `{}`, string(b))
}

func TestPageRequest_IncludesPopulated(t *testing.T) {
	in := derive.PageRequest{Page: 2, PageSize: 25}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	assert.JSONEq(t, `{"page":2,"page_size":25}`, string(b))
}
func TestPosition_DecodeFull(t *testing.T) {
	payload := `{
		"instrument_name": "BTC-PERP",
		"instrument_type": "perp",
		"amount": "0.5",
		"average_price": "65000",
		"mark_price": "65500",
		"mark_value": "32750",
		"index_price": "65500",
		"leverage": "5",
		"liquidation_price": "10000",
		"unrealized_pnl": "250",
		"realized_pnl": "10",
		"open_orders_margin": "100",
		"cumulative_funding": "1",
		"pending_funding": "0.1"
	}`
	var p derive.Position
	require.NoError(t, json.Unmarshal([]byte(payload), &p))
	assert.Equal(t, "BTC-PERP", p.InstrumentName)
	assert.Equal(t, derive.InstrumentTypePerp, p.InstrumentType)
	assert.Equal(t, "0.5", p.Amount.String())
	assert.Equal(t, "250", p.UnrealizedPNL.String())
	assert.Equal(t, "5", p.Leverage.String())
}

func TestPosition_DecodeMinimal(t *testing.T) {
	payload := `{
		"instrument_name": "BTC-PERP",
		"instrument_type": "perp",
		"amount": "0",
		"average_price": "0",
		"mark_price": "0",
		"mark_value": "0",
		"unrealized_pnl": "0",
		"realized_pnl": "0"
	}`
	var p derive.Position
	require.NoError(t, json.Unmarshal([]byte(payload), &p))
	assert.True(t, p.Amount.IsZero())
}

func TestPosition_RoundTrip(t *testing.T) {
	in := derive.Position{
		InstrumentName: "BTC-PERP",
		InstrumentType: derive.InstrumentTypePerp,
		Amount:         derive.MustDecimal("1"),
		AveragePrice:   derive.MustDecimal("65000"),
		MarkPrice:      derive.MustDecimal("65500"),
		MarkValue:      derive.MustDecimal("65500"),
		UnrealizedPNL:  derive.MustDecimal("500"),
		RealizedPNL:    derive.MustDecimal("0"),
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	var out derive.Position
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.InstrumentName, out.InstrumentName)
	assert.Equal(t, in.UnrealizedPNL.String(), out.UnrealizedPNL.String())
}
func validLeg() derive.RFQLeg {
	return derive.RFQLeg{
		InstrumentName: "BTC-PERP",
		Direction:      derive.DirectionBuy,
		Amount:         derive.MustDecimal("1"),
	}
}

func TestRFQLeg_Validate_Happy(t *testing.T) {
	require.NoError(t, validLeg().Validate())
}

func TestRFQLeg_Validate_Rejects(t *testing.T) {
	cases := []struct {
		name string
		mut  func(*derive.RFQLeg)
		want string
	}{
		{"empty instrument", func(l *derive.RFQLeg) { l.InstrumentName = "" }, "instrument_name"},
		{"bad direction", func(l *derive.RFQLeg) { l.Direction = derive.Direction("sideways") }, "direction"},
		{"zero amount", func(l *derive.RFQLeg) { l.Amount = derive.MustDecimal("0") }, "amount"},
		{"negative amount", func(l *derive.RFQLeg) { l.Amount = derive.MustDecimal("-1") }, "amount"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := validLeg()
			c.mut(&l)
			err := l.Validate()
			require.Error(t, err)
			assert.True(t, errors.Is(err, derive.ErrInvalidParams))
			assert.Contains(t, err.Error(), c.want)
		})
	}
}

func TestRFQLeg_RoundTrip(t *testing.T) {

	in := derive.RFQLeg{
		InstrumentName: "BTC-PERP",
		Direction:      derive.DirectionBuy,
		Amount:         derive.MustDecimal("1"),
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	var out derive.RFQLeg
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.InstrumentName, out.InstrumentName)
	assert.Equal(t, in.Direction, out.Direction)
}

func TestQuoteLeg_RoundTrip(t *testing.T) {
	in := derive.QuoteLeg{
		InstrumentName: "BTC-PERP",
		Direction:      derive.DirectionBuy,
		Amount:         derive.MustDecimal("1"),
		Price:          derive.MustDecimal("65000"),
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	var out derive.QuoteLeg
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, "65000", out.Price.String())
}

func TestRFQ_Decode(t *testing.T) {
	payload := `{
		"rfq_id": "R1",
		"subaccount_id": 1,
		"status": "open",
		"legs": [
			{"instrument_name":"BTC-PERP","direction":"buy","amount":"1"}
		],
		"max_total_fee": "10",
		"creation_timestamp": 1700000000000,
		"last_update_timestamp": 1700000000001
	}`
	var rfq derive.RFQ
	require.NoError(t, json.Unmarshal([]byte(payload), &rfq))
	assert.Equal(t, "R1", rfq.RFQID)
	assert.Equal(t, derive.QuoteStatusOpen, rfq.Status)
	require.Len(t, rfq.Legs, 1)
}

func TestQuote_Decode(t *testing.T) {
	payload := `{
		"quote_id": "Q1",
		"rfq_id": "R1",
		"subaccount_id": 1,
		"direction": "sell",
		"legs": [{"instrument_name":"BTC-PERP","direction":"sell","amount":"1","price":"65000"}],
		"price": "65000",
		"status": "open",
		"creation_timestamp": 1700000000000
	}`
	var q derive.Quote
	require.NoError(t, json.Unmarshal([]byte(payload), &q))
	assert.Equal(t, "Q1", q.QuoteID)
	assert.Equal(t, derive.DirectionSell, q.Direction)
	assert.Equal(t, derive.QuoteStatusOpen, q.Status)
	require.Len(t, q.Legs, 1)
	assert.Equal(t, "65000", q.Legs[0].Price.String())
}
func TestSpotFeed_Decode(t *testing.T) {

	raw := `{
		"timestamp": 1777842232556,
		"feeds": {
			"BTC": {
				"price": "78908.29",
				"confidence": "1",
				"price_prev_daily": "78689.04",
				"confidence_prev_daily": "1",
				"timestamp_prev_daily": 1777755832556
			}
		}
	}`
	var sf derive.SpotFeed
	require.NoError(t, json.Unmarshal([]byte(raw), &sf))
	assert.Equal(t, int64(1777842232556), sf.Timestamp.Millis())
	require.Contains(t, sf.Feeds, "BTC")
	btc := sf.Feeds["BTC"]
	assert.Equal(t, "78908.29", btc.Price.String())
	assert.Equal(t, "78689.04", btc.PricePrevDaily.String())
	assert.Equal(t, int64(1777755832556), btc.TimestampPrevDaily.Millis())
}
func TestSubAccount_Decode(t *testing.T) {
	payload := `{
		"subaccount_id": 7,
		"owner_address": "0x1111111111111111111111111111111111111111",
		"margin_type": "PM",
		"is_under_liquidation": false,
		"subaccount_value": "100",
		"initial_margin": "50",
		"maintenance_margin": "30"
	}`
	var sa derive.SubAccount
	require.NoError(t, json.Unmarshal([]byte(payload), &sa))
	assert.Equal(t, int64(7), sa.SubaccountID)
	assert.Equal(t, "PM", sa.MarginType)
	assert.False(t, sa.IsUnderLiquidation)
}

func TestSubAccount_RoundTrip(t *testing.T) {
	in := derive.SubAccount{
		SubaccountID:      1,
		OwnerAddress:      derive.MustAddress("0x1111111111111111111111111111111111111111"),
		MarginType:        "SM",
		SubaccountValue:   derive.MustDecimal("0"),
		InitialMargin:     derive.MustDecimal("0"),
		MaintenanceMargin: derive.MustDecimal("0"),
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	var out derive.SubAccount
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.SubaccountID, out.SubaccountID)
	assert.Equal(t, in.OwnerAddress.String(), out.OwnerAddress.String())
}
func TestTickerSlim_DecodeFull(t *testing.T) {

	payload := `{
		"timestamp": 1700000001000,
		"instrument_ticker": {
			"t": 1700000001000,
			"A": "0.5",  "a": "78758.5",
			"B": "0.4",  "b": "78752.1",
			"I": "78760", "M": "78755",
			"f": "0.0001",
			"stats": {"v": "100"},
			"option_pricing": null
		}
	}`
	var ts derive.TickerSlim
	require.NoError(t, json.Unmarshal([]byte(payload), &ts))
	assert.Equal(t, int64(1700000001000), ts.Timestamp.Millis())
	assert.Equal(t, int64(1700000001000), ts.Ticker.Timestamp.Millis())
	assert.Equal(t, "0.5", ts.Ticker.BestAskAmount.String())
	assert.Equal(t, "78758.5", ts.Ticker.BestAskPrice.String())
	assert.Equal(t, "0.4", ts.Ticker.BestBidAmount.String())
	assert.Equal(t, "78752.1", ts.Ticker.BestBidPrice.String())
	assert.Equal(t, "78760", ts.Ticker.IndexPrice.String())
	assert.Equal(t, "78755", ts.Ticker.MarkPrice.String())
	assert.Equal(t, "0.0001", ts.Ticker.FundingRate.String())
	assert.JSONEq(t, `{"v":"100"}`, string(ts.Ticker.Stats))

	assert.Equal(t, "null", string(ts.Ticker.OptionPricing))
}

func TestTickerSlim_DecodeMinimal(t *testing.T) {

	payload := `{
		"timestamp": 1,
		"instrument_ticker": {
			"t": 1,
			"A": "0", "a": "0",
			"B": "0", "b": "0"
		}
	}`
	var ts derive.TickerSlim
	require.NoError(t, json.Unmarshal([]byte(payload), &ts))
	assert.Equal(t, "0", ts.Ticker.BestAskAmount.String())
	assert.Equal(t, json.RawMessage(nil), ts.Ticker.Stats)
	assert.Equal(t, json.RawMessage(nil), ts.Ticker.OptionPricing)
}

func TestTickerSlim_RoundTrip_PreservesTopFields(t *testing.T) {
	var in derive.TickerSlim
	in.Ticker.MarkPrice = derive.MustDecimal("100.5")
	in.Ticker.IndexPrice = derive.MustDecimal("100.5")
	b, err := json.Marshal(in)
	require.NoError(t, err)
	var out derive.TickerSlim
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.Ticker.MarkPrice.String(), out.Ticker.MarkPrice.String())
}
func TestTicker_DecodeFull(t *testing.T) {

	payload := `{
		"instrument_name": "BTC-PERP",
		"instrument_type": "perp",
		"is_active": true,
		"best_bid_price": "100",
		"best_bid_amount": "1",
		"best_ask_price": "101",
		"best_ask_amount": "2",
		"five_percent_bid_depth": "10",
		"five_percent_ask_depth": "12",
		"mark_price": "100.5",
		"index_price": "100.5",
		"min_price": "90",
		"max_price": "110",
		"open_interest": {"PM": [{"current_open_interest": "1", "interest_cap": "100", "manager_currency": "BTC"}]},
		"timestamp": 1700000000000
	}`
	var tk derive.Ticker
	require.NoError(t, json.Unmarshal([]byte(payload), &tk))
	assert.Equal(t, "BTC-PERP", tk.InstrumentName)
	assert.Equal(t, "perp", tk.InstrumentType)
	assert.True(t, tk.IsActive)
	assert.Equal(t, "100", tk.BestBidPrice.String())
	assert.Equal(t, "10", tk.FivePercentBidDepth.String())
	assert.Equal(t, "90", tk.MinPrice.String())
	assert.Contains(t, string(tk.OpenInterest), "current_open_interest")
}

func TestTicker_DecodeMinimal(t *testing.T) {
	payload := `{
		"instrument_name": "BTC-PERP",
		"best_bid_price": "0",
		"best_bid_amount": "0",
		"best_ask_price": "0",
		"best_ask_amount": "0",
		"mark_price": "0",
		"index_price": "0",
		"timestamp": 0
	}`
	var tk derive.Ticker
	require.NoError(t, json.Unmarshal([]byte(payload), &tk))
	assert.Equal(t, "BTC-PERP", tk.InstrumentName)
}

func TestTicker_RoundTrip(t *testing.T) {
	in := derive.Ticker{
		InstrumentName: "BTC-PERP",
		BestBidPrice:   derive.MustDecimal("100"),
		BestAskPrice:   derive.MustDecimal("101"),
		MarkPrice:      derive.MustDecimal("100.5"),
		IndexPrice:     derive.MustDecimal("100.5"),
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	var out derive.Ticker
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.InstrumentName, out.InstrumentName)
	assert.Equal(t, in.BestBidPrice.String(), out.BestBidPrice.String())
}

// FuzzMillisTime_UnmarshalJSON ensures the time unmarshaler is panic-free
// on arbitrary input, including non-string-non-number JSON.
func FuzzMillisTime_UnmarshalJSON(f *testing.F) {
	f.Add([]byte(`1700000000000`))
	f.Add([]byte(`"1700000000000"`))
	f.Add([]byte(`null`))
	f.Add([]byte(`""`))
	f.Add([]byte(`"abc"`))
	f.Add([]byte(`-1`))
	f.Add([]byte(`{}`))

	f.Fuzz(func(t *testing.T, raw []byte) {
		var mt derive.MillisTime
		_ = mt.UnmarshalJSON(raw)
	})
}
func TestMillisTime_RoundTripFromNumber(t *testing.T) {
	now := time.Date(2026, 5, 2, 10, 0, 0, 0, time.UTC)
	mt := derive.NewMillisTime(now)
	b, err := json.Marshal(mt)
	require.NoError(t, err)

	var got derive.MillisTime
	require.NoError(t, json.Unmarshal(b, &got))
	assert.Equal(t, now.UnixMilli(), got.Millis())
	assert.Equal(t, now.UnixMilli(), got.Time().UnixMilli())
}

func TestMillisTime_UnmarshalString(t *testing.T) {
	var mt derive.MillisTime
	require.NoError(t, json.Unmarshal([]byte(`"1700000000000"`), &mt))
	assert.Equal(t, int64(1700000000000), mt.Millis())
}

func TestMillisTime_UnmarshalNullEmpty(t *testing.T) {
	var mt derive.MillisTime
	require.NoError(t, json.Unmarshal([]byte(`null`), &mt))
	assert.True(t, mt.Time().IsZero())

	require.NoError(t, json.Unmarshal([]byte(`""`), &mt))
	assert.True(t, mt.Time().IsZero())
}

func TestMillisTime_UnmarshalInvalid(t *testing.T) {
	var mt derive.MillisTime
	assert.Error(t, json.Unmarshal([]byte(`"abc"`), &mt))
	assert.Error(t, json.Unmarshal([]byte(`{`), &mt))
}
func TestTrade_Decode(t *testing.T) {
	payload := `{
		"trade_id": "T1",
		"order_id": "O1",
		"subaccount_id": 1,
		"instrument_name": "BTC-PERP",
		"direction": "buy",
		"trade_price": "65000",
		"trade_amount": "0.1",
		"mark_price": "65000",
		"index_price": "64999",
		"trade_fee": "0.5",
		"liquidity_role": "taker",
		"realized_pnl": "10",
		"timestamp": 1700000000000
	}`
	var tr derive.Trade
	require.NoError(t, json.Unmarshal([]byte(payload), &tr))
	assert.Equal(t, "T1", tr.TradeID)
	assert.Equal(t, derive.DirectionBuy, tr.Direction)
	assert.Equal(t, derive.LiquidityRoleTaker, tr.LiquidityRole)
}

func TestTrade_OmitsEmpty(t *testing.T) {
	in := derive.Trade{
		TradeID:        "T1",
		InstrumentName: "BTC-PERP",
		Direction:      derive.DirectionBuy,
		TradePrice:     derive.MustDecimal("100"),
		TradeAmount:    derive.MustDecimal("1"),
		MarkPrice:      derive.MustDecimal("100"),
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	s := string(b)
	assert.NotContains(t, s, `"order_id"`)
	assert.NotContains(t, s, `"liquidity_role"`)
}
func TestDepositTx_Decode(t *testing.T) {
	payload := `{
		"tx_hash": "0x1111111111111111111111111111111111111111111111111111111111111111",
		"asset": "USDC",
		"amount": "1000",
		"subaccount_id": 1,
		"status": "completed",
		"timestamp": 1700000000000
	}`
	var d derive.DepositTx
	require.NoError(t, json.Unmarshal([]byte(payload), &d))
	assert.Equal(t, "USDC", d.Asset)
	assert.Equal(t, "completed", d.Status)
}

func TestWithdrawTx_Decode(t *testing.T) {
	payload := `{
		"tx_hash": "0x2222222222222222222222222222222222222222222222222222222222222222",
		"asset": "USDC",
		"amount": "500",
		"subaccount_id": 1,
		"status": "pending",
		"timestamp": 1700000000000
	}`
	var w derive.WithdrawTx
	require.NoError(t, json.Unmarshal([]byte(payload), &w))
	assert.Equal(t, "pending", w.Status)
}

func TestDepositTx_RoundTrip(t *testing.T) {
	in := derive.DepositTx{Asset: "USDC", Amount: derive.MustDecimal("1"), Status: "completed"}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	var out derive.DepositTx
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.Asset, out.Asset)
	assert.Equal(t, in.Status, out.Status)
}

const sampleHash = "0x1111111111111111111111111111111111111111111111111111111111111111"

func TestNewTxHash_Valid(t *testing.T) {
	h, err := derive.NewTxHash(sampleHash)
	require.NoError(t, err)
	assert.False(t, h.IsZero())
	assert.Equal(t, sampleHash, h.String())
}

func TestNewTxHash_Empty(t *testing.T) {
	h, err := derive.NewTxHash("")
	require.NoError(t, err)
	assert.True(t, h.IsZero())
}

func TestNewTxHash_Invalid(t *testing.T) {
	_, err := derive.NewTxHash("0xabc")
	assert.Error(t, err)
	_, err = derive.NewTxHash("notahash")
	assert.Error(t, err)
}

func TestTxHash_JSONRoundTrip(t *testing.T) {
	type wrap struct {
		H derive.TxHash `json:"h"`
	}
	in := wrap{H: must(derive.NewTxHash(sampleHash))}
	b, err := json.Marshal(in)
	require.NoError(t, err)

	var out wrap
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.H, out.H)
}

func TestTxHash_UnmarshalEmpty(t *testing.T) {
	var h derive.TxHash
	require.NoError(t, json.Unmarshal([]byte(`""`), &h))
	assert.True(t, h.IsZero())
}

func TestTxHash_UnmarshalInvalid(t *testing.T) {
	var h derive.TxHash
	assert.Error(t, json.Unmarshal([]byte(`"0xabc"`), &h))
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

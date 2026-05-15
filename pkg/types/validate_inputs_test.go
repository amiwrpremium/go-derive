package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

// Each test covers the happy path and one representative failure case
// per Input/Query type. The Validate methods are simple required-field
// checks, so deeper coverage doesn't add value.

func TestCancelOrderInput_Validate(t *testing.T) {
	require.NoError(t, types.CancelOrderInput{InstrumentName: "BTC-PERP", OrderID: "O1"}.Validate())
	err := types.CancelOrderInput{OrderID: "O1"}.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "instrument_name")
	err = types.CancelOrderInput{InstrumentName: "BTC-PERP"}.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "order_id")
}

func TestCancelByInstrumentInput_Validate(t *testing.T) {
	require.NoError(t, types.CancelByInstrumentInput{InstrumentName: "BTC-PERP"}.Validate())
	err := types.CancelByInstrumentInput{}.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "instrument_name")
}

func TestCancelByLabelInput_Validate(t *testing.T) {
	require.NoError(t, types.CancelByLabelInput{Label: "alpha"}.Validate())
	err := types.CancelByLabelInput{}.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "label")
}

func TestCancelByNonceInput_Validate(t *testing.T) {
	require.NoError(t, types.CancelByNonceInput{InstrumentName: "BTC-PERP", Nonce: 1}.Validate())
	err := types.CancelByNonceInput{Nonce: 1}.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "instrument_name")
	err = types.CancelByNonceInput{InstrumentName: "BTC-PERP"}.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nonce")
}

func TestCancelAlgoOrderInput_Validate(t *testing.T) {
	require.NoError(t, types.CancelAlgoOrderInput{OrderID: "A1"}.Validate())
	err := types.CancelAlgoOrderInput{}.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "order_id")
}

func TestCancelTriggerOrderInput_Validate(t *testing.T) {
	require.NoError(t, types.CancelTriggerOrderInput{OrderID: "T1"}.Validate())
	err := types.CancelTriggerOrderInput{}.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "order_id")
}

func TestCancelRFQInput_Validate(t *testing.T) {
	require.NoError(t, types.CancelRFQInput{RFQID: "R1"}.Validate())
	err := types.CancelRFQInput{}.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "rfq_id")
}

func TestCancelQuoteInput_Validate(t *testing.T) {
	require.NoError(t, types.CancelQuoteInput{QuoteID: "Q1"}.Validate())
	err := types.CancelQuoteInput{}.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "quote_id")
}

func TestChangeSubaccountLabelInput_Validate(t *testing.T) {
	require.NoError(t, types.ChangeSubaccountLabelInput{Label: "alpha"}.Validate())
	err := types.ChangeSubaccountLabelInput{}.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "label")
}

func TestCreateContactInfoInput_Validate(t *testing.T) {
	require.NoError(t, types.CreateContactInfoInput{ContactType: "email", ContactValue: "a@b.c"}.Validate())
	err := types.CreateContactInfoInput{ContactValue: "a@b.c"}.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "contact_type")
	err = types.CreateContactInfoInput{ContactType: "email"}.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "contact_value")
}

func TestUpdateContactInfoInput_Validate(t *testing.T) {
	require.NoError(t, types.UpdateContactInfoInput{ContactID: 7, NewValue: "new"}.Validate())
	err := types.UpdateContactInfoInput{NewValue: "new"}.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "contact_id")
	err = types.UpdateContactInfoInput{ContactID: 7}.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "contact_value")
}

func TestDeleteContactInfoInput_Validate(t *testing.T) {
	require.NoError(t, types.DeleteContactInfoInput{ContactID: 7}.Validate())
	err := types.DeleteContactInfoInput{}.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "contact_id")
}

func TestResetMMPInput_Validate(t *testing.T) {
	require.NoError(t, types.ResetMMPInput{Currency: "BTC"}.Validate())
	err := types.ResetMMPInput{}.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "currency")
}

func TestMarginWatchQuery_Validate(t *testing.T) {
	require.NoError(t, types.MarginWatchQuery{SubaccountID: 7}.Validate())
	err := types.MarginWatchQuery{}.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "subaccount_id")
}

func TestSendRFQInput_Validate(t *testing.T) {
	good := types.SendRFQInput{Legs: []types.RFQLeg{{
		InstrumentName: "BTC-PERP",
		Direction:      enums.DirectionBuy,
		Amount:         types.MustDecimal("1"),
	}}}
	require.NoError(t, good.Validate())

	err := types.SendRFQInput{}.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "legs")

	// A leg with an invalid field bubbles up with the leg index.
	bad := types.SendRFQInput{Legs: []types.RFQLeg{
		{InstrumentName: "BTC-PERP", Direction: enums.DirectionBuy, Amount: types.MustDecimal("1")},
		{InstrumentName: "", Direction: enums.DirectionBuy, Amount: types.MustDecimal("1")},
	}}
	err = bad.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "leg 1")
}

func TestCancelBatchInput_Validate(t *testing.T) {
	// Each filter satisfies the requirement on its own.
	for _, in := range []types.CancelBatchInput{
		{RFQID: "R1"},
		{QuoteID: "Q1"},
		{Label: "alpha"},
		{Nonce: 7},
	} {
		require.NoError(t, in.Validate())
	}
	// All-zero (or only SubaccountID set) fails — would no-op server-side.
	err := types.CancelBatchInput{}.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "rfq_id")
	err = types.CancelBatchInput{SubaccountID: 1}.Validate()
	require.Error(t, err)
}

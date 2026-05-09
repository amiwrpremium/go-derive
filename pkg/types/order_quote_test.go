package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestOrderQuoteResult_Valid(t *testing.T) {
	raw := []byte(`{
		"is_valid":true,
		"invalid_reason":null,
		"estimated_fill_amount":"0.5",
		"estimated_fill_price":"50000",
		"estimated_fee":"5",
		"estimated_realized_pnl":"0",
		"estimated_order_status":"filled",
		"suggested_max_fee":"10",
		"pre_initial_margin":"100",
		"post_initial_margin":"120",
		"post_liquidation_price":null
	}`)
	var got types.OrderQuoteResult
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.True(t, got.IsValid)
	assert.Equal(t, enums.RFQInvalidReason(""), got.InvalidReason, "valid request decodes to empty zero-value")
	assert.Equal(t, "0.5", got.EstimatedFillAmount.String())
	assert.Equal(t, enums.OrderStatusFilled, got.EstimatedOrderStatus)
	assert.Equal(t, "0", got.PostLiquidationPrice.String(), "null decimal decodes to zero-value")
}

func TestOrderQuoteResult_Invalid(t *testing.T) {
	raw := []byte(`{
		"is_valid":false,
		"invalid_reason":"Insufficient buying power.",
		"estimated_fill_amount":"0",
		"estimated_fill_price":"0",
		"estimated_fee":"0",
		"estimated_realized_pnl":"0",
		"estimated_order_status":"rejected",
		"suggested_max_fee":"0",
		"pre_initial_margin":"100",
		"post_initial_margin":"100",
		"post_liquidation_price":"45000"
	}`)
	var got types.OrderQuoteResult
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.False(t, got.IsValid)
	assert.Equal(t, enums.RFQInvalidReasonInsufficientBuyingPower, got.InvalidReason)
	assert.Equal(t, enums.OrderStatusRejected, got.EstimatedOrderStatus)
	assert.Equal(t, "45000", got.PostLiquidationPrice.String())
}

package enums_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/enums"
)

func TestDirection_Valid_Buy(t *testing.T) {
	assert.True(t, enums.DirectionBuy.Valid())
}

func TestDirection_Valid_Sell(t *testing.T) {
	assert.True(t, enums.DirectionSell.Valid())
}

func TestDirection_Valid_RejectsUnknown(t *testing.T) {
	assert.False(t, enums.Direction("unknown").Valid())
	assert.False(t, enums.Direction("").Valid())
	assert.False(t, enums.Direction("BUY").Valid(), "case-sensitive")
}

func TestDirection_Sign_Buy(t *testing.T) {
	assert.Equal(t, 1, enums.DirectionBuy.Sign())
}

func TestDirection_Sign_Sell(t *testing.T) {
	assert.Equal(t, -1, enums.DirectionSell.Sign())
}

func TestDirection_Sign_PanicsOnInvalid(t *testing.T) {
	assert.Panics(t, func() { _ = enums.Direction("nope").Sign() })
	assert.Panics(t, func() { _ = enums.Direction("").Sign() })
}

func TestDirection_Opposite_BuyToSell(t *testing.T) {
	assert.Equal(t, enums.DirectionSell, enums.DirectionBuy.Opposite())
}

func TestDirection_Opposite_SellToBuy(t *testing.T) {
	assert.Equal(t, enums.DirectionBuy, enums.DirectionSell.Opposite())
}

// Opposite returns Buy for any non-Buy value (the else arm), so an
// unknown value also flips to Buy. We document and verify that contract.
func TestDirection_Opposite_UnknownTreatedAsSell(t *testing.T) {
	assert.Equal(t, enums.DirectionBuy, enums.Direction("unknown").Opposite())
}

func TestDirection_JSONMarshal(t *testing.T) {
	type wrap struct {
		D enums.Direction `json:"d"`
	}
	b, err := json.Marshal(wrap{D: enums.DirectionBuy})
	require.NoError(t, err)
	assert.JSONEq(t, `{"d":"buy"}`, string(b))
}

func TestDirection_JSONUnmarshal(t *testing.T) {
	type wrap struct {
		D enums.Direction `json:"d"`
	}
	var got wrap
	require.NoError(t, json.Unmarshal([]byte(`{"d":"sell"}`), &got))
	assert.Equal(t, enums.DirectionSell, got.D)
}

package private_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/channels/private"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestBalances_Name(t *testing.T) {
	assert.Equal(t, "subaccount.5.balances", private.Balances{SubaccountID: 5}.Name())
}

func TestBalances_Decode_Happy(t *testing.T) {
	raw := json.RawMessage(`{"subaccount_id":5,"subaccount_value":"100","initial_margin":"50","maintenance_margin":"30","collaterals":[]}`)
	v, err := private.Balances{}.Decode(raw)
	require.NoError(t, err)
	bal, ok := v.(types.Balance)
	require.True(t, ok)
	assert.Equal(t, int64(5), bal.SubaccountID)
}

func TestBalances_Decode_Malformed(t *testing.T) {
	_, err := private.Balances{}.Decode([]byte(`[`))
	assert.Error(t, err)
}

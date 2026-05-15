package methods_test

import (
	"context"
	"testing"

	"github.com/amiwrpremium/go-derive/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllReferralCodes_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_all_referral_codes", []any{
		map[string]any{
			"referral_code":    "ALICE",
			"wallet":           "0x1111111111111111111111111111111111111111",
			"receiving_wallet": "0x2222222222222222222222222222222222222222",
		},
	})
	got, err := api.GetAllReferralCodes(context.Background())
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "ALICE", got[0].ReferralCode)
}

func TestGetReferralCode_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_referral_code", "ALICE")
	got, err := api.GetReferralCode(context.Background(), types.ReferralCodeQuery{Wallet: "0x1111111111111111111111111111111111111111"})
	require.NoError(t, err)
	assert.Equal(t, "ALICE", got)
}

func TestGetInviteCode_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/get_invite_code", map[string]any{
		"code":           "INVITE-ALICE",
		"remaining_uses": int64(-1),
	})
	got, err := api.GetInviteCode(context.Background(), types.InviteCodeQuery{Wallet: "0x1111111111111111111111111111111111111111"})
	require.NoError(t, err)
	assert.Equal(t, "INVITE-ALICE", got.Code)
	assert.Equal(t, int64(-1), got.RemainingUses)
}

func TestValidateInviteCode_Decode(t *testing.T) {
	api, ft := newAPI(t, false, 0)
	ft.HandleResult("public/validate_invite_code", "valid")
	got, err := api.ValidateInviteCode(context.Background(), types.ValidateInviteCodeQuery{Code: "INVITE-ALICE"})
	require.NoError(t, err)
	assert.Equal(t, "valid", got)
}

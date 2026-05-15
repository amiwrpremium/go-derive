// Package methods is the shared implementation of every JSON-RPC method
// Derive exposes. Both pkg/rest.Client and pkg/ws.Client embed *API so that
// each method is defined exactly once, parameterised by the underlying
// transport.
//
// Public methods are unauthenticated; private methods require Signer to be
// non-nil. Private methods that mutate orders also use the Domain to sign
// the per-action EIP-712 hash.
package methods

import (
	"context"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

// GetReferralCode returns the referral code currently associated
// with one wallet. Public.
//
// Pass an empty `wallet` to default to the configured signer's
// wallet (if any).
func (a *API) GetReferralCode(ctx context.Context, q types.ReferralCodeQuery) (string, error) {
	wallet := q.Wallet
	if wallet == "" && a.Signer != nil {
		wallet = a.Signer.OwnerAddress().Hex()
	}
	params := map[string]any{}
	if wallet != "" {
		params["wallet"] = wallet
	}
	var code string
	if err := a.call(ctx, "public/get_referral_code", params, &code); err != nil {
		return "", err
	}
	return code, nil
}

// GetInviteCode returns the invite code currently allocated to one
// wallet, plus its remaining-uses counter (`-1` for unlimited).
// Public.
//
// Pass an empty `wallet` to default to the configured signer's
// wallet (if any).
func (a *API) GetInviteCode(ctx context.Context, q types.InviteCodeQuery) (types.InviteCode, error) {
	wallet := q.Wallet
	if wallet == "" && a.Signer != nil {
		wallet = a.Signer.OwnerAddress().Hex()
	}
	params := map[string]any{}
	if wallet != "" {
		params["wallet"] = wallet
	}
	var resp types.InviteCode
	if err := a.call(ctx, "public/get_invite_code", params, &resp); err != nil {
		return types.InviteCode{}, err
	}
	return resp, nil
}

// ValidateInviteCode checks whether one invite code is valid and
// has remaining uses. Public.
func (a *API) ValidateInviteCode(ctx context.Context, q types.ValidateInviteCodeQuery) (string, error) {
	var status string
	if err := a.call(ctx, "public/validate_invite_code", map[string]any{
		"code": q.Code,
	}, &status); err != nil {
		return "", err
	}
	return status, nil
}

// GetAllReferralCodes returns every valid referral code on the
// configured signer's wallet. Public — but the wallet param is
// sourced from the signer if available; otherwise the engine
// applies its server-side default.
func (a *API) GetAllReferralCodes(ctx context.Context) ([]types.ReferralCodeRecord, error) {
	params := map[string]any{}
	if a.Signer != nil {
		params["wallet"] = a.Signer.OwnerAddress().Hex()
	}
	var resp []types.ReferralCodeRecord
	if err := a.call(ctx, "public/get_all_referral_codes", params, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

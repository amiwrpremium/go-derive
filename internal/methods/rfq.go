// Package methods — see collateral.go for the overview.
package methods

import (
	"context"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

// SendRFQ broadcasts a request-for-quote to market makers. Private.
func (a *API) SendRFQ(ctx context.Context, legs []types.RFQLeg, maxFee types.Decimal) (types.RFQ, error) {
	if err := a.requireSubaccount(); err != nil {
		return types.RFQ{}, err
	}
	params := map[string]any{
		"subaccount_id": a.Subaccount,
		"legs":          legs,
		"max_total_fee": maxFee,
	}
	var rfq types.RFQ
	err := a.call(ctx, "private/send_rfq", params, &rfq)
	return rfq, err
}

// PollRFQs returns the status of recent RFQs initiated by this subaccount. Private.
func (a *API) PollRFQs(ctx context.Context) ([]types.RFQ, error) {
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	var resp struct {
		RFQs []types.RFQ `json:"rfqs"`
	}
	err := a.call(ctx, "private/poll_rfqs", map[string]any{
		"subaccount_id": a.Subaccount,
	}, &resp)
	return resp.RFQs, err
}

// CancelRFQ cancels an outstanding RFQ. Private.
func (a *API) CancelRFQ(ctx context.Context, rfqID string) error {
	if err := a.requireSubaccount(); err != nil {
		return err
	}
	return a.call(ctx, "private/cancel_rfq", map[string]any{
		"subaccount_id": a.Subaccount,
		"rfq_id":        rfqID,
	}, nil)
}

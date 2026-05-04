package methods_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestSendRFQ_Happy(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/send_rfq", map[string]any{
		"rfq_id": "R1", "subaccount_id": 1, "status": "open",
		"legs": []any{}, "creation_timestamp": 1, "last_update_timestamp": 1,
	})
	rfq, err := api.SendRFQ(context.Background(), nil, types.MustDecimal("100"))
	require.NoError(t, err)
	assert.Equal(t, "R1", rfq.RFQID)
}

func TestSendRFQ_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	_, err := api.SendRFQ(context.Background(), nil, types.MustDecimal("0"))
	assert.ErrorIs(t, err, derrors.ErrSubaccountRequired)
}

func TestPollRFQs_Happy(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/poll_rfqs", map[string]any{"rfqs": []any{}})
	got, err := api.PollRFQs(context.Background())
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestPollRFQs_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	_, err := api.PollRFQs(context.Background())
	assert.ErrorIs(t, err, derrors.ErrSubaccountRequired)
}

func TestCancelRFQ_Happy(t *testing.T) {
	api, ft := newAPI(t, true, 1)
	ft.HandleResult("private/cancel_rfq", nil)
	require.NoError(t, api.CancelRFQ(context.Background(), "R1"))
	params := paramsAsMap(t, ft.LastCall().Params)
	assert.Equal(t, "R1", params["rfq_id"])
}

func TestCancelRFQ_RequiresSubaccount(t *testing.T) {
	api, _ := newAPI(t, true, 0)
	err := api.CancelRFQ(context.Background(), "R1")
	assert.ErrorIs(t, err, derrors.ErrSubaccountRequired)
}

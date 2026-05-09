package methods_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive"
)

func TestRFQExtras_AllMethods(t *testing.T) {
	api, ft := newAPI(t, true, 9)
	cases := []struct {
		name    string
		method  string
		invoke  func() (json.RawMessage, error)
		mockOut any
	}{
		{"GetRFQs", "private/get_rfqs",
			func() (json.RawMessage, error) { return api.GetRFQs(context.Background(), nil) },
			[]any{}},
		{"GetQuotes", "private/get_quotes",
			func() (json.RawMessage, error) { return api.GetQuotes(context.Background(), nil) },
			[]any{}},
		{"PollQuotes", "private/poll_quotes",
			func() (json.RawMessage, error) { return api.PollQuotes(context.Background(), nil) },
			[]any{}},
		{"SendQuote", "private/send_quote",
			func() (json.RawMessage, error) {
				return api.SendQuote(context.Background(), map[string]any{"rfq_id": "R"})
			},
			map[string]any{"quote_id": "Q"}},
		{"ExecuteQuote", "private/execute_quote",
			func() (json.RawMessage, error) {
				return api.ExecuteQuote(context.Background(), map[string]any{"quote_id": "Q"})
			},
			map[string]any{"status": "filled"}},
		{"CancelQuote", "private/cancel_quote",
			func() (json.RawMessage, error) {
				return api.CancelQuote(context.Background(), map[string]any{"quote_id": "Q"})
			},
			map[string]any{}},
		{"CancelBatchQuotes", "private/cancel_batch_quotes",
			func() (json.RawMessage, error) {
				return api.CancelBatchQuotes(context.Background(), nil)
			},
			map[string]any{}},
		{"CancelBatchRFQs", "private/cancel_batch_rfqs",
			func() (json.RawMessage, error) {
				return api.CancelBatchRFQs(context.Background(), nil)
			},
			map[string]any{}},
		{"RFQGetBestQuote", "private/rfq_get_best_quote",
			func() (json.RawMessage, error) {
				return api.RFQGetBestQuote(context.Background(), map[string]any{"rfq_id": "R"})
			},
			map[string]any{"price": "1"}},
		{"OrderQuote", "private/order_quote",
			func() (json.RawMessage, error) {
				return api.OrderQuote(context.Background(), map[string]any{"instrument_name": "BTC-PERP"})
			},
			map[string]any{}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ft.HandleResult(c.method, c.mockOut)
			raw, err := c.invoke()
			require.NoError(t, err, "method %s", c.method)
			assert.NotEmpty(t, raw)
		})
	}
}

func TestRFQExtras_RequireSigner(t *testing.T) {
	api, _ := newAPI(t, false, 0)
	checks := []func() (json.RawMessage, error){
		func() (json.RawMessage, error) { return api.GetRFQs(context.Background(), nil) },
		func() (json.RawMessage, error) { return api.GetQuotes(context.Background(), nil) },
		func() (json.RawMessage, error) { return api.SendQuote(context.Background(), nil) },
		func() (json.RawMessage, error) { return api.ExecuteQuote(context.Background(), nil) },
		func() (json.RawMessage, error) { return api.OrderQuote(context.Background(), nil) },
	}
	for _, fn := range checks {
		_, err := fn()
		assert.ErrorIs(t, err, derive.ErrUnauthorized)
	}
}

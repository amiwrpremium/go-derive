package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

func TestPage_RoundTrip(t *testing.T) {
	in := types.Page{NumPages: 5, Count: 100}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	var out types.Page
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in, out)
}

// Derive may add fields like `current_page` later — the SDK should accept
// them silently rather than fail the decode.
func TestPage_IgnoresUnknownFields(t *testing.T) {
	raw := []byte(`{"num_pages":3,"count":50,"current_page":2,"page_size":20}`)
	var p types.Page
	require.NoError(t, json.Unmarshal(raw, &p))
	assert.Equal(t, 3, p.NumPages)
	assert.Equal(t, 50, p.Count)
}

func TestPageRequest_OmitsZeroFields(t *testing.T) {
	in := types.PageRequest{}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	assert.JSONEq(t, `{}`, string(b))
}

func TestPageRequest_IncludesPopulated(t *testing.T) {
	in := types.PageRequest{Page: 2, PageSize: 25}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	assert.JSONEq(t, `{"page":2,"page_size":25}`, string(b))
}

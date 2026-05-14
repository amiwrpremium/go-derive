package types_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/types"
)

// makePages returns a fetch closure that emits pre-baked pages.
// pages[i] is the slice returned for request page i+1.
// numPages is what every response carries in the Page envelope (so the
// driver knows when to stop). Use numPages=0 for the "empty result"
// edge case.
func makePages(pages [][]int, numPages int) func(context.Context, types.PageRequest) ([]int, types.Page, error) {
	return func(_ context.Context, p types.PageRequest) ([]int, types.Page, error) {
		idx := p.Page - 1
		if idx < 0 || idx >= len(pages) {
			return nil, types.Page{NumPages: numPages}, nil
		}
		return pages[idx], types.Page{NumPages: numPages, Count: countAll(pages)}, nil
	}
}

func countAll(pages [][]int) int {
	n := 0
	for _, p := range pages {
		n += len(p)
	}
	return n
}

func TestPaginate_EmptyResult(t *testing.T) {
	got, err := types.Paginate(context.Background(), types.PaginateOptions{},
		makePages([][]int{}, 0))
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestPaginate_SinglePage(t *testing.T) {
	got, err := types.Paginate(context.Background(), types.PaginateOptions{},
		makePages([][]int{{1, 2, 3}}, 1))
	require.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, got)
}

func TestPaginate_MultiPage(t *testing.T) {
	got, err := types.Paginate(context.Background(), types.PaginateOptions{},
		makePages([][]int{{1, 2}, {3, 4}, {5}}, 3))
	require.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, got)
}

func TestPaginate_MaxItemsTruncates(t *testing.T) {
	got, err := types.Paginate(context.Background(),
		types.PaginateOptions{MaxItems: 3},
		makePages([][]int{{1, 2}, {3, 4}, {5, 6}}, 3))
	require.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, got)
}

func TestPaginate_MaxPagesCaps(t *testing.T) {
	got, err := types.Paginate(context.Background(),
		types.PaginateOptions{MaxPages: 2},
		makePages([][]int{{1, 2}, {3, 4}, {5, 6}}, 3))
	require.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4}, got)
}

func TestPaginate_PageSizeIsThreadedThrough(t *testing.T) {
	var seenSizes []int
	fetch := func(_ context.Context, p types.PageRequest) ([]int, types.Page, error) {
		seenSizes = append(seenSizes, p.PageSize)
		if p.Page == 1 {
			return []int{1, 2}, types.Page{NumPages: 2}, nil
		}
		return []int{3}, types.Page{NumPages: 2}, nil
	}
	_, err := types.Paginate(context.Background(),
		types.PaginateOptions{PageSize: 50}, fetch)
	require.NoError(t, err)
	assert.Equal(t, []int{50, 50}, seenSizes)
}

func TestPaginate_PropagatesFetchError(t *testing.T) {
	boom := errors.New("boom")
	calls := 0
	fetch := func(_ context.Context, _ types.PageRequest) ([]int, types.Page, error) {
		calls++
		if calls == 2 {
			return nil, types.Page{}, boom
		}
		return []int{calls}, types.Page{NumPages: 3}, nil
	}
	got, err := types.Paginate(context.Background(), types.PaginateOptions{}, fetch)
	require.ErrorIs(t, err, boom)
	assert.Nil(t, got, "partial results are discarded on error")
}

func TestPaginate_RespectsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	fetch := func(_ context.Context, _ types.PageRequest) ([]int, types.Page, error) {
		calls++
		if calls == 1 {
			cancel()
			return []int{1, 2}, types.Page{NumPages: 5}, nil
		}
		t.Fatal("fetch called after ctx cancelled")
		return nil, types.Page{}, nil
	}
	got, err := types.Paginate(ctx, types.PaginateOptions{}, fetch)
	require.ErrorIs(t, err, context.Canceled)
	assert.Nil(t, got)
}

func TestPaginate_StopsOnEmptyMidStream(t *testing.T) {
	// Some endpoints incorrectly report NumPages=100 but return an empty
	// page after the actual data is exhausted. The driver should treat
	// the first empty page as a terminator, not loop forever.
	got, err := types.Paginate(context.Background(), types.PaginateOptions{},
		makePages([][]int{{1, 2}, {}}, 100))
	require.NoError(t, err)
	assert.Equal(t, []int{1, 2}, got)
}

func TestPaginate_NumPagesZeroFallsBackToEmptyTerminator(t *testing.T) {
	// NumPages=0 in the response means the engine didn't compute totals;
	// the driver should keep paging until an empty page arrives.
	got, err := types.Paginate(context.Background(), types.PaginateOptions{},
		makePages([][]int{{1, 2}, {3, 4}, {}}, 0))
	require.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4}, got)
}

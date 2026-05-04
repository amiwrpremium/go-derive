package jsonrpc_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/internal/jsonrpc"
)

func TestNewIDGen_FirstNextIsOne(t *testing.T) {
	g := jsonrpc.NewIDGen()
	assert.Equal(t, uint64(1), g.Next())
}

func TestIDGen_StrictlyIncreasing(t *testing.T) {
	g := jsonrpc.NewIDGen()
	require.Equal(t, uint64(1), g.Next())
	require.Equal(t, uint64(2), g.Next())
	require.Equal(t, uint64(3), g.Next())
}

func TestIDGen_ConcurrentUniqueness(t *testing.T) {
	g := jsonrpc.NewIDGen()
	const n = 1000
	const goroutines = 8
	var wg sync.WaitGroup
	results := make([][]uint64, goroutines)
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			out := make([]uint64, n)
			for j := 0; j < n; j++ {
				out[j] = g.Next()
			}
			results[i] = out
		}(i)
	}
	wg.Wait()

	seen := map[uint64]struct{}{}
	for _, r := range results {
		for _, id := range r {
			_, dup := seen[id]
			require.False(t, dup, "duplicate id %d", id)
			seen[id] = struct{}{}
		}
	}
	assert.Equal(t, goroutines*n, len(seen))
}

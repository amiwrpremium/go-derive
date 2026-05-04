package auth_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amiwrpremium/go-derive/pkg/auth"
)

func TestNonceGen_StrictlyMonotonic(t *testing.T) {
	g := auth.NewNonceGen()
	prev := g.Next()
	for i := 0; i < 1000; i++ {
		n := g.Next()
		require.Greater(t, n, prev, "iteration %d", i)
		prev = n
	}
}

func TestNonceGen_ConcurrentUniqueness(t *testing.T) {
	g := auth.NewNonceGen()
	const goroutines = 16
	const perG = 200
	results := make([][]uint64, goroutines)
	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			out := make([]uint64, perG)
			for j := 0; j < perG; j++ {
				out[j] = g.Next()
			}
			results[i] = out
		}(i)
	}
	wg.Wait()

	seen := map[uint64]struct{}{}
	for _, r := range results {
		for _, n := range r {
			_, dup := seen[n]
			require.False(t, dup, "duplicate nonce %d", n)
			seen[n] = struct{}{}
		}
	}
	assert.Equal(t, goroutines*perG, len(seen))
}

func TestNonceGen_NewIsNonZeroFirst(t *testing.T) {
	g := auth.NewNonceGen()
	assert.NotZero(t, g.Next())
}

// Package methods is the shared implementation of every JSON-RPC method
// Derive exposes.
package methods

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

// instrumentMeta is the slice of an [types.Instrument] that signed
// actions need: the on-chain base-asset address and per-asset sub-id.
// These are not sent on the wire — they feed the EIP-712
// TradeModuleData / RFQQuoteModuleData payload that the SDK hashes
// before signing.
type instrumentMeta struct {
	Asset types.Address
	SubID uint64
}

// instrumentCache is a process-lifetime map of instrument_name → on-chain
// metadata. The cache exists so callers can place orders / send RFQs by
// instrument name without first having to call public/get_instrument and
// thread the resulting BaseAssetAddress / BaseAssetSubID through every
// signed action.
//
// Three population paths are supported, in increasing eagerness:
//
//  1. Lazy. When a signed action is built with Asset.IsZero() the cache
//     is consulted; on miss [API.GetInstrument] is called and the result
//     is cached.
//  2. Side-effect. [API.GetInstrument], [API.GetInstruments], and
//     [API.GetAllInstruments] all populate the cache as a side-effect of
//     a successful response — regardless of who called them. A caller
//     listing instruments for a UI warms the cache for free.
//  3. Eager. [API.PreloadInstruments] and [API.PreloadAllInstruments]
//     fetch and cache up front. The client constructors expose a
//     WithInstrumentPreload option that kicks off the latter in a
//     background goroutine immediately after construction.
//
// Caller-supplied Asset/SubID on a signed action are always honoured —
// the cache is only consulted when the caller leaves them zero.
type instrumentCache struct {
	mu sync.RWMutex
	m  map[string]instrumentMeta
}

func newInstrumentCache() *instrumentCache {
	return &instrumentCache{m: make(map[string]instrumentMeta)}
}

func (c *instrumentCache) get(name string) (instrumentMeta, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.m[name]
	return v, ok
}

// populateOne extracts on-chain metadata from an [types.Instrument] and
// stores it under the instrument's name. Returns the stored value plus
// ok=true on success; ok=false when the instrument lacked the required
// fields (the cache is left untouched in that case).
func (c *instrumentCache) populateOne(inst types.Instrument) (instrumentMeta, bool) {
	if inst.Name == "" || inst.BaseAssetAddress.IsZero() {
		return instrumentMeta{}, false
	}
	sub, err := strconv.ParseUint(inst.BaseAssetSubID, 10, 64)
	if err != nil {
		return instrumentMeta{}, false
	}
	meta := instrumentMeta{Asset: inst.BaseAssetAddress, SubID: sub}
	c.mu.Lock()
	c.m[inst.Name] = meta
	c.mu.Unlock()
	return meta, true
}

// populate loops populateOne over a slice. Entries that lack on-chain
// metadata are silently skipped.
func (c *instrumentCache) populate(insts []types.Instrument) {
	for _, inst := range insts {
		c.populateOne(inst)
	}
}

func (c *instrumentCache) invalidate(name string) {
	c.mu.Lock()
	delete(c.m, name)
	c.mu.Unlock()
}

// instCache returns the API's cache, lazily creating it on first use.
// Self-initialisation lets transport-layer constructors leave the field
// zero — see [API.cacheStorage].
func (a *API) instCache() *instrumentCache {
	a.instrumentCacheOnce.Do(func() {
		a.instrumentCacheStorage = newInstrumentCache()
	})
	return a.instrumentCacheStorage
}

// resolveInstrument returns the on-chain metadata for one instrument,
// fetching from public/get_instrument on cache miss. Used by signed
// actions when the caller leaves Asset zero on the input.
func (a *API) resolveInstrument(ctx context.Context, name string) (instrumentMeta, error) {
	if meta, ok := a.instCache().get(name); ok {
		return meta, nil
	}
	inst, err := a.GetInstrument(ctx, name)
	if err != nil {
		return instrumentMeta{}, err
	}
	if meta, ok := a.instCache().populateOne(inst); ok {
		return meta, nil
	}
	return instrumentMeta{}, fmt.Errorf(
		"methods: instrument %q has no on-chain metadata (base_asset_address or base_asset_sub_id missing); pass Asset/SubID explicitly",
		name,
	)
}

// InvalidateInstrumentCache removes one entry from the instrument
// metadata cache. Use this if the engine reports a stale-instrument
// error and a refresh is needed; the next signed action against the
// instrument will refetch via public/get_instrument.
func (a *API) InvalidateInstrumentCache(name string) {
	a.instCache().invalidate(name)
}

// PreloadInstruments fetches every live instrument in the supplied
// currencies and populates the metadata cache. Each currency × kind
// (perp, option, erc20) is fetched in a single
// public/get_instruments call.
//
// Returns the first error encountered (subsequent fetches are
// abandoned). The cache reflects whatever fetches completed before
// the error.
func (a *API) PreloadInstruments(ctx context.Context, currencies ...string) error {
	kinds := []enums.InstrumentType{
		enums.InstrumentTypePerp,
		enums.InstrumentTypeOption,
		enums.InstrumentTypeERC20,
	}
	for _, ccy := range currencies {
		for _, kind := range kinds {
			if _, err := a.GetInstruments(ctx, ccy, kind); err != nil {
				return fmt.Errorf("methods: preload %s %s: %w", ccy, kind, err)
			}
		}
	}
	return nil
}

// PreloadAllInstruments paginates every live instrument across all
// currencies via public/get_all_instruments and populates the cache.
// Calls GetAllInstruments once per kind (perp, option, erc20) and
// pages until exhaustion.
//
// Expired instruments are not included — they can't be traded, so
// caching them is wasteful.
//
// Returns the first error encountered.
func (a *API) PreloadAllInstruments(ctx context.Context) error {
	kinds := []enums.InstrumentType{
		enums.InstrumentTypePerp,
		enums.InstrumentTypeOption,
		enums.InstrumentTypeERC20,
	}
	for _, kind := range kinds {
		page := types.PageRequest{Page: 1, PageSize: 1000}
		for {
			_, info, err := a.GetAllInstruments(ctx, kind, false, page)
			if err != nil {
				return fmt.Errorf("methods: preload all %s: %w", kind, err)
			}
			if info.NumPages == 0 || page.Page >= info.NumPages {
				break
			}
			page.Page++
		}
	}
	return nil
}

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

// GetInstruments lists active instruments matching the filter. Public.
//
// Derive returns the result as a bare JSON array of instrument objects.
// As a side effect, the returned instruments' on-chain metadata is
// stored in the SDK's instrument cache so subsequent signed actions on
// any of them skip the per-instrument get_instrument lookup.
func (a *API) GetInstruments(ctx context.Context, q types.InstrumentsQuery) ([]types.Instrument, error) {
	params := map[string]any{}
	if q.Currency != "" {
		params["currency"] = q.Currency
	}
	if q.Kind != "" {
		params["instrument_type"] = q.Kind
	}
	params["expired"] = false
	var insts []types.Instrument
	if err := a.call(ctx, "public/get_instruments", params, &insts); err != nil {
		return nil, err
	}
	a.instCache().populate(insts)
	return insts, nil
}

// GetInstrument fetches one instrument by name. Public.
//
// As a side effect, the returned instrument's on-chain metadata is
// stored in the SDK's instrument cache so subsequent signed actions on
// it skip the lookup.
func (a *API) GetInstrument(ctx context.Context, q types.InstrumentQuery) (types.Instrument, error) {
	var inst types.Instrument
	err := a.call(ctx, "public/get_instrument", map[string]any{"instrument_name": q.Name}, &inst)
	if err == nil {
		a.instCache().populateOne(inst)
	}
	return inst, err
}

// GetTicker fetches the public ticker for one instrument. Public.
func (a *API) GetTicker(ctx context.Context, q types.TickerQuery) (types.Ticker, error) {
	var t types.Ticker
	err := a.call(ctx, "public/get_ticker", map[string]any{"instrument_name": q.Name}, &t)
	return t, err
}

// GetPublicTradeHistory returns recent trades on the instrument. Public.
//
// All fields on the query are optional and AND together. Pass an
// instrument name to scope to one market, a currency to scope to one
// underlying, or both — leave fields zero to ask the engine for an
// unfiltered tape.
func (a *API) GetPublicTradeHistory(ctx context.Context, q types.PublicTradeHistoryQuery, page types.PageRequest) ([]types.Trade, types.Page, error) {
	params := map[string]any{}
	if q.InstrumentName != "" {
		params["instrument_name"] = q.InstrumentName
	}
	if q.Currency != "" {
		params["currency"] = q.Currency
	}
	if q.InstrumentType != "" {
		params["instrument_type"] = q.InstrumentType
	}
	if q.SubaccountID != 0 {
		params["subaccount_id"] = q.SubaccountID
	}
	if q.TradeID != "" {
		params["trade_id"] = q.TradeID
	}
	if q.TxHash != "" {
		params["tx_hash"] = q.TxHash
	}
	if q.TxStatus != "" {
		params["tx_status"] = q.TxStatus
	}
	if !q.FromTimestamp.Time().IsZero() {
		params["from_timestamp"] = q.FromTimestamp.Millis()
	}
	if !q.ToTimestamp.Time().IsZero() {
		params["to_timestamp"] = q.ToTimestamp.Millis()
	}
	if page.Page > 0 {
		params["page"] = page.Page
	}
	if page.PageSize > 0 {
		params["page_size"] = page.PageSize
	}
	var resp struct {
		Trades     []types.Trade `json:"trades"`
		Pagination types.Page    `json:"pagination"`
	}
	if err := a.call(ctx, "public/get_trade_history", params, &resp); err != nil {
		return nil, types.Page{}, err
	}
	return resp.Trades, resp.Pagination, nil
}

// GetTime returns the server clock in milliseconds. Public.
func (a *API) GetTime(ctx context.Context) (int64, error) {
	var t int64
	err := a.call(ctx, "public/get_time", map[string]any{}, &t)
	return t, err
}

// GetCurrencies returns the list of supported currency names. Public.
//
// Derive's `public/get_all_currencies` result is a bare JSON array of
// rich currency objects (margin parameters, manager addresses, etc.);
// this method extracts the `currency` name field from each. Callers
// that need the full object should call the raw transport directly.
func (a *API) GetCurrencies(ctx context.Context) ([]string, error) {
	var raw []struct {
		Currency string `json:"currency"`
	}
	if err := a.call(ctx, "public/get_all_currencies", map[string]any{}, &raw); err != nil {
		return nil, err
	}
	out := make([]string, 0, len(raw))
	for _, c := range raw {
		out = append(out, c.Currency)
	}
	return out, nil
}

// GetCurrency returns the per-asset margin parameters, manager
// addresses, and protocol-asset addresses for one underlying
// currency. Public.
//
// Counterpart to the plural [API.GetCurrencies] (which returns just
// the currency names).
func (a *API) GetCurrency(ctx context.Context, q types.CurrencyQuery) (types.Currency, error) {
	var c types.Currency
	if err := a.call(ctx, "public/get_currency", map[string]any{"currency": q.Currency}, &c); err != nil {
		return types.Currency{}, err
	}
	return c, nil
}

// GetAllInstruments lists every instrument matching the supplied
// kind, paginated. Public.
//
// Distinct from [API.GetInstruments] — `public/get_instruments` is
// for the live, currency-filtered list a UI uses; this method backs
// `public/get_all_instruments`, which paginates across all
// currencies and can include expired instruments via `includeExpired`.
func (a *API) GetAllInstruments(ctx context.Context, q types.AllInstrumentsQuery, page types.PageRequest) ([]types.Instrument, types.Page, error) {
	params := map[string]any{
		"instrument_type": q.Kind,
		"expired":         q.IncludeExpired,
	}
	if page.Page > 0 {
		params["page"] = page.Page
	}
	if page.PageSize > 0 {
		params["page_size"] = page.PageSize
	}
	var resp struct {
		Instruments []types.Instrument `json:"instruments"`
		Pagination  types.Page         `json:"pagination"`
	}
	if err := a.call(ctx, "public/get_all_instruments", params, &resp); err != nil {
		return nil, types.Page{}, err
	}
	a.instCache().populate(resp.Instruments)
	return resp.Instruments, resp.Pagination, nil
}

// GetTickers returns the ticker snapshot keyed by instrument name
// for every instrument matching the filter. Public.
//
// `instrumentType` is required; `currency` and `expiryDate` (in the
// engine's YYYYMMDD numeric form) are both required for option
// queries and ignored for perpetuals / ERC-20s. Pass empty / zero
// for either when not applicable.
//
// Each value is a [types.InstrumentTickerSlim] — the bare per-
// instrument compact-wire payload — not the WS subscription envelope
// [types.TickerSlim], which wraps the same payload with an outer
// `{timestamp, instrument_ticker}` shape that this REST endpoint
// does not emit.
func (a *API) GetTickers(ctx context.Context, q types.TickersQuery) (map[string]types.InstrumentTickerSlim, error) {
	params := map[string]any{
		"instrument_type": q.InstrumentType,
	}
	if q.Currency != "" {
		params["currency"] = q.Currency
	}
	if q.ExpiryDate != 0 {
		params["expiry_date"] = q.ExpiryDate
	}
	var resp struct {
		Tickers map[string]types.InstrumentTickerSlim `json:"tickers"`
	}
	if err := a.call(ctx, "public/get_tickers", params, &resp); err != nil {
		return nil, err
	}
	return resp.Tickers, nil
}

// GetOptionSettlementPrices returns the per-expiry settlement prices
// for one currency's option market. Public.
//
// Pre-settlement entries return Price as the zero-value Decimal
// (the wire field is null until the expiry settles on chain).
func (a *API) GetOptionSettlementPrices(ctx context.Context, q types.OptionSettlementPricesQuery) ([]types.OptionSettlementPrice, error) {
	var resp struct {
		Expiries []types.OptionSettlementPrice `json:"expiries"`
	}
	if err := a.call(ctx, "public/get_option_settlement_prices", map[string]any{"currency": q.Currency}, &resp); err != nil {
		return nil, err
	}
	return resp.Expiries, nil
}

// GetLiveIncidents returns the list of platform incidents currently
// in progress. Public.
//
// Empty list means no active incidents.
func (a *API) GetLiveIncidents(ctx context.Context) ([]types.Incident, error) {
	var resp struct {
		Incidents []types.Incident `json:"incidents"`
	}
	if err := a.call(ctx, "public/get_live_incidents", map[string]any{}, &resp); err != nil {
		return nil, err
	}
	return resp.Incidents, nil
}

// GetAllStatistics returns the per-(currency, instrument_type)
// aggregate of rolling 24h and all-time statistics across every
// instrument. Public.
//
// Optional `endTime` (Unix seconds) — pass 0 for the engine's
// default (now).
func (a *API) GetAllStatistics(ctx context.Context, q types.AllStatisticsQuery) ([]types.AggregateStatistics, error) {
	params := map[string]any{}
	if q.EndTime > 0 {
		params["end_time"] = q.EndTime
	}
	var resp []types.AggregateStatistics
	if err := a.call(ctx, "public/all_statistics", params, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetAllUserStatistics returns the per-wallet trading statistics
// for every wallet, optionally scoped to a `[0, endTime]` window
// (Unix seconds). Pass zero for the engine's "current time" default.
// Public.
//
// The engine accepts additional optional filters (currency,
// instrument_name, is_rfq, start_time) per docs.derive.xyz that
// this overload does not surface; if you need them, route through
// the lower-level call helpers.
func (a *API) GetAllUserStatistics(ctx context.Context, q types.AllUserStatisticsQuery) ([]types.UserStatistics, error) {
	params := map[string]any{}
	if q.EndTimeSec != 0 {
		params["end_time"] = q.EndTimeSec
	}
	var resp []types.UserStatistics
	if err := a.call(ctx, "public/all_user_statistics", params, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetUserStatistics returns the trading statistics for one wallet.
// Public.
//
// The engine accepts additional optional filters (currency,
// end_time, instrument_name, is_rfq, start_time) per
// docs.derive.xyz that this overload does not surface.
func (a *API) GetUserStatistics(ctx context.Context, q types.UserStatisticsQuery) (types.UserStatistics, error) {
	params := map[string]any{
		"wallet": q.Wallet,
	}
	var resp types.UserStatistics
	if err := a.call(ctx, "public/user_statistics", params, &resp); err != nil {
		return types.UserStatistics{}, err
	}
	return resp, nil
}

// GetAsset fetches one Asset record by name. Public.
//
// Distinct from [API.GetInstrument]: an asset is the on-chain
// ERC-1155 entity (token id, decimals, address); an instrument adds
// the orderbook / pricing surface.
func (a *API) GetAsset(ctx context.Context, q types.AssetQuery) (types.Asset, error) {
	var resp types.Asset
	if err := a.call(ctx, "public/get_asset", map[string]any{"asset_name": q.Name}, &resp); err != nil {
		return types.Asset{}, err
	}
	return resp, nil
}

// GetAssets lists Asset records matching the supplied filter.
// Public.
//
// Per docs.derive.xyz all three arguments are required:
// `assetType` ("erc20", "option", "perp"), `currency` (e.g.
// "ETH"), and `expired` (include expired assets).
func (a *API) GetAssets(ctx context.Context, q types.AssetsQuery) ([]types.Asset, error) {
	params := map[string]any{
		"asset_type": q.AssetType,
		"currency":   q.Currency,
		"expired":    q.Expired,
	}
	var resp []types.Asset
	if err := a.call(ctx, "public/get_assets", params, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetBridgeBalances lists every bridge / cross-chain balance the
// engine tracks. Public.
func (a *API) GetBridgeBalances(ctx context.Context) ([]types.BridgeBalance, error) {
	var resp []types.BridgeBalance
	if err := a.call(ctx, "public/get_bridge_balances", map[string]any{}, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetStDRVSnapshots returns one wallet's staked-DRV balance snapshots
// over a time window. Public.
//
// Note: timestamps on this endpoint are in seconds, not
// milliseconds.
func (a *API) GetStDRVSnapshots(ctx context.Context, q types.STDRVSnapshotsQuery) (types.StDRVSnapshots, error) {
	params := map[string]any{
		"wallet": q.Wallet,
	}
	if q.FromSec != 0 {
		params["from_sec"] = q.FromSec
	}
	if q.ToSec != 0 {
		params["to_sec"] = q.ToSec
	}
	var resp types.StDRVSnapshots
	if err := a.call(ctx, "public/get_stdrv_snapshots", params, &resp); err != nil {
		return types.StDRVSnapshots{}, err
	}
	return resp, nil
}

// GetDescendantTree returns the referee tree rooted at one wallet
// (or invite code). Public.
//
// The `Descendants` field is preserved as raw JSON because the wire
// shape is recursive; decode further at the call site.
func (a *API) GetDescendantTree(ctx context.Context, q types.DescendantTreeQuery) (types.DescendantTree, error) {
	var resp types.DescendantTree
	if err := a.call(ctx, "public/get_descendant_tree", map[string]any{
		"wallet_or_invite_code": q.WalletOrInviteCode,
	}, &resp); err != nil {
		return types.DescendantTree{}, err
	}
	return resp, nil
}

// GetTreeRoots returns every root wallet (top-of-tree referrer) the
// engine tracks. Public.
//
// The `Roots` field is preserved as raw JSON because the inner
// shape varies per program.
func (a *API) GetTreeRoots(ctx context.Context) (types.TreeRoots, error) {
	var resp types.TreeRoots
	if err := a.call(ctx, "public/get_tree_roots", map[string]any{}, &resp); err != nil {
		return types.TreeRoots{}, err
	}
	return resp, nil
}

// MarginWatch calculates the mark-to-market and maintenance-margin
// snapshot for one subaccount. Public.
//
// `subaccountID` is required. `forceOnchain` forces an on-chain
// balance fetch (slower, freshest data); `isDelayedLiquidation`
// lowers the maintenance-margin requirement under a
// delayed-liquidation grace period — only valid when the
// subaccount has delayed liquidation enabled.
//
// Note: this is the RPC counterpart to the platform-wide
// `margin_watch` WebSocket channel. The RPC returns a single
// snapshot for one subaccount; the channel emits a stream of
// at-risk subaccounts engine-wide.
func (a *API) MarginWatch(ctx context.Context, q types.MarginWatchQuery) (types.MarginSnapshot, error) {
	params := map[string]any{
		"subaccount_id": q.SubaccountID,
	}
	if q.ForceOnchain {
		params["force_onchain"] = true
	}
	if q.IsDelayedLiquidation {
		params["is_delayed_liquidation"] = true
	}
	var resp types.MarginSnapshot
	if err := a.call(ctx, "public/margin_watch", params, &resp); err != nil {
		return types.MarginSnapshot{}, err
	}
	return resp, nil
}

// GetStatistics returns rolling 24-hour and all-time statistics for
// one instrument: volume, premium volume, fees, trades count, plus
// total open interest. Public.
//
// Currency and EndTime on the query are optional. EndTime, when
// non-zero, anchors the rollup window to a Unix-seconds timestamp;
// otherwise the engine uses "now".
func (a *API) GetStatistics(ctx context.Context, q types.StatisticsQuery) (types.Statistics, error) {
	params := map[string]any{"instrument_name": q.InstrumentName}
	if q.Currency != "" {
		params["currency"] = q.Currency
	}
	if q.EndTime != 0 {
		params["end_time"] = q.EndTime
	}
	var resp types.Statistics
	if err := a.call(ctx, "public/statistics", params, &resp); err != nil {
		return types.Statistics{}, err
	}
	return resp, nil
}

// Shared implementation of every JSON-RPC method Derive exposes lives in
// this file. Both [RestClient] and [WsClient] embed *API so that each
// method is defined exactly once, parameterised by the underlying
// transport.
//
// Public methods are unauthenticated; private methods require [Signer] to
// be non-nil. Private methods that mutate orders also use the EIP-712
// [Domain] to sign the per-action hash.

package derive

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/amiwrpremium/go-derive/internal/transport"
	"github.com/ethereum/go-ethereum/common"
	"time"
)

// API is a transport-agnostic facade that holds the ambient configuration
// (signer, subaccount id, EIP-712 domain) used by the methods defined in this
// package. Construct it once per client.
type API struct {
	T      transport.Transport
	Signer Signer
	Domain Domain
	// Subaccount is the default subaccount id used by private endpoints. It
	// can be 0 for public-only clients.
	Subaccount int64
	// Nonces is the source of action nonces; one per client is fine.
	Nonces *NonceGen
	// SignatureExpiry is added to time.Now() to populate signature_expiry_sec
	// on signed actions.
	SignatureExpiry int64

	// tradeModule is the on-chain TradeModule contract — used as
	// Action.Module when signing order placement. Set via SetTradeModule.
	tradeModule common.Address
}

// requireSigner returns ErrUnauthorized if no signer is configured.
func (a *API) requireSigner() error {
	if a.Signer == nil {
		return ErrUnauthorized
	}
	return nil
}

// requireSubaccount returns ErrSubaccountRequired when the action needs one.
func (a *API) requireSubaccount() error {
	if a.Subaccount == 0 {
		return ErrSubaccountRequired
	}
	return nil
}

// call is a shortcut for the common case of one method, one params, one out.
// It also re-wraps a transport-level [transport.JSONRPCError] into a
// public [APIError] so callers receive the rich-typed sentinel-aware
// error.
func (a *API) call(ctx context.Context, method string, params, out any) error {
	return wrapTransportError(a.T.Call(ctx, method, params, out))
}

// wrapTransportError converts a transport-layer `*transport.JSONRPCError`
// into a public `*APIError`. Errors of any other type pass through
// unchanged. This sits at the methods/transport boundary and is the reason
// transport can stay free of an import on `pkg/errors` — closing the
// otherwise-inevitable `pkg/errors → transport → pkg/errors` cycle once
// pkg/errors lifts to root.
func wrapTransportError(err error) error {
	if rpcErr, ok := err.(*transport.JSONRPCError); ok {
		return &APIError{
			Code:    rpcErr.Code,
			Message: rpcErr.Message,
			Data:    rpcErr.Data,
		}
	}
	return err
}

// GetCollateral returns the collateral breakdown for the subaccount. Private.
func (a *API) GetCollateral(ctx context.Context) ([]Collateral, error) {
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	var resp struct {
		Collaterals []Collateral `json:"collaterals"`
	}
	err := a.call(ctx, "private/get_collaterals", map[string]any{
		"subaccount_id": a.Subaccount,
	}, &resp)
	return resp.Collaterals, err
}

// GetFundingRateHistory returns historical funding rate prints for one
// perpetual instrument over the requested window.
//
// Required params: `instrument_name`. Optional: `start_timestamp`,
// `end_timestamp`, `period`. Public.
func (a *API) GetFundingRateHistory(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	var raw json.RawMessage
	err := a.call(ctx, "public/get_funding_rate_history", params, &raw)
	return raw, err
}

// GetPerpImpactTWAP returns the time-weighted average impact price for
// one currency's perpetual book over the requested window.
//
// Required params: `currency`, `start_time`, `end_time`. Public.
func (a *API) GetPerpImpactTWAP(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	var raw json.RawMessage
	err := a.call(ctx, "public/get_perp_impact_twap", params, &raw)
	return raw, err
}

// GetPublicMargin runs Derive's risk-engine margin calculation against a
// user-supplied set of `simulated_collaterals` and `simulated_positions`,
// returning the resulting margin requirement.
//
// Required params: `simulated_collaterals`, `simulated_positions`,
// `margin_type` ("PM" / "PM2" / "SM"). Public — no signer required.
func (a *API) GetPublicMargin(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	var raw json.RawMessage
	err := a.call(ctx, "public/get_margin", params, &raw)
	return raw, err
}

// GetLatestSignedFeeds returns the latest oracle signed-feed snapshot.
//
// Optional params: `currency`. Pass nil to get every currency the venue
// publishes. Public.
func (a *API) GetLatestSignedFeeds(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if params == nil {
		params = map[string]any{}
	}
	var raw json.RawMessage
	err := a.call(ctx, "public/get_latest_signed_feeds", params, &raw)
	return raw, err
}

// GetSpotFeedHistory returns historical oracle spot prices for one
// currency over the requested window at the given period.
//
// Required params: `currency`, `period`, `start_timestamp`,
// `end_timestamp`. Public.
func (a *API) GetSpotFeedHistory(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	var raw json.RawMessage
	err := a.call(ctx, "public/get_spot_feed_history", params, &raw)
	return raw, err
}

// GetStatistics returns rolling 24h volume / OI / volatility statistics
// for one instrument.
//
// Required params: `instrument_name`. Public.
func (a *API) GetStatistics(ctx context.Context, instrument string) (json.RawMessage, error) {
	var raw json.RawMessage
	err := a.call(ctx, "public/statistics", map[string]any{"instrument_name": instrument}, &raw)
	return raw, err
}

// GetTransaction returns the on-chain status / receipt of one Derive
// transaction by its transaction id.
//
// Required params: `transaction_id`. Public.
func (a *API) GetTransaction(ctx context.Context, transactionID string) (json.RawMessage, error) {
	var raw json.RawMessage
	err := a.call(ctx, "public/get_transaction", map[string]any{"transaction_id": transactionID}, &raw)
	return raw, err
}

// GetPublicOptionSettlementHistory returns the network-wide option
// settlement history.
//
// Optional params: pagination. Pass nil for the default range. Public.
func (a *API) GetPublicOptionSettlementHistory(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if params == nil {
		params = map[string]any{}
	}
	var raw json.RawMessage
	err := a.call(ctx, "public/get_option_settlement_history", params, &raw)
	return raw, err
}

// GetAccount returns wallet-level account information for the signer.
//
// No params. Private.
func (a *API) GetAccount(ctx context.Context) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	params := map[string]any{"wallet": a.Signer.Owner().Hex()}
	var raw json.RawMessage
	err := a.call(ctx, "private/get_account", params, &raw)
	return raw, err
}

// GetMargin returns the live margin breakdown for the configured
// subaccount.
//
// No params. Private.
func (a *API) GetMargin(ctx context.Context) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/get_margin", map[string]any{"subaccount_id": a.Subaccount}, &raw)
	return raw, err
}

// GetFundingHistory returns funding payments received / paid by the
// configured subaccount.
//
// Optional params: `start_timestamp`, `end_timestamp`, `instrument_name`,
// pagination. Private.
func (a *API) GetFundingHistory(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/get_funding_history", params, &raw)
	return raw, err
}

// GetLiquidationHistory returns the configured subaccount's past
// liquidation events.
//
// Optional params: `start_timestamp`, `end_timestamp`, pagination. Private.
func (a *API) GetLiquidationHistory(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/get_liquidation_history", params, &raw)
	return raw, err
}

// GetOptionSettlementHistory returns the configured subaccount's past
// option-settlement events.
//
// Optional params: `start_timestamp`, `end_timestamp`, pagination. Private.
func (a *API) GetOptionSettlementHistory(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/get_option_settlement_history", params, &raw)
	return raw, err
}

// GetSubaccountValueHistory returns the equity-curve series for the
// configured subaccount.
//
// Required params: `period`, `start_timestamp`, `end_timestamp`. Private.
func (a *API) GetSubaccountValueHistory(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/get_subaccount_value_history", params, &raw)
	return raw, err
}

// GetERC20TransferHistory returns deposit / withdrawal-style ERC-20
// transfers attributed to the configured subaccount.
func (a *API) GetERC20TransferHistory(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/get_erc20_transfer_history", params, &raw)
	return raw, err
}

// GetInterestHistory returns the configured subaccount's interest charges
// and rebates.
func (a *API) GetInterestHistory(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/get_interest_history", params, &raw)
	return raw, err
}

// ExpiredAndCancelledHistory returns the configured subaccount's expired
// and cancelled orders.
func (a *API) ExpiredAndCancelledHistory(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/expired_and_cancelled_history", params, &raw)
	return raw, err
}

// GetMMPConfig returns the active market-maker-protection config for the
// configured subaccount.
func (a *API) GetMMPConfig(ctx context.Context) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/get_mmp_config", map[string]any{"subaccount_id": a.Subaccount}, &raw)
	return raw, err
}

// GetNotifications returns the wallet's notification feed.
func (a *API) GetNotifications(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok && a.Subaccount != 0 {
		params["subaccount_id"] = a.Subaccount
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/get_notifications", params, &raw)
	return raw, err
}

// UpdateNotifications marks one or more notifications as seen / dismissed.
//
// Required params: `notification_ids` ([]int) and `status`. Private.
func (a *API) UpdateNotifications(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/update_notifications", params, &raw)
	return raw, err
}

// Replace cancels one outstanding order and submits a replacement in a
// single round trip — the standard maker pattern for re-pricing without a
// race against the matching engine.
//
// Params should include `order_id_to_cancel` and the same fields PlaceOrder
// would take. The full param shape is documented at docs.derive.xyz.
//
// Private; requires signer + subaccount.
func (a *API) Replace(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/replace", params, &raw)
	return raw, err
}

// OrderDebug previews an order without submitting it — Derive returns the
// validated request and any synthetic fees / margin impacts the engine
// computes. Use this to sanity-check signed payloads in CI.
//
// Params mirror PlaceOrder. Private.
func (a *API) OrderDebug(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/order_debug", params, &raw)
	return raw, err
}

// CancelByNonce cancels an order by the nonce on its signed payload —
// useful when the caller has not received the order id back yet.
//
// Required params: `instrument_name`, `nonce`, `wallet`. Private.
func (a *API) CancelByNonce(ctx context.Context, instrument string, nonce uint64) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	params := map[string]any{
		"instrument_name": instrument,
		"nonce":           nonce,
		"wallet":          a.Signer.Owner().Hex(),
		"subaccount_id":   a.Subaccount,
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/cancel_by_nonce", params, &raw)
	return raw, err
}

// SetCancelOnDisconnect arms or disarms the kill-switch that cancels every
// open order on the wallet if the WebSocket session disconnects.
//
// Pass enabled=true to arm; false to disarm. Private.
func (a *API) SetCancelOnDisconnect(ctx context.Context, enabled bool) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	params := map[string]any{
		"wallet":  a.Signer.Owner().Hex(),
		"enabled": enabled,
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/set_cancel_on_disconnect", params, &raw)
	return raw, err
}

// ChangeSubaccountLabel sets the human-readable label on the configured
// subaccount.
func (a *API) ChangeSubaccountLabel(ctx context.Context, label string) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	params := map[string]any{
		"subaccount_id": a.Subaccount,
		"label":         label,
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/change_subaccount_label", params, &raw)
	return raw, err
}

// GetInstruments lists active instruments matching the filter. Public.
//
// Derive returns the result as a bare JSON array of instrument objects.
func (a *API) GetInstruments(ctx context.Context, currency string, kind InstrumentType) ([]Instrument, error) {
	params := map[string]any{}
	if currency != "" {
		params["currency"] = currency
	}
	if kind != "" {
		params["instrument_type"] = kind
	}
	params["expired"] = false
	var insts []Instrument
	if err := a.call(ctx, "public/get_instruments", params, &insts); err != nil {
		return nil, err
	}
	return insts, nil
}

// GetInstrument fetches one instrument by name. Public.
func (a *API) GetInstrument(ctx context.Context, name string) (Instrument, error) {
	var inst Instrument
	err := a.call(ctx, "public/get_instrument", map[string]any{"instrument_name": name}, &inst)
	return inst, err
}

// GetTicker fetches the public ticker for one instrument. Public.
func (a *API) GetTicker(ctx context.Context, name string) (Ticker, error) {
	var t Ticker
	err := a.call(ctx, "public/get_ticker", map[string]any{"instrument_name": name}, &t)
	return t, err
}

// GetPublicTradeHistory returns recent trades on the instrument. Public.
func (a *API) GetPublicTradeHistory(ctx context.Context, instrument string, page PageRequest) ([]Trade, Page, error) {
	params := map[string]any{"instrument_name": instrument}
	if page.Page > 0 {
		params["page"] = page.Page
	}
	if page.PageSize > 0 {
		params["page_size"] = page.PageSize
	}
	var resp struct {
		Trades     []Trade `json:"trades"`
		Pagination Page    `json:"pagination"`
	}
	if err := a.call(ctx, "public/get_trade_history", params, &resp); err != nil {
		return nil, Page{}, err
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

// MMPConfig is the input to SetMMPConfig — Market Maker Protection rules.
type MMPConfig struct {
	Currency        string `json:"currency"`
	MMPFrozenTimeMs int64  `json:"mmp_frozen_time"`
	MMPIntervalMs   int64  `json:"mmp_interval"`
	MMPAmountLimit  string `json:"mmp_amount_limit,omitempty"`
	MMPDeltaLimit   string `json:"mmp_delta_limit,omitempty"`
}

// Validate performs schema-level checks on the receiver. Returns nil on
// success or an error wrapping [types.ErrInvalidParams]. The two limit
// fields are decimal strings on the wire and remain unparsed here.
func (c MMPConfig) Validate() error {
	if c.Currency == "" {
		return invalidInput("currency", "required")
	}
	if c.MMPFrozenTimeMs < 0 {
		return invalidInput("mmp_frozen_time", "must be non-negative")
	}
	if c.MMPIntervalMs < 0 {
		return invalidInput("mmp_interval", "must be non-negative")
	}
	return nil
}

// SetMMPConfig configures market-maker protection for a currency. Private.
func (a *API) SetMMPConfig(ctx context.Context, cfg MMPConfig) error {
	if err := a.requireSubaccount(); err != nil {
		return err
	}
	params := map[string]any{
		"subaccount_id":   a.Subaccount,
		"currency":        cfg.Currency,
		"mmp_frozen_time": cfg.MMPFrozenTimeMs,
		"mmp_interval":    cfg.MMPIntervalMs,
	}
	if cfg.MMPAmountLimit != "" {
		params["mmp_amount_limit"] = cfg.MMPAmountLimit
	}
	if cfg.MMPDeltaLimit != "" {
		params["mmp_delta_limit"] = cfg.MMPDeltaLimit
	}
	return a.call(ctx, "private/set_mmp_config", params, nil)
}

// ResetMMP unfreezes the subaccount's MMP for a currency. Private.
func (a *API) ResetMMP(ctx context.Context, currency string) error {
	if err := a.requireSubaccount(); err != nil {
		return err
	}
	return a.call(ctx, "private/reset_mmp", map[string]any{
		"subaccount_id": a.Subaccount,
		"currency":      currency,
	}, nil)
}

// invalidInput wraps [ErrInvalidParams] for input DTOs declared in
// this package, so callers can match every Validate failure with one
// errors.Is regardless of where the DTO was declared.
func invalidInput(field, reason string) error {
	return fmt.Errorf("%w: %s: %s", ErrInvalidParams, field, reason)
}

// PlaceOrderInput is a thin convenience wrapper for the user-facing
// PlaceOrder. It contains only the strategically-relevant fields; the SDK
// fills in subaccount id, signature, signer, nonce and expiry from the
// configured signer and ambient state.
type PlaceOrderInput struct {
	InstrumentName string
	Asset          common.Address
	SubID          uint64
	Direction      Direction
	OrderType      OrderType
	TimeInForce    TimeInForce
	Amount         Decimal
	LimitPrice     Decimal
	MaxFee         Decimal
	Label          string
	MMP            bool
	ReduceOnly     bool
}

// Validate performs schema-level checks on the receiver: required fields
// populated, enum values in range, numeric fields in bounds. It does not
// validate against an instrument's tick / amount step (those live on
// [Instrument] and require a network round-trip).
//
// Returns nil on success or an error wrapping [ErrInvalidParams].
func (in PlaceOrderInput) Validate() error {
	if in.InstrumentName == "" {
		return invalidInput("instrument_name", "required")
	}
	if in.Asset == (common.Address{}) {
		return invalidInput("asset", "required")
	}
	if err := in.Direction.Validate(); err != nil {
		return invalidInput("direction", err.Error())
	}
	if err := in.OrderType.Validate(); err != nil {
		return invalidInput("order_type", err.Error())
	}
	if in.TimeInForce != "" {
		if err := in.TimeInForce.Validate(); err != nil {
			return invalidInput("time_in_force", err.Error())
		}
	}
	if in.Amount.Sign() <= 0 {
		return invalidInput("amount", "must be positive")
	}
	if in.LimitPrice.Sign() <= 0 {
		return invalidInput("limit_price", "must be positive")
	}
	if in.MaxFee.Sign() < 0 {
		return invalidInput("max_fee", "must be non-negative")
	}
	return nil
}

// PlaceOrder builds, signs and submits an order. Private.
//
// The session key signs the action; the resulting signature, signer address,
// nonce and expiry are embedded in the JSON-RPC params so the matching engine
// can recompute the EIP-712 hash and verify.
func (a *API) PlaceOrder(ctx context.Context, in PlaceOrderInput) (Order, error) {
	if err := a.requireSigner(); err != nil {
		return Order{}, err
	}
	if err := a.requireSubaccount(); err != nil {
		return Order{}, err
	}

	nonce := a.Nonces.Next()
	expiry := time.Now().Unix() + a.SignatureExpiry

	module := common.HexToAddress(a.Domain.VerifyingContract)

	if a.tradeModuleOverride() != (common.Address{}) {
		module = a.tradeModuleOverride()
	}

	tmd := TradeModuleData{
		Asset:       in.Asset,
		SubID:       in.SubID,
		LimitPrice:  in.LimitPrice.Inner(),
		Amount:      in.Amount.Inner(),
		MaxFee:      in.MaxFee.Inner(),
		RecipientID: a.Subaccount,
		IsBid:       in.Direction == DirectionBuy,
	}
	dataHash, err := tmd.Hash()
	if err != nil {
		return Order{}, err
	}

	action := ActionData{
		SubaccountID: a.Subaccount,
		Nonce:        nonce,
		Module:       module,
		Data:         dataHash,
		Expiry:       expiry,
		Owner:        a.Signer.Owner(),
		Signer:       a.Signer.Address(),
	}
	sig, err := a.Signer.SignAction(ctx, a.Domain, action)
	if err != nil {
		return Order{}, err
	}

	params := OrderParams{
		InstrumentName:  in.InstrumentName,
		Direction:       in.Direction,
		OrderType:       in.OrderType,
		TimeInForce:     in.TimeInForce,
		Amount:          in.Amount,
		LimitPrice:      in.LimitPrice,
		MaxFee:          in.MaxFee,
		SubaccountID:    a.Subaccount,
		Nonce:           nonce,
		Signer:          Address(a.Signer.Address()),
		Signature:       sig.Hex(),
		SignatureExpiry: expiry,
		Label:           in.Label,
		MMP:             in.MMP,
		ReduceOnly:      in.ReduceOnly,
	}
	var resp struct {
		Order Order `json:"order"`
	}
	if err := a.call(ctx, "private/order", params, &resp); err != nil {
		return Order{}, err
	}
	return resp.Order, nil
}

// tradeModuleOverride returns the TradeModule address from the ambient
// Contracts struct if available. The API struct doesn't carry the
// full config to keep its size small; we expose it via a setter (see
// SetTradeModule below) that pkg/rest and pkg/ws set up at construction.
func (a *API) tradeModuleOverride() common.Address { return a.tradeModule }

// SetTradeModule is called by the client constructors to thread through the
// per-network TradeModule contract address.
func (a *API) SetTradeModule(addr common.Address) { a.tradeModule = addr }

// CancelOrder cancels one open order by id. Private.
func (a *API) CancelOrder(ctx context.Context, instrument, orderID string) error {
	if err := a.requireSubaccount(); err != nil {
		return err
	}
	params := map[string]any{
		"subaccount_id":   a.Subaccount,
		"instrument_name": instrument,
		"order_id":        orderID,
	}
	return a.call(ctx, "private/cancel", params, nil)
}

// CancelByLabel cancels all orders carrying the given label. Private.
func (a *API) CancelByLabel(ctx context.Context, label string) (cancelled int, err error) {
	if err := a.requireSubaccount(); err != nil {
		return 0, err
	}
	params := map[string]any{
		"subaccount_id": a.Subaccount,
		"label":         label,
	}
	var resp struct {
		CancelledOrders int `json:"cancelled_orders"`
	}
	if err := a.call(ctx, "private/cancel_by_label", params, &resp); err != nil {
		return 0, err
	}
	return resp.CancelledOrders, nil
}

// CancelByInstrument cancels all open orders on the instrument. Private.
func (a *API) CancelByInstrument(ctx context.Context, instrument string) (cancelled int, err error) {
	if err := a.requireSubaccount(); err != nil {
		return 0, err
	}
	params := map[string]any{
		"subaccount_id":   a.Subaccount,
		"instrument_name": instrument,
	}
	var resp struct {
		CancelledOrders int `json:"cancelled_orders"`
	}
	if err := a.call(ctx, "private/cancel_by_instrument", params, &resp); err != nil {
		return 0, err
	}
	return resp.CancelledOrders, nil
}

// CancelAll cancels every open order on the subaccount. Private.
func (a *API) CancelAll(ctx context.Context) (cancelled int, err error) {
	if err := a.requireSubaccount(); err != nil {
		return 0, err
	}
	params := map[string]any{"subaccount_id": a.Subaccount}
	var resp struct {
		CancelledOrders int `json:"cancelled_orders"`
	}
	if err := a.call(ctx, "private/cancel_all", params, &resp); err != nil {
		return 0, err
	}
	return resp.CancelledOrders, nil
}

// GetOrder fetches one order by id. Private.
func (a *API) GetOrder(ctx context.Context, orderID string) (Order, error) {
	if err := a.requireSubaccount(); err != nil {
		return Order{}, err
	}
	params := map[string]any{
		"subaccount_id": a.Subaccount,
		"order_id":      orderID,
	}
	var resp struct {
		Order Order `json:"order"`
	}
	err := a.call(ctx, "private/get_order", params, &resp)
	return resp.Order, err
}

// GetOpenOrders lists currently-open orders on the subaccount. Private.
func (a *API) GetOpenOrders(ctx context.Context) ([]Order, error) {
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	params := map[string]any{"subaccount_id": a.Subaccount}
	var resp struct {
		Orders []Order `json:"orders"`
	}
	err := a.call(ctx, "private/get_open_orders", params, &resp)
	return resp.Orders, err
}

// GetOrderHistory paginates past orders. Private.
func (a *API) GetOrderHistory(ctx context.Context, page PageRequest) ([]Order, Page, error) {
	if err := a.requireSubaccount(); err != nil {
		return nil, Page{}, err
	}
	params := map[string]any{"subaccount_id": a.Subaccount}
	if page.Page > 0 {
		params["page"] = page.Page
	}
	if page.PageSize > 0 {
		params["page_size"] = page.PageSize
	}
	var resp struct {
		Orders     []Order `json:"orders"`
		Pagination Page    `json:"pagination"`
	}
	if err := a.call(ctx, "private/get_orders", params, &resp); err != nil {
		return nil, Page{}, err
	}
	return resp.Orders, resp.Pagination, nil
}

// GetPositions lists open positions on the subaccount. Private.
func (a *API) GetPositions(ctx context.Context) ([]Position, error) {
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	params := map[string]any{"subaccount_id": a.Subaccount}
	var resp struct {
		Positions []Position `json:"positions"`
	}
	err := a.call(ctx, "private/get_positions", params, &resp)
	return resp.Positions, err
}

// SendRFQ broadcasts a request-for-quote to market makers. Private.
func (a *API) SendRFQ(ctx context.Context, legs []RFQLeg, maxFee Decimal) (RFQ, error) {
	if err := a.requireSubaccount(); err != nil {
		return RFQ{}, err
	}
	params := map[string]any{
		"subaccount_id": a.Subaccount,
		"legs":          legs,
		"max_total_fee": maxFee,
	}
	var rfq RFQ
	err := a.call(ctx, "private/send_rfq", params, &rfq)
	return rfq, err
}

// PollRFQs returns the status of recent RFQs initiated by this subaccount. Private.
func (a *API) PollRFQs(ctx context.Context) ([]RFQ, error) {
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	var resp struct {
		RFQs []RFQ `json:"rfqs"`
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

// GetRFQs returns the configured subaccount's outstanding (open / done)
// RFQs.
func (a *API) GetRFQs(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/get_rfqs", params, &raw)
	return raw, err
}

// GetQuotes returns quotes the configured subaccount has issued or
// received against open RFQs.
func (a *API) GetQuotes(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/get_quotes", params, &raw)
	return raw, err
}

// PollQuotes is the long-poll variant of GetQuotes — used by makers who
// want to be woken on new RFQs without holding a WebSocket open.
func (a *API) PollQuotes(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/poll_quotes", params, &raw)
	return raw, err
}

// SendQuote responds to an open RFQ with a maker quote. The signed
// payload covers the multi-leg quote price and a per-leg side direction.
//
// Required params include the RFQ id, the per-leg quote prices, and the
// signature/nonce/expiry triple. Private.
func (a *API) SendQuote(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/send_quote", params, &raw)
	return raw, err
}

// ExecuteQuote picks one quote response and trades against it. Used by
// the taker once `send_rfq` has surfaced acceptable quotes.
func (a *API) ExecuteQuote(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/execute_quote", params, &raw)
	return raw, err
}

// CancelQuote cancels one outstanding maker quote by id.
func (a *API) CancelQuote(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/cancel_quote", params, &raw)
	return raw, err
}

// CancelBatchQuotes cancels every quote whose id appears in `quote_ids`,
// or every open quote on the subaccount when the field is omitted.
func (a *API) CancelBatchQuotes(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/cancel_batch_quotes", params, &raw)
	return raw, err
}

// CancelBatchRFQs cancels every RFQ whose id appears in `rfq_ids`, or
// every open RFQ on the subaccount when the field is omitted.
func (a *API) CancelBatchRFQs(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/cancel_batch_rfqs", params, &raw)
	return raw, err
}

// RFQGetBestQuote returns the best quote currently outstanding on one
// RFQ — the helper a taker uses to pick a counterparty before calling
// ExecuteQuote.
func (a *API) RFQGetBestQuote(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/rfq_get_best_quote", params, &raw)
	return raw, err
}

// OrderQuote routes an order through the RFQ matching path instead of
// the central order book. Useful for instruments with thin books where
// makers respond on demand.
func (a *API) OrderQuote(ctx context.Context, params map[string]any) (json.RawMessage, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["subaccount_id"]; !ok {
		params["subaccount_id"] = a.Subaccount
	}
	var raw json.RawMessage
	err := a.call(ctx, "private/order_quote", params, &raw)
	return raw, err
}

// GetSubaccount fetches the configured subaccount snapshot. Private.
func (a *API) GetSubaccount(ctx context.Context) (SubAccount, error) {
	if err := a.requireSubaccount(); err != nil {
		return SubAccount{}, err
	}
	var sa SubAccount
	err := a.call(ctx, "private/get_subaccount", map[string]any{
		"subaccount_id": a.Subaccount,
	}, &sa)
	return sa, err
}

// GetSubaccounts lists every subaccount owned by the wallet. Private.
func (a *API) GetSubaccounts(ctx context.Context) ([]int64, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	var resp struct {
		SubaccountIDs []int64 `json:"subaccount_ids"`
	}
	err := a.call(ctx, "private/get_subaccounts", map[string]any{
		"wallet": a.Signer.Owner().Hex(),
	}, &resp)
	return resp.SubaccountIDs, err
}

// GetTradeHistory paginates the user's fills. Private.
func (a *API) GetTradeHistory(ctx context.Context, page PageRequest) ([]Trade, Page, error) {
	if err := a.requireSubaccount(); err != nil {
		return nil, Page{}, err
	}
	params := map[string]any{"subaccount_id": a.Subaccount}
	if page.Page > 0 {
		params["page"] = page.Page
	}
	if page.PageSize > 0 {
		params["page_size"] = page.PageSize
	}
	var resp struct {
		Trades     []Trade `json:"trades"`
		Pagination Page    `json:"pagination"`
	}
	if err := a.call(ctx, "private/get_trade_history", params, &resp); err != nil {
		return nil, Page{}, err
	}
	return resp.Trades, resp.Pagination, nil
}

// GetDepositHistory paginates deposit transactions. Private.
func (a *API) GetDepositHistory(ctx context.Context, page PageRequest) ([]DepositTx, Page, error) {
	if err := a.requireSubaccount(); err != nil {
		return nil, Page{}, err
	}
	params := map[string]any{"subaccount_id": a.Subaccount}
	if page.Page > 0 {
		params["page"] = page.Page
	}
	if page.PageSize > 0 {
		params["page_size"] = page.PageSize
	}
	var resp struct {
		Events     []DepositTx `json:"events"`
		Pagination Page        `json:"pagination"`
	}
	if err := a.call(ctx, "private/get_deposit_history", params, &resp); err != nil {
		return nil, Page{}, err
	}
	return resp.Events, resp.Pagination, nil
}

// GetWithdrawalHistory paginates withdrawal transactions. Private.
func (a *API) GetWithdrawalHistory(ctx context.Context, page PageRequest) ([]WithdrawTx, Page, error) {
	if err := a.requireSubaccount(); err != nil {
		return nil, Page{}, err
	}
	params := map[string]any{"subaccount_id": a.Subaccount}
	if page.Page > 0 {
		params["page"] = page.Page
	}
	if page.PageSize > 0 {
		params["page_size"] = page.PageSize
	}
	var resp struct {
		Events     []WithdrawTx `json:"events"`
		Pagination Page         `json:"pagination"`
	}
	if err := a.call(ctx, "private/get_withdrawal_history", params, &resp); err != nil {
		return nil, Page{}, err
	}
	return resp.Events, resp.Pagination, nil
}

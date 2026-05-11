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
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/amiwrpremium/go-derive/pkg/auth"
	"github.com/amiwrpremium/go-derive/pkg/enums"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

// PlaceOrder builds, signs and submits an order. Private.
//
// The session key signs the action; the resulting signature, signer address,
// nonce and expiry are embedded in the JSON-RPC params so the matching engine
// can recompute the EIP-712 hash and verify.
func (a *API) PlaceOrder(ctx context.Context, in types.PlaceOrderInput) (types.Order, error) {
	params, err := a.signedOrderParams(ctx, in)
	if err != nil {
		return types.Order{}, err
	}
	var resp struct {
		Order types.Order `json:"order"`
	}
	if err := a.call(ctx, "private/order", params, &resp); err != nil {
		return types.Order{}, err
	}
	return resp.Order, nil
}

// signedOrderParams is the shared signing block for every order-
// submission endpoint (private/order, private/algo_order,
// private/trigger_order). It produces the JSON-RPC params map with
// the standard order fields plus the EIP-712 action signature.
// Method-specific extras (algo_*, trigger_*) are added by the caller.
func (a *API) signedOrderParams(ctx context.Context, in types.PlaceOrderInput) (map[string]any, error) {
	if err := a.requireSigner(); err != nil {
		return nil, err
	}
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}

	nonce := a.Nonces.Next()
	expiry := time.Now().Unix() + a.SignatureExpiry

	module := common.HexToAddress(a.Domain.VerifyingContract) // override below
	// The TradeModule address differs from the matching engine domain; the
	// caller-side wiring fills it in via the netconf.Contracts struct. For
	// safety we read it from a hidden field on the action input.
	if a.tradeModuleOverride() != (common.Address{}) {
		module = a.tradeModuleOverride()
	}

	tmd := auth.TradeModuleData{
		Asset:       common.Address(in.Asset),
		SubID:       in.SubID,
		LimitPrice:  in.LimitPrice.Inner(),
		Amount:      in.Amount.Inner(),
		MaxFee:      in.MaxFee.Inner(),
		RecipientID: a.Subaccount,
		IsBid:       in.Direction == enums.DirectionBuy,
	}
	dataHash, err := tmd.Hash()
	if err != nil {
		return nil, err
	}

	action := auth.ActionData{
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
		return nil, err
	}

	params := map[string]any{
		"instrument_name":      in.InstrumentName,
		"direction":            in.Direction,
		"order_type":           in.OrderType,
		"amount":               in.Amount,
		"limit_price":          in.LimitPrice,
		"max_fee":              in.MaxFee,
		"subaccount_id":        a.Subaccount,
		"nonce":                nonce,
		"signer":               a.Signer.Address().Hex(),
		"signature":            sig.Hex(),
		"signature_expiry_sec": expiry,
	}
	if in.TimeInForce != "" {
		params["time_in_force"] = in.TimeInForce
	}
	if in.Label != "" {
		params["label"] = in.Label
	}
	if in.MMP {
		params["mmp"] = true
	}
	if in.ReduceOnly {
		params["reduce_only"] = true
	}
	return params, nil
}

// PlaceAlgoOrder builds, signs and submits a TWAP-style algorithmic
// order. Private. Wraps `private/algo_order`.
//
// Signing flow is identical to [API.PlaceOrder] — the algo-specific
// fields (algo_type, algo_duration_sec, algo_num_slices) are added
// alongside the standard order params; they are not part of the
// EIP-712 signature payload.
//
// Returns the engine's order record in `algo_active` state. Cancel
// with [API.CancelAlgoOrder] or [API.CancelAllAlgoOrders].
func (a *API) PlaceAlgoOrder(ctx context.Context, in types.AlgoOrderInput) (types.Order, error) {
	params, err := a.signedOrderParams(ctx, in.PlaceOrderInput)
	if err != nil {
		return types.Order{}, err
	}
	params["algo_type"] = in.AlgoType
	params["algo_duration_sec"] = in.AlgoDurationSec
	params["algo_num_slices"] = in.AlgoNumSlices

	var resp struct {
		Order types.Order `json:"order"`
	}
	if err := a.call(ctx, "private/algo_order", params, &resp); err != nil {
		return types.Order{}, err
	}
	return resp.Order, nil
}

// PlaceTriggerOrder builds, signs and submits a stop-loss or
// take-profit trigger order. Private. Wraps `private/trigger_order`.
//
// Signing flow is identical to [API.PlaceOrder]; the trigger-
// specific fields (trigger_type, trigger_price_type, trigger_price)
// are added alongside the standard order params and are not part of
// the EIP-712 signature payload.
//
// The order is saved server-side in `untriggered` state until the
// matching engine sees the watched price cross the trigger level.
// Cancel ahead of time with [API.CancelTriggerOrder] or
// [API.CancelAllTriggerOrders].
func (a *API) PlaceTriggerOrder(ctx context.Context, in types.TriggerOrderInput) (types.Order, error) {
	params, err := a.signedOrderParams(ctx, in.PlaceOrderInput)
	if err != nil {
		return types.Order{}, err
	}
	params["trigger_type"] = in.TriggerType
	params["trigger_price_type"] = in.TriggerPriceType
	params["trigger_price"] = in.TriggerPrice

	var resp struct {
		Order types.Order `json:"order"`
	}
	if err := a.call(ctx, "private/trigger_order", params, &resp); err != nil {
		return types.Order{}, err
	}
	return resp.Order, nil
}

// tradeModuleOverride returns the TradeModule address from the ambient
// netconf.Contracts struct if available. The API struct doesn't carry the
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
func (a *API) GetOrder(ctx context.Context, orderID string) (types.Order, error) {
	if err := a.requireSubaccount(); err != nil {
		return types.Order{}, err
	}
	params := map[string]any{
		"subaccount_id": a.Subaccount,
		"order_id":      orderID,
	}
	var resp struct {
		Order types.Order `json:"order"`
	}
	err := a.call(ctx, "private/get_order", params, &resp)
	return resp.Order, err
}

// GetOpenOrders lists currently-open orders on the subaccount. Private.
func (a *API) GetOpenOrders(ctx context.Context) ([]types.Order, error) {
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	params := map[string]any{"subaccount_id": a.Subaccount}
	var resp struct {
		Orders []types.Order `json:"orders"`
	}
	err := a.call(ctx, "private/get_open_orders", params, &resp)
	return resp.Orders, err
}

// GetOrders paginates orders on the configured subaccount, narrowed
// by the supplied filter. Private.
//
// Wraps `/private/get_orders`. Pass `filter == nil` to omit all
// filters and page through every order on the subaccount.
//
// To page through orders by time window, use [API.GetOrderHistory]
// instead — `/private/get_orders` only filters by status /
// instrument / label, not by time.
func (a *API) GetOrders(ctx context.Context, page types.PageRequest, filter *types.GetOrdersFilter) ([]types.Order, types.Page, error) {
	if err := a.requireSubaccount(); err != nil {
		return nil, types.Page{}, err
	}
	params := map[string]any{"subaccount_id": a.Subaccount}
	if page.Page > 0 {
		params["page"] = page.Page
	}
	if page.PageSize > 0 {
		params["page_size"] = page.PageSize
	}
	if filter != nil {
		if filter.InstrumentName != "" {
			params["instrument_name"] = filter.InstrumentName
		}
		if filter.Label != "" {
			params["label"] = filter.Label
		}
		if filter.Status != "" {
			params["status"] = filter.Status
		}
	}
	var resp struct {
		Orders     []types.Order `json:"orders"`
		Pagination types.Page    `json:"pagination"`
	}
	if err := a.call(ctx, "private/get_orders", params, &resp); err != nil {
		return nil, types.Page{}, err
	}
	return resp.Orders, resp.Pagination, nil
}

// GetOrderHistory paginates past orders for the configured
// subaccount or wallet, filtered by a time window. Private.
//
// Wraps `/private/get_order_history`. The configured subaccount is
// threaded through when both Wallet is empty and the subaccount is
// non-zero — supply a Wallet to query across every subaccount the
// wallet owns.
func (a *API) GetOrderHistory(ctx context.Context, page types.PageRequest, q types.OrderHistoryQuery) ([]types.Order, types.Page, error) {
	if err := a.requireSigner(); err != nil {
		return nil, types.Page{}, err
	}
	params := map[string]any{}
	if q.Wallet != "" {
		params["wallet"] = q.Wallet
	} else {
		if err := a.requireSubaccount(); err != nil {
			return nil, types.Page{}, err
		}
		params["subaccount_id"] = a.Subaccount
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
		Orders     []types.Order `json:"orders"`
		Pagination types.Page    `json:"pagination"`
	}
	if err := a.call(ctx, "private/get_order_history", params, &resp); err != nil {
		return nil, types.Page{}, err
	}
	return resp.Orders, resp.Pagination, nil
}

// Replace cancels one outstanding order and submits a replacement
// in a single round trip — the standard maker pattern for
// re-pricing without a race against the matching engine. Private.
//
// `params` should include `order_id_to_cancel` (or
// `nonce_to_cancel`) plus the same fields PlaceOrder takes for
// the replacement. The full param shape is documented at
// docs.derive.xyz.
//
// The response carries the cancelled order, the (optional)
// replacement order, the engine's error if the replacement was
// rejected, and the trades the new order matched.
func (a *API) Replace(ctx context.Context, params map[string]any) (types.ReplaceResult, error) {
	if err := a.requireSigner(); err != nil {
		return types.ReplaceResult{}, err
	}
	if err := a.requireSubaccount(); err != nil {
		return types.ReplaceResult{}, err
	}
	var resp types.ReplaceResult
	if err := a.call(ctx, "private/replace", params, &resp); err != nil {
		return types.ReplaceResult{}, err
	}
	return resp, nil
}

// OrderDebug previews an order without submitting it. Private.
//
// `params` mirror PlaceOrder's. The response carries the
// engine's internal hashing artefacts — useful for validating
// signatures in CI.
func (a *API) OrderDebug(ctx context.Context, params map[string]any) (types.OrderDebugResult, error) {
	if err := a.requireSigner(); err != nil {
		return types.OrderDebugResult{}, err
	}
	if err := a.requireSubaccount(); err != nil {
		return types.OrderDebugResult{}, err
	}
	var resp types.OrderDebugResult
	if err := a.call(ctx, "private/order_debug", params, &resp); err != nil {
		return types.OrderDebugResult{}, err
	}
	return resp, nil
}

// CancelByNonce cancels an order by the nonce on its signed
// payload — useful when the caller has not yet received the
// order id back. Private.
//
// Returns the number of orders that matched the (instrument,
// nonce) tuple and were cancelled.
func (a *API) CancelByNonce(ctx context.Context, instrument string, nonce uint64) (types.CancelByNonceResult, error) {
	if err := a.requireSigner(); err != nil {
		return types.CancelByNonceResult{}, err
	}
	if err := a.requireSubaccount(); err != nil {
		return types.CancelByNonceResult{}, err
	}
	params := map[string]any{
		"instrument_name": instrument,
		"nonce":           nonce,
		"wallet":          a.Signer.Owner().Hex(),
		"subaccount_id":   a.Subaccount,
	}
	var resp types.CancelByNonceResult
	if err := a.call(ctx, "private/cancel_by_nonce", params, &resp); err != nil {
		return types.CancelByNonceResult{}, err
	}
	return resp, nil
}

// CancelAlgoOrder cancels one in-flight algo order by id. Private.
//
// Returns the cancelled order (in `algo_active` -> `cancelled`
// state). Counterpart to [API.CancelTriggerOrder] for algo orders
// (e.g. TWAP) that have started slicing into the market.
func (a *API) CancelAlgoOrder(ctx context.Context, orderID string) (types.Order, error) {
	if err := a.requireSubaccount(); err != nil {
		return types.Order{}, err
	}
	params := map[string]any{
		"subaccount_id": a.Subaccount,
		"order_id":      orderID,
	}
	var resp types.Order
	if err := a.call(ctx, "private/cancel_algo_order", params, &resp); err != nil {
		return types.Order{}, err
	}
	return resp, nil
}

// CancelAllAlgoOrders cancels every in-flight algo order on the
// configured subaccount. Private.
//
// Returns nil on success; the wire response is a fixed "ok" string
// surfaced as a nil error.
func (a *API) CancelAllAlgoOrders(ctx context.Context) error {
	if err := a.requireSubaccount(); err != nil {
		return err
	}
	return a.call(ctx, "private/cancel_all_algo_orders", map[string]any{
		"subaccount_id": a.Subaccount,
	}, nil)
}

// GetAlgoOrders lists every active algo order on the configured
// subaccount. Private.
//
// Counterpart to [API.GetOpenOrders] for algo orders. Returns a bare
// list — no pagination wrapper.
func (a *API) GetAlgoOrders(ctx context.Context) ([]types.Order, error) {
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	var resp []types.Order
	if err := a.call(ctx, "private/get_algo_orders", map[string]any{
		"subaccount_id": a.Subaccount,
	}, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetTriggerOrders lists every untriggered trigger order on the
// configured subaccount. Private.
//
// The wire response wraps the list in a `{subaccount_id, orders[]}`
// envelope, mirroring [API.GetOpenOrders].
func (a *API) GetTriggerOrders(ctx context.Context) ([]types.Order, error) {
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	var resp struct {
		Orders []types.Order `json:"orders"`
	}
	if err := a.call(ctx, "private/get_trigger_orders", map[string]any{
		"subaccount_id": a.Subaccount,
	}, &resp); err != nil {
		return nil, err
	}
	return resp.Orders, nil
}

// CancelTriggerOrder cancels one untriggered trigger order by id.
// Private.
//
// Returns the cancelled order (in `untriggered` -> `cancelled`
// state). Counterpart to [API.CancelOrder] for trigger orders that
// have not yet fired.
func (a *API) CancelTriggerOrder(ctx context.Context, orderID string) (types.Order, error) {
	if err := a.requireSubaccount(); err != nil {
		return types.Order{}, err
	}
	params := map[string]any{
		"subaccount_id": a.Subaccount,
		"order_id":      orderID,
	}
	var resp types.Order
	if err := a.call(ctx, "private/cancel_trigger_order", params, &resp); err != nil {
		return types.Order{}, err
	}
	return resp, nil
}

// CancelAllTriggerOrders cancels every untriggered trigger order on
// the configured subaccount. Private.
//
// Returns nil on success; the wire response is a fixed "ok" string
// surfaced as a nil error.
func (a *API) CancelAllTriggerOrders(ctx context.Context) error {
	if err := a.requireSubaccount(); err != nil {
		return err
	}
	return a.call(ctx, "private/cancel_all_trigger_orders", map[string]any{
		"subaccount_id": a.Subaccount,
	}, nil)
}

// SetCancelOnDisconnect arms or disarms the kill-switch that
// cancels every open order on the wallet if the WebSocket
// session disconnects. Private.
//
// Pass `enabled=true` to arm; `false` to disarm. The endpoint
// returns plain `"ok"` on success — surfaced as a nil error.
func (a *API) SetCancelOnDisconnect(ctx context.Context, enabled bool) error {
	if err := a.requireSigner(); err != nil {
		return err
	}
	params := map[string]any{
		"wallet":  a.Signer.Owner().Hex(),
		"enabled": enabled,
	}
	return a.call(ctx, "private/set_cancel_on_disconnect", params, nil)
}

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
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/pkg/auth"
)

// invalidInput wraps [derive.ErrInvalidParams] for input DTOs declared in
// this package, so callers can match every Validate failure with one
// errors.Is regardless of where the DTO was declared.
func invalidInput(field, reason string) error {
	return fmt.Errorf("%w: %s: %s", derive.ErrInvalidParams, field, reason)
}

// PlaceOrderInput is a thin convenience wrapper for the user-facing
// PlaceOrder. It contains only the strategically-relevant fields; the SDK
// fills in subaccount id, signature, signer, nonce and expiry from the
// configured signer and ambient state.
type PlaceOrderInput struct {
	InstrumentName string
	Asset          common.Address
	SubID          uint64
	Direction      derive.Direction
	OrderType      derive.OrderType
	TimeInForce    derive.TimeInForce
	Amount         derive.Decimal
	LimitPrice     derive.Decimal
	MaxFee         derive.Decimal
	Label          string
	MMP            bool
	ReduceOnly     bool
}

// Validate performs schema-level checks on the receiver: required fields
// populated, enum values in range, numeric fields in bounds. It does not
// validate against an instrument's tick / amount step (those live on
// [derive.Instrument] and require a network round-trip).
//
// Returns nil on success or an error wrapping [derive.ErrInvalidParams].
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
func (a *API) PlaceOrder(ctx context.Context, in PlaceOrderInput) (derive.Order, error) {
	if err := a.requireSigner(); err != nil {
		return derive.Order{}, err
	}
	if err := a.requireSubaccount(); err != nil {
		return derive.Order{}, err
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
		Asset:       in.Asset,
		SubID:       in.SubID,
		LimitPrice:  in.LimitPrice.Inner(),
		Amount:      in.Amount.Inner(),
		MaxFee:      in.MaxFee.Inner(),
		RecipientID: a.Subaccount,
		IsBid:       in.Direction == derive.DirectionBuy,
	}
	dataHash, err := tmd.Hash()
	if err != nil {
		return derive.Order{}, err
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
		return derive.Order{}, err
	}

	params := derive.OrderParams{
		InstrumentName:  in.InstrumentName,
		Direction:       in.Direction,
		OrderType:       in.OrderType,
		TimeInForce:     in.TimeInForce,
		Amount:          in.Amount,
		LimitPrice:      in.LimitPrice,
		MaxFee:          in.MaxFee,
		SubaccountID:    a.Subaccount,
		Nonce:           nonce,
		Signer:          derive.Address(a.Signer.Address()),
		Signature:       sig.Hex(),
		SignatureExpiry: expiry,
		Label:           in.Label,
		MMP:             in.MMP,
		ReduceOnly:      in.ReduceOnly,
	}
	var resp struct {
		Order derive.Order `json:"order"`
	}
	if err := a.call(ctx, "private/order", params, &resp); err != nil {
		return derive.Order{}, err
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
func (a *API) GetOrder(ctx context.Context, orderID string) (derive.Order, error) {
	if err := a.requireSubaccount(); err != nil {
		return derive.Order{}, err
	}
	params := map[string]any{
		"subaccount_id": a.Subaccount,
		"order_id":      orderID,
	}
	var resp struct {
		Order derive.Order `json:"order"`
	}
	err := a.call(ctx, "private/get_order", params, &resp)
	return resp.Order, err
}

// GetOpenOrders lists currently-open orders on the subaccount. Private.
func (a *API) GetOpenOrders(ctx context.Context) ([]derive.Order, error) {
	if err := a.requireSubaccount(); err != nil {
		return nil, err
	}
	params := map[string]any{"subaccount_id": a.Subaccount}
	var resp struct {
		Orders []derive.Order `json:"orders"`
	}
	err := a.call(ctx, "private/get_open_orders", params, &resp)
	return resp.Orders, err
}

// GetOrderHistory paginates past orders. Private.
func (a *API) GetOrderHistory(ctx context.Context, page derive.PageRequest) ([]derive.Order, derive.Page, error) {
	if err := a.requireSubaccount(); err != nil {
		return nil, derive.Page{}, err
	}
	params := map[string]any{"subaccount_id": a.Subaccount}
	if page.Page > 0 {
		params["page"] = page.Page
	}
	if page.PageSize > 0 {
		params["page_size"] = page.PageSize
	}
	var resp struct {
		Orders     []derive.Order `json:"orders"`
		Pagination derive.Page    `json:"pagination"`
	}
	if err := a.call(ctx, "private/get_orders", params, &resp); err != nil {
		return nil, derive.Page{}, err
	}
	return resp.Orders, resp.Pagination, nil
}

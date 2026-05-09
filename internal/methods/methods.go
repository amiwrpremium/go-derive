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

	"github.com/ethereum/go-ethereum/common"

	"github.com/amiwrpremium/go-derive"
	"github.com/amiwrpremium/go-derive/internal/transport"
)

// API is a transport-agnostic facade that holds the ambient configuration
// (signer, subaccount id, EIP-712 domain) used by the methods defined in this
// package. Construct it once per client.
type API struct {
	T      transport.Transport
	Signer derive.Signer
	Domain derive.Domain
	// Subaccount is the default subaccount id used by private endpoints. It
	// can be 0 for public-only clients.
	Subaccount int64
	// Nonces is the source of action nonces; one per client is fine.
	Nonces *derive.NonceGen
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
		return derive.ErrUnauthorized
	}
	return nil
}

// requireSubaccount returns ErrSubaccountRequired when the action needs one.
func (a *API) requireSubaccount() error {
	if a.Subaccount == 0 {
		return derive.ErrSubaccountRequired
	}
	return nil
}

// call is a shortcut for the common case of one method, one params, one out.
// It also re-wraps a transport-level [transport.JSONRPCError] into a
// public [derive.APIError] so callers receive the rich-typed sentinel-aware
// error.
func (a *API) call(ctx context.Context, method string, params, out any) error {
	return wrapTransportError(a.T.Call(ctx, method, params, out))
}

// wrapTransportError converts a transport-layer `*transport.JSONRPCError`
// into a public `*derive.APIError`. Errors of any other type pass through
// unchanged. This sits at the methods/transport boundary and is the reason
// transport can stay free of an import on `pkg/errors` — closing the
// otherwise-inevitable `pkg/errors → transport → pkg/errors` cycle once
// pkg/errors lifts to root.
func wrapTransportError(err error) error {
	if rpcErr, ok := err.(*transport.JSONRPCError); ok {
		return &derive.APIError{
			Code:    rpcErr.Code,
			Message: rpcErr.Message,
			Data:    rpcErr.Data,
		}
	}
	return err
}

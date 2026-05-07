// Package methods.
package methods

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/amiwrpremium/go-derive/internal/netconf"
	"github.com/amiwrpremium/go-derive/internal/transport"
	"github.com/amiwrpremium/go-derive/pkg/auth"
	derrors "github.com/amiwrpremium/go-derive/pkg/errors"
)

// API is a transport-agnostic facade that holds the ambient configuration
// (signer, subaccount id, EIP-712 domain) used by the methods defined in this
// package. Construct it once per client.
type API struct {
	T      transport.Transport
	Signer auth.Signer
	Domain netconf.Domain
	// Subaccount is the default subaccount id used by private endpoints. It
	// can be 0 for public-only clients.
	Subaccount int64
	// Nonces is the source of action nonces; one per client is fine.
	Nonces *auth.NonceGen
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
		return derrors.ErrUnauthorized
	}
	return nil
}

// requireSubaccount returns ErrSubaccountRequired when the action needs one.
func (a *API) requireSubaccount() error {
	if a.Subaccount == 0 {
		return derrors.ErrSubaccountRequired
	}
	return nil
}

// call is a shortcut for the common case of one method, one params, one out.
func (a *API) call(ctx context.Context, method string, params, out any) error {
	return a.T.Call(ctx, method, params, out)
}

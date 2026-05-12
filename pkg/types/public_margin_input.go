// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the input DTO for `public/get_margin`, the
// unauthenticated risk-engine margin simulator.
package types

import "github.com/amiwrpremium/go-derive/pkg/enums"

// SimulatedCollateral is one entry in the collateral list passed to
// `public/get_margin`. Used both as the baseline portfolio
// ([PublicMarginInput.SimulatedCollaterals]) and as a simulated
// change ([PublicMarginInput.SimulatedCollateralChanges]).
type SimulatedCollateral struct {
	// AssetName is the ERC-20 asset symbol, e.g. "ETH", "USDC",
	// "WSTETH".
	AssetName string `json:"asset_name"`
	// Amount is the collateral quantity in base units.
	Amount Decimal `json:"amount"`
}

// SimulatedPosition is one entry in the position list passed to
// `public/get_margin`. Used both as the baseline portfolio and as a
// simulated change.
type SimulatedPosition struct {
	// InstrumentName identifies the position's market.
	InstrumentName string `json:"instrument_name"`
	// Amount is the position quantity in base-currency units; signed
	// for shorts.
	Amount Decimal `json:"amount"`
	// EntryPrice is the entry price to use in the simulation. Only
	// honoured for perps; defaults to the mark price when zero.
	EntryPrice Decimal `json:"entry_price,omitempty"`
}

// PublicMarginInput is the typed body for `public/get_margin`. The
// risk engine evaluates the supplied portfolio (collaterals +
// positions) against the requested margin model and returns the
// pre/post initial- and maintenance-margin requirements.
//
// Required: [MarginType], [SimulatedCollaterals], [SimulatedPositions].
// Required only for portfolio margin: [Market]. The simulated-change
// arrays are optional and let the caller layer "what if I added /
// closed this" on top of the baseline.
type PublicMarginInput struct {
	// MarginType is the margin model: SM, PM, or PM2.
	MarginType enums.MarginType
	// Market scopes the calculation to one market — required for
	// portfolio margin and ignored for SM.
	Market string
	// SimulatedCollaterals is the baseline collateral portfolio.
	SimulatedCollaterals []SimulatedCollateral
	// SimulatedPositions is the baseline position portfolio.
	SimulatedPositions []SimulatedPosition
	// SimulatedCollateralChanges layers deposits / withdrawals /
	// spot trades on top of the baseline. Optional.
	SimulatedCollateralChanges []SimulatedCollateral
	// SimulatedPositionChanges layers perp / option trades on top of
	// the baseline. Optional.
	SimulatedPositionChanges []SimulatedPosition
}

// Validate performs schema-level checks on the receiver. Returns
// nil on success or a wrapped [ErrInvalidParams].
func (in PublicMarginInput) Validate() error {
	if err := in.MarginType.Validate(); err != nil {
		return invalidParam("margin_type", err.Error())
	}
	if len(in.SimulatedCollaterals) == 0 {
		return invalidParam("simulated_collaterals", "must have at least one entry")
	}
	if in.SimulatedPositions == nil {
		return invalidParam("simulated_positions", "required (may be empty)")
	}
	return nil
}

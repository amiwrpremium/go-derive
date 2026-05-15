// Package types declares the domain types used in REST and WebSocket
// requests and responses.
//
// This file holds the input DTOs for the single-target cancel
// endpoints (private/cancel, private/cancel_by_*, private/cancel_rfq,
// private/cancel_quote, etc.). Multi-cancel inputs live elsewhere —
// see [CancelBatchInput] for the batch RFQ/quote cancel surface.
package types

// CancelOrderInput parameterises a single-order cancel
// (private/cancel).
type CancelOrderInput struct {
	// InstrumentName is the market on which the order rests.
	InstrumentName string
	// OrderID is the engine-assigned id of the order to cancel.
	OrderID string
}

// Validate enforces both required fields are populated. Returns nil on
// success or a wrapped [ErrInvalidParams].
func (in CancelOrderInput) Validate() error {
	if in.InstrumentName == "" {
		return invalidParam("instrument_name", "required")
	}
	if in.OrderID == "" {
		return invalidParam("order_id", "required")
	}
	return nil
}

// CancelByInstrumentInput parameterises a cancel-by-instrument call
// (private/cancel_by_instrument).
type CancelByInstrumentInput struct {
	// InstrumentName scopes the cancel to one market.
	InstrumentName string
}

// Validate enforces InstrumentName is populated.
func (in CancelByInstrumentInput) Validate() error {
	if in.InstrumentName == "" {
		return invalidParam("instrument_name", "required")
	}
	return nil
}

// CancelByLabelInput parameterises a cancel-by-label call
// (private/cancel_by_label).
type CancelByLabelInput struct {
	// Label is the user-defined label assigned at order-submission
	// time.
	Label string
}

// Validate enforces Label is populated.
func (in CancelByLabelInput) Validate() error {
	if in.Label == "" {
		return invalidParam("label", "required")
	}
	return nil
}

// CancelByNonceInput parameterises a cancel-by-nonce call
// (private/cancel_by_nonce).
type CancelByNonceInput struct {
	// InstrumentName scopes the cancel to one market.
	InstrumentName string
	// Nonce is the signed nonce of the order to cancel.
	Nonce uint64
}

// Validate enforces both fields are populated.
func (in CancelByNonceInput) Validate() error {
	if in.InstrumentName == "" {
		return invalidParam("instrument_name", "required")
	}
	if in.Nonce == 0 {
		return invalidParam("nonce", "required")
	}
	return nil
}

// CancelAlgoOrderInput parameterises a single-algo-order cancel
// (private/cancel_algo_order).
type CancelAlgoOrderInput struct {
	// OrderID is the engine-assigned id of the algo order.
	OrderID string
}

// Validate enforces OrderID is populated.
func (in CancelAlgoOrderInput) Validate() error {
	if in.OrderID == "" {
		return invalidParam("order_id", "required")
	}
	return nil
}

// CancelTriggerOrderInput parameterises a single-trigger-order cancel
// (private/cancel_trigger_order).
type CancelTriggerOrderInput struct {
	// OrderID is the engine-assigned id of the trigger order.
	OrderID string
}

// Validate enforces OrderID is populated.
func (in CancelTriggerOrderInput) Validate() error {
	if in.OrderID == "" {
		return invalidParam("order_id", "required")
	}
	return nil
}

// CancelRFQInput parameterises a single-RFQ cancel
// (private/cancel_rfq).
type CancelRFQInput struct {
	// RFQID identifies the RFQ to cancel.
	RFQID string
}

// Validate enforces RFQID is populated.
func (in CancelRFQInput) Validate() error {
	if in.RFQID == "" {
		return invalidParam("rfq_id", "required")
	}
	return nil
}

// CancelQuoteInput parameterises a single-quote cancel
// (private/cancel_quote).
type CancelQuoteInput struct {
	// QuoteID identifies the quote to cancel.
	QuoteID string
}

// Validate enforces QuoteID is populated.
func (in CancelQuoteInput) Validate() error {
	if in.QuoteID == "" {
		return invalidParam("quote_id", "required")
	}
	return nil
}

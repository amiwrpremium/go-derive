package enums

// TxStatus is the on-chain transaction lifecycle as reported by Derive.
//
// Trades, transfers, and other actions that emit an on-chain transaction
// carry a `tx_status` field that walks through these values: requested
// → pending → settled (or reverted/ignored on failure).
type TxStatus string

const (
	// TxStatusRequested — request accepted, transaction not yet submitted.
	TxStatusRequested TxStatus = "requested"
	// TxStatusPending — transaction submitted on-chain, awaiting confirmation.
	TxStatusPending TxStatus = "pending"
	// TxStatusSettled — transaction confirmed; final.
	TxStatusSettled TxStatus = "settled"
	// TxStatusReverted — on-chain execution reverted; final.
	TxStatusReverted TxStatus = "reverted"
	// TxStatusIgnored — request superseded or de-duplicated; final.
	TxStatusIgnored TxStatus = "ignored"
)

// Valid reports whether the receiver is one of the defined statuses.
func (s TxStatus) Valid() bool {
	switch s {
	case TxStatusRequested, TxStatusPending, TxStatusSettled, TxStatusReverted, TxStatusIgnored:
		return true
	default:
		return false
	}
}

// Terminal reports whether the status is final.
func (s TxStatus) Terminal() bool {
	switch s {
	case TxStatusSettled, TxStatusReverted, TxStatusIgnored:
		return true
	default:
		return false
	}
}

// Validate returns nil when the receiver is one of the defined wire
// values, or an error wrapping [ErrInvalidEnum] when it isn't.
func (s TxStatus) Validate() error {
	if s.Valid() {
		return nil
	}
	return invalid("TxStatus", string(s))
}

// Package codec is a small bag of low-level encoding helpers shared between
// pkg/auth (action signing) and pkg/types. None of it leaks to user code.
package codec

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// ParseAddress accepts hex with or without the 0x prefix and returns the
// checksummed common.Address. It returns an error on malformed input.
func ParseAddress(s string) (common.Address, error) {
	if !common.IsHexAddress(s) {
		return common.Address{}, fmt.Errorf("codec: invalid address %q", s)
	}
	return common.HexToAddress(s), nil
}

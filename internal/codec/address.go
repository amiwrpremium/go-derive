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

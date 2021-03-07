package asset

import (
	"fmt"
	"gobacktrader/btutil"
)

// ValidatePair takes a currency pair and returns a
// cleaned pair string along with an error if invalid.
func ValidatePair(pair string) (string, error) {
	pair = btutil.CleanString(pair)
	if len(pair) != 6 {
		return pair, fmt.Errorf("expecting a six character currency pair, got '%s'", pair)
	}

	return pair, nil
}

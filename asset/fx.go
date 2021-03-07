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

// SplitPair returns the two individual components of an FX pair.
func SplitPair(pair string) (string, string, error) {
	var ccy1, ccy2 string
	pair, err := ValidatePair(pair)
	if err != nil {
		return ccy1, ccy2, err
	}

	ccy1 = pair[:3]
	ccy2 = pair[3:]
	return ccy1, ccy2, nil
}

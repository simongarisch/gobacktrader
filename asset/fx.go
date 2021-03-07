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

// IsEquivalentPair returns true where we expect the rate to be static.
// For example, AUDUSD = 1.0, USDUSD = 1.0.
func IsEquivalentPair(pair string) (bool, error) {
	ccy1, ccy2, err := SplitPair(pair)
	if err != nil {
		return false, err
	}

	if ccy1 == ccy2 {
		return true, nil
	}
	return false, nil
}

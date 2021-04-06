package asset

import (
	"fmt"
	"gobacktrader/btutil"
)

// ValidateCurrency takes a currency code and returns
// a cleaned code along with an error if invalid.
func ValidateCurrency(ccy string) (string, error) {
	ccy = btutil.CleanString(ccy)
	if len(ccy) != 3 {
		return ccy, fmt.Errorf("'%s' is not a valid currency code", ccy)
	}
	return ccy, nil
}

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

// GetInversePair returns the inverse of some currency pair
func GetInversePair(pair string) (string, error) {
	ccy1, ccy2, err := SplitPair(pair)
	if err != nil {
		return "", err
	}

	return ccy2 + ccy1, nil
}

// FxRate associates an FX pair with a rate.
type FxRate struct {
	priceHistory
	pair string
}

// NewFxRate returns a new instance of FxRate
func NewFxRate(pair string, price Price) (*FxRate, error) {
	pair, err := ValidatePair(pair)
	fxrate := FxRate{pair: pair}
	fxrate.SetRate(price)
	return &fxrate, err
}

// GetPair returns the FX pair as a string
func (r FxRate) GetPair() string {
	return r.pair
}

// GetRate returns the FX rate for this pair.
func (r FxRate) GetRate() Price {
	return r.price
}

// SetRate sets the FX rate for this pair.
func (r *FxRate) SetRate(price Price) {
	r.price = price
}

// SetPrice does the same as SetRate and sets the rate for this pair.
func (r *FxRate) SetPrice(price Price) {
	r.SetRate(price)
}

// FxRates keeps track of FXRate instances.
// There should only ever be one instance of a pair
// (or its inverse) that is registered.
type FxRates struct {
	rates []*FxRate
}

// NewFxRates returns a new instance of FxRates.
func NewFxRates() *FxRates {
	return &FxRates{}
}

// Register adds an FXRate to the available FxRates.
// We cannot register a pair or its inverse more than once.
// If we have an FX pair then we implicitly have its inverse
// e.g. USDAUD = 1.0 / AUDUSD.
func (fxRates *FxRates) Register(rate *FxRate) error {
	pair := rate.GetPair()
	pair, err := ValidatePair(pair)
	if err != nil {
		return err
	}

	for _, fxRate := range fxRates.rates {
		registeredPair := btutil.CleanString(fxRate.GetPair())
		registeredInversePair, err := GetInversePair(registeredPair)
		if err != nil {
			return err
		}

		if pair == registeredPair || pair == registeredInversePair {
			return fmt.Errorf("'%s' fx rate instance already exists", pair)
		}
	}

	fxRates.rates = append(fxRates.rates, rate)
	return nil
}

// GetRate returns three items: an FX rate, a boolean value to show whether
// this is rate is available and an error.
func (fxRates *FxRates) GetRate(pair string) (float64, bool, error) {
	pair, err1 := ValidatePair(pair)
	isEquivalentPair, err2 := IsEquivalentPair(pair)
	inversePair, err3 := GetInversePair(pair)

	if err := btutil.AnyValidError(err1, err2, err3); err != nil {
		return 0.0, false, err
	}

	if isEquivalentPair { // e.g. AUDAUD, USDUSD, GBPGBP
		return 1.0, true, nil
	}

	for _, fxRate := range fxRates.rates {
		registeredPair := btutil.CleanString(fxRate.GetPair())
		// if we find the FX pair and it's valid then return the rate
		if registeredPair == pair {
			rate := fxRate.GetRate()
			if rate.Valid {
				if rate.Float64 == 0.0 {
					return 0.0, false, fmt.Errorf("'%s' FX rate is zero", pair)
				}
				return rate.Float64, true, nil
			}
		}

		// also look for the inverse pair
		if registeredPair == inversePair {
			rate := fxRate.GetRate()
			if rate.Valid {
				if rate.Float64 == 0.0 {
					return 0.0, false, fmt.Errorf("'%s' FX rate is zero", pair)
				}
				return 1 / rate.Float64, true, nil
			}
		}
	}

	return 0.0, false, nil
}

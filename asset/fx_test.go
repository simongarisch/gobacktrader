package asset

import (
	"gobacktrader/btutil"
	"testing"
)

func TestValidateCurrency(t *testing.T) {
	// first look at a valid currency code
	ccy, err := ValidateCurrency(" usd ")
	if ccy != "USD" {
		t.Errorf("Expecting currency 'USD', got '%s'", ccy)
	}
	if err != nil {
		t.Errorf("Unexpected error in ValidateCurrency - %s", err)
	}

	// then an invalid code
	_, err = ValidateCurrency("usda")
	if err == nil {
		t.Error("Expecting an error for invalid currency 'USDA'")
	}
	if err.Error() != "'USDA' is not a valid currency code" {
		t.Error("Unexpected error message string.")
	}
}

func TestValidatePair(t *testing.T) {
	// start with a valid pair
	pair, err := ValidatePair(" audusd ")
	if err != nil {
		t.Errorf("Error in ValidatePair - %s", err)
	}
	if pair != "AUDUSD" {
		t.Errorf("Expecting 'AUDUSD', got '%s'", pair)
	}

	// then an invalid pair
	pair, err = ValidatePair("AUDUSDX")
	if err.Error() != "expecting a six character currency pair, got 'AUDUSDX'" {
		t.Error("Unexpected error string")
	}
}

func TestSplitPair(t *testing.T) {
	ccy1, ccy2, err := SplitPair("AUDUSD")
	if err != nil {
		t.Errorf("Error in SplitPair - %s", err)
	}
	if ccy1 != "AUD" {
		t.Errorf("Expecting 'AUD' as ccy1, got '%s'", ccy1)
	}
	if ccy2 != "USD" {
		t.Errorf("Expecting 'USD' as ccy2, got '%s'", ccy2)
	}

	_, _, err = SplitPair("AUDUSDX")
	if err.Error() != "expecting a six character currency pair, got 'AUDUSDX'" {
		t.Error("Unexpected error string")
	}
}

func TestIsEquivalentPair(t *testing.T) {
	isEquivalent, err := IsEquivalentPair("AUDAUD")
	if err != nil {
		t.Errorf("Error in IsEquivalentPair - %s", err)
	}
	if !isEquivalent {
		t.Error("AUDAUD should be an equivalent pair")
	}

	isEquivalent, err = IsEquivalentPair("AUDUSD")
	if err != nil {
		t.Errorf("Error in IsEquivalentPair - %s", err)
	}
	if isEquivalent {
		t.Error("AUDUSD should not be an equivalent pair")
	}
}

func TestGetInversePair(t *testing.T) {
	inversePair, err := GetInversePair("AUDUSD")
	if err != nil {
		t.Errorf("Error in GetInversePair - %s", err)
	}
	if inversePair != "USDAUD" {
		t.Errorf("inverse pair of 'AUDUSD' is 'USDAUD', got '%s'", inversePair)
	}

	_, err = GetInversePair("AUDUSDX")
	if err.Error() != "expecting a six character currency pair, got 'AUDUSDX'" {
		t.Error("Unexpected error string")
	}
}

func TestFxRate(t *testing.T) {
	price := Price{Float64: 0.75, Valid: true}
	fxrate := NewFxRate("audusd", price)

	pair := fxrate.GetPair()
	if pair != "AUDUSD" {
		t.Errorf("Unexpected pair: wanted 'AUDUSD', got '%s'", pair)
	}

	rate := fxrate.GetRate()
	if rate.Float64 != 0.75 {
		t.Errorf("Unexpected rate: wanted 0.75, got %0.2f", rate.Float64)
	}

	price = Price{Float64: 0.8, Valid: true}
	fxrate.SetRate(price)
	rate = fxrate.GetRate()
	if rate.Float64 != 0.8 {
		t.Errorf("Unexpected rate: wanted 0.8, got %0.2f", rate.Float64)
	}
}

func TestFxRates(t *testing.T) {
	fxRates := FxRates{}

	price := Price{Float64: 0.75, Valid: true}
	audusd := NewFxRate("AUDUSD", price)
	err := fxRates.Register(&audusd)
	if err != nil {
		t.Errorf("Error in fxRates.Register - %s", err)
	}

	// check the AUDUSD rate
	rate, ok, err := fxRates.GetRate("AUDUSD")
	if err != nil {
		t.Errorf("Error in GetRate - %s", err)
	}
	if !ok {
		t.Error("'AUDUSD' rate should be available")
	}
	if rate != 0.75 {
		t.Errorf("Unexpected FX rate: wanted 0.75, got %0.4f", rate)
	}

	// the inverse rate of USDAUD
	rate, ok, err = fxRates.GetRate("USDAUD")
	if err != nil {
		t.Errorf("Error in GetRate - %s", err)
	}
	if !ok {
		t.Error("'USDAUD' rate should be available.")
	}

	actualRate := btutil.Round4dp(rate)
	expectedRate := btutil.Round4dp(1 / 0.75)
	if actualRate != expectedRate {
		t.Errorf("Unexpected inverse FX rate: wanted %0.4f, got %0.4f", expectedRate, actualRate)
	}
}

func TestFxRatesEmpty(t *testing.T) {
	fxRates := FxRates{}
	_, ok, err := fxRates.GetRate("AUDUSD")
	if ok {
		t.Error("'AUDUSD' rate shouldn't be available.")
	}
	if err != nil {
		t.Errorf("Error in GetRate - %s", err)
	}
}

func TestFxRatesRegistering(t *testing.T) {
	fxRates := FxRates{}
	xxxyyy := NewFxRate("XXXYYY", Price{Float64: 0.5, Valid: true})
	yyyxxx := NewFxRate("YYYXXX", Price{Float64: 2.0, Valid: true})

	var err error
	err = fxRates.Register(&xxxyyy)
	if err != nil {
		t.Errorf("Error in fxRates.Register - %s", err)
	}

	err = fxRates.Register(&yyyxxx) // we already implicitly have this rate
	if err.Error() != "'YYYXXX' fx rate instance already exists" {
		t.Error("Unexpected error when registering inverse rate.")
	}
}

func TestFxRatesEquivalentPairs(t *testing.T) {
	fxRates := FxRates{}
	equivalentPairs := []string{"AUDAUD", "USDUSD", "GBPGBP"}
	for _, pair := range equivalentPairs {
		rate, ok, err := fxRates.GetRate(pair)
		if err != nil {
			t.Errorf("Error in GetRate - %s", err)
		}
		if !ok {
			t.Errorf("'%s' should have an available rate.", pair)
		}
		if rate != 1.0 {
			t.Errorf("Expecting an FX rate of 1.0 for '%s'", pair)
		}
	}
}

func TestFxRateChanges(t *testing.T) {
	fxRates := FxRates{}
	startRate, endRate := 0.75, 0.80
	audusd := NewFxRate("AUDUSD", Price{Float64: startRate, Valid: true})

	err := fxRates.Register(&audusd)
	if err != nil {
		t.Errorf("Error in fxRates.Register - %s", err)
	}

	// check the starting rate
	rate, ok, err := fxRates.GetRate("AUDUSD")
	if err != nil {
		t.Errorf("Error in GetRate - %s", err)
	}
	if !ok {
		t.Error("'AUDUSD' should be an available rate")
	}
	if rate != startRate {
		t.Errorf("Unexpected FX rate: wanted %0.4f, got %0.4f", startRate, rate)
	}

	audusd.SetRate(Price{Float64: endRate, Valid: true})
	// check the ending rate
	rate, ok, err = fxRates.GetRate("AUDUSD")
	if err != nil {
		t.Errorf("Error in GetRate - %s", err)
	}
	if !ok {
		t.Error("'AUDUSD' should be an available rate")
	}
	if rate != endRate {
		t.Errorf("Unexpected FX rate: wanted %0.4f, got %0.4f", endRate, rate)
	}
}

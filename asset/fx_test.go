package asset

import (
	"gobacktrader/btutil"
	"testing"
	"time"
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
	if btutil.GetErrorString(err) != "'USDA' is not a valid currency code" {
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
	if btutil.GetErrorString(err) != "expecting a six character currency pair, got 'AUDUSDX'" {
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
	if btutil.GetErrorString(err) != "expecting a six character currency pair, got 'AUDUSDX'" {
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

	_, err = IsEquivalentPair("AUDUSDA")
	if btutil.GetErrorString(err) != "expecting a six character currency pair, got 'AUDUSDA'" {
		t.Error("Unexpected error string.")
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
	if btutil.GetErrorString(err) != "expecting a six character currency pair, got 'AUDUSDX'" {
		t.Error("Unexpected error string")
	}
}

func TestFxRate(t *testing.T) {
	price := Price{Float64: 0.75, Valid: true}
	fxrate, err := NewFxRate("audusd", price)
	if err != nil {
		t.Errorf("Error in NewFxRate - %s", err)
	}

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
	audusd, err := NewFxRate("AUDUSD", price)
	if err != nil {
		t.Errorf("Error in NewFxRate - %s", err)
	}
	err = fxRates.Register(&audusd)
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
	fxRates := NewFxRates()
	xxxyyy, _ := NewFxRate("XXXYYY", Price{Float64: 0.5, Valid: true})
	yyyxxx, _ := NewFxRate("YYYXXX", Price{Float64: 2.0, Valid: true})

	var err error
	err = fxRates.Register(&xxxyyy)
	if err != nil {
		t.Errorf("Error in fxRates.Register - %s", err)
	}

	err = fxRates.Register(&yyyxxx) // we already implicitly have this rate
	if btutil.GetErrorString(err) != "'YYYXXX' fx rate instance already exists" {
		t.Error("Unexpected error when registering inverse rate.")
	}

	// try an invalid FxRate
	fxRate := FxRate{pair: "AUDUSDA"}
	err = fxRates.Register(&fxRate)
	if btutil.GetErrorString(err) != "expecting a six character currency pair, got 'AUDUSDA'" {
		t.Error("Unexpected error string.")
	}
}

func TestFxRatesBadRegister(t *testing.T) {
	badRate := FxRate{pair: "AUDUSDA"}
	goodRate := FxRate{pair: "AUDUSD"}
	fxRates := FxRates{rates: []*FxRate{&badRate}}

	err := fxRates.Register(&goodRate)
	if btutil.GetErrorString(err) != "expecting a six character currency pair, got 'AUDUSDA'" {
		t.Error("Unexpected error string.")
	}
}

func TestFxRatesEquivalentPairs(t *testing.T) {
	fxRates := NewFxRates()
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

func TestZeroRates(t *testing.T) {
	fxRates := NewFxRates()
	fxRate, err := NewFxRate("AUDUSD", Price{Float64: 0.0, Valid: true})
	if err != nil {
		t.Errorf("Error in NewFxRate - %s", err)
	}
	fxRates.Register(&fxRate)

	_, _, err = fxRates.GetRate("AUDUSD")
	if btutil.GetErrorString(err) != "'AUDUSD' FX rate is zero" {
		t.Error("Unexpected error string.")
	}
	_, _, err = fxRates.GetRate("USDAUD")
	if btutil.GetErrorString(err) != "'USDAUD' FX rate is zero" {
		t.Error("Unexpected error string.")
	}
}

func TestGetRateInvalidPair(t *testing.T) {
	fxRates := NewFxRates()
	_, _, err := fxRates.GetRate("AUDUSDA")
	if btutil.GetErrorString(err) != "expecting a six character currency pair, got 'AUDUSDA'" {
		t.Error("Unexpected error string.")
	}
}

func TestFxRateChanges(t *testing.T) {
	fxRates := NewFxRates()
	startRate, endRate := 0.75, 0.80
	audusd, err := NewFxRate("AUDUSD", Price{Float64: startRate, Valid: true})
	if err != nil {
		t.Errorf("Error in NewFxRate - %s", err)
	}

	err = fxRates.Register(&audusd)
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

func TestFxHistory(t *testing.T) {
	fxrate, err := NewFxRate("audusd", nullPrice)
	if err != nil {
		t.Errorf("Error in NewFxRate - %s", err)
	}

	time1 := time.Date(2020, time.December, 14, 0, 0, 0, 0, time.UTC)
	time2 := time.Date(2020, time.December, 15, 0, 0, 0, 0, time.UTC)
	price1 := Price{Float64: 0.75, Valid: true}
	price2 := Price{Float64: 0.80, Valid: true}

	fxrate.SetPrice(price1)
	fxrate.TakeSnapshot(time1, &fxrate)
	fxrate.SetPrice(price2)
	fxrate.TakeSnapshot(time2, &fxrate)

	history := fxrate.GetHistory()

	// check our first snapshot
	snap1 := history[time1]
	if !snap1.GetTime().Equal(time1) {
		t.Error("snap1 - unexpected time.")
	}
	if snap1.GetPrice().Float64 != 0.75 {
		t.Error("snap1 - unexpected price.")
	}

	// and our second snapshot
	snap2 := history[time2]
	if !snap2.GetTime().Equal(time2) {
		t.Error("snap2 - unexpected time.")
	}
	if snap2.GetPrice().Float64 != 0.8 {
		t.Error("snap2 - unexpected price.")
	}
}

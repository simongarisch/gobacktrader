package asset

import (
	"gobacktrader/btutil"
	"testing"
	"time"
)

func TestNewPortfolio(t *testing.T) {
	p, err := NewPortfolio("XXX", "AUD")
	if err != nil {
		t.Errorf("Error in NewPortfolio - %s", err)
	}
	if p.GetCode() != "XXX" {
		t.Error("Unexpected portfolio code.")
	}
	if p.GetBaseCurrency() != "AUD" {
		t.Error("Unexpected base currency")
	}
	if p.NumPositions() != 0 {
		t.Error("New portfolios shouldn't have any positions.")
	}
}

func TestPortfolioValuation(t *testing.T) {
	p, err := NewPortfolio("XXX", "AUD")
	if err != nil {
		t.Errorf("Error in NewPortfolio - %s", err)
	}

	// portfolios should have a zero value to start with
	value, err := p.GetValue()
	if err != nil {
		t.Errorf("Error in GetValue - %s", err)
	}
	if !value.Valid {
		t.Error("Expecting a valid value.")
	}
	if value.Float64 != 0.0 {
		t.Error("Expecting a zero value.")
	}

	// Start adding positions
	a, err := NewStock("ZZB AU", "AUD")
	if err != nil {
		t.Errorf("Error in NewStock - %s", err)
	}
	if p.HasAsset(&a) {
		t.Error("The portfolio shouldn't have this asset position yet.")
	}

	p.ModifyPositions(&a, 100)
	if !p.HasAsset(&a) {
		t.Error("The portfolio should have this asset position.")
	}

	// the portfolio has 100 units of asset 'ZZB AU'
	// however, this asset doesn't yet have a valid price
	// and the portfolio value should still be invalid.
	value, err = p.GetValue()
	if err != nil {
		t.Errorf("Error in GetValue - %s", err)
	}
	if value.Valid {
		t.Error("Expecting an ivalid portfolio value.")
	}

	// we should have a valid portfolio value after setting
	// the price of 'ZZB AU'
	a.SetPrice(Price{Float64: 2.5, Valid: true})
	value, err = p.GetValue()
	if err != nil {
		t.Errorf("Error in GetValue - %s", err)
	}
	if !value.Valid {
		t.Error("We should have a valid portfolio value.")
	}

	portfolioValue := value.Float64
	expectedValue := 250.0
	if portfolioValue != expectedValue {
		t.Errorf("Unexpected portfolio value: wanted %0.2f, got %0.2f", expectedValue, portfolioValue)
	}
}

func TestGetWeight(t *testing.T) {
	p, err := NewPortfolio("XXX", "AUD")
	if err != nil {
		t.Errorf("Error in NewPortfolio - %s", err)
	}

	a, _ := NewStock("AAA AU", "AUD")
	_, err = p.GetWeight(&a)
	if btutil.GetErrorString(err) != "cannot calculate weights for portfolio with zero value" {
		t.Errorf("Unexpected error string.")
	}

	stock, err1 := NewStock("ZZB AU", "AUD")
	cash, err2 := NewCash("AUD")
	if err := btutil.AnyValidError(err1, err2); err != nil {
		t.Errorf("Error in portfolio asset creation - %s", err)
	}

	// put 100% of the portfolio weight in our stock
	stock.SetPrice(Price{Float64: 10.0, Valid: true})
	p.ModifyPositions(&stock, 100)
	weight, err := p.GetWeight(&stock)
	if err != nil {
		t.Errorf("Error in GetWeight - %s", err)
	}
	if weight.Float64 != 1.0 {
		t.Error("Expecting a stock weight of 100%")
	}

	weight, err = p.GetWeight(&cash)
	if err != nil {
		t.Errorf("Error in GetWeight - %s", err)
	}
	if weight.Float64 != 0.0 {
		t.Error("Expecting a zero weight in cash")
	}

	// add $200 cash by incrementing positions twice
	p.ModifyPositions(&cash, 150)
	p.ModifyPositions(&cash, 50)
	stock.SetPrice(Price{Float64: 8.0, Valid: true})

	// stock value = 100 * 8 = 800
	// cash value = 200
	// total value = 1000
	// stock weight = 80%, cash weight = 20%
	weight, err = p.GetWeight(&stock)
	if err != nil {
		t.Errorf("Error in GetWeight - %s", err)
	}
	if weight.Float64 != 0.8 {
		t.Error("Expecting a stock weight of 80%")
	}

	weight, err = p.GetWeight(&cash)
	if err != nil {
		t.Errorf("Error in GetWeight - %s", err)
	}
	if weight.Float64 != 0.2 {
		t.Error("Expecting a cash weight of 20%")
	}

	str := p.Show()
	expected := "---Portfolio('XXX')---\nAUD        200.00\nZZB AU     100.00\n"
	if str != expected {
		t.Errorf("Unexpected format in portfolio.Show")
	}
}

func TestPortfolioValuationCurrency(t *testing.T) {
	p1, err1 := NewPortfolio("XXX", "AUD")
	p2, err2 := NewPortfolio("YYY", "USD")
	if err := btutil.AnyValidError(err1, err2); err != nil {
		t.Errorf("Error in NewPortfolio - %s", err)
	}

	// create a stock
	stock, err1 := NewStock("ZZB AU", "AUD")
	cash, err2 := NewCash("AUD")
	if err := btutil.AnyValidError(err1, err2); err != nil {
		t.Errorf("Error in portfolio asset creation - %s", err)
	}

	// Give both p1 and p2 200 shares of stock and $100 cash.
	p1.ModifyPositions(&stock, 200)
	p1.ModifyPositions(&cash, 100)
	p2.ModifyPositions(&stock, 200)
	p2.ModifyPositions(&cash, 100)

	stock.SetPrice(Price{Float64: 2.5, Valid: true})

	fxRates := FxRates{}
	audusd, err := NewFxRate("AUDUSD", Price{Float64: 0.75, Valid: true})
	if err != nil {
		t.Errorf("Error in NewFxRate - %s", err)
	}

	fxRates.Register(&audusd)
	p1.SetFxRates(&fxRates)
	p2.SetFxRates(&fxRates)

	// this is an AUD denominated stock and AUD cash
	// AUD portfolio value = 200 * 2.50 + 100 = 600 AUD
	// USD portfolio value = 600 * 0.75 = 450 USD
	value, err := p1.GetValue()
	if err != nil {
		t.Errorf("Error in portfolio.GetValue - %s", err)
	}
	if !value.Valid {
		t.Error("Expecting a valid portfolio value.")
	}
	if value.Float64 != 600.0 {
		t.Errorf("Unexpected portfolio value: wanted 600, got %0.2f", value.Float64)
	}

	value, err = p2.GetValue()
	if err != nil {
		t.Errorf("Error in portfolio.GetValue - %s", err)
	}
	if !value.Valid {
		t.Error("Expecting a valid portfolio value.")
	}
	if value.Float64 != 450.0 {
		t.Errorf("Unexpected portfolio value: wanted 450, got %0.2f", value.Float64)
	}
}

func TestLargerPortfolio(t *testing.T) {
	p, err := NewPortfolio("XXX", "AUD")
	if err != nil {
		t.Errorf("Error in NewPortfolio - %s", err)
	}
	code := p.GetCode()
	if code != "XXX" {
		t.Errorf("Unexpected portfolio code: wanted 'XXX', got '%s'", code)
	}

	stock1, err1 := NewStock("ZZA AU", "AUD")
	stock2, err2 := NewStock("ZZU US", "USD")
	stock3, err3 := NewStock("ZZG AU", "GBP")
	if err := btutil.AnyValidError(err1, err2, err3); err != nil {
		t.Errorf("Error in NewStock - %s", err)
	}

	aud, err1 := NewCash("AUD")
	usd, err2 := NewCash("USD")
	gbp, err3 := NewCash("GBP")
	if err := btutil.AnyValidError(err1, err2, err3); err != nil {
		t.Errorf("Error in NewCash - %s", err)
	}

	// add 100 shares of each stock and $100 for each currency
	p.ModifyPositions(&stock1, 100)
	p.ModifyPositions(&stock2, 100)
	p.ModifyPositions(&stock3, 100)
	p.ModifyPositions(&aud, 100)
	p.ModifyPositions(&usd, 100)
	p.ModifyPositions(&gbp, 100)

	stock1.SetPrice(Price{Float64: 1.5, Valid: true})
	stock2.SetPrice(Price{Float64: 2.5, Valid: true})
	stock3.SetPrice(Price{Float64: 3.5, Valid: true})

	fxRates := FxRates{}
	audusd, err1 := NewFxRate("AUDUSD", Price{Float64: 0.75, Valid: true})
	gbpaud, err2 := NewFxRate("GBPAUD", Price{Float64: 1.80, Valid: true})
	if err := btutil.AnyValidError(err1, err2); err != nil {
		t.Errorf("Error in NewFxRate - %s", err)
	}

	fxRates.Register(&audusd)
	fxRates.Register(&gbpaud)
	p.SetFxRates(&fxRates)

	// this portfolio has a base currency of AUD
	// stock value = (1.5 * 100) + (2.5 * 100) / 0.75 + (3.5 * 100) * 1.80 = 1113.3333
	// cash value = 100 + 100 / 0.75 + 100 * 1.8 = 413.3333
	// total value = AUD 1526.6667
	value, err := p.GetValue()
	if err != nil {
		t.Errorf("Error in portfolio.GetValue - %s", err)
	}
	if !value.Valid {
		t.Error("Expecting a valid portfolio value.")
	}

	actualValue := btutil.Round2dp(value.Float64)
	expectedValue := 1526.67
	if actualValue != expectedValue {
		t.Errorf("Unexpected portfolio value: wanted %0.2f, got %0.2f", expectedValue, actualValue)
	}

	// increase a stock price and check this flows through.
	stock1.SetPrice(Price{Float64: 2.5, Valid: true}) // increasing price by $1 (will add $100 AUD to portfolio value)
	value, err = p.GetValue()
	if err != nil {
		t.Errorf("Error in portfolio.GetValue - %s", err)
	}
	if !value.Valid {
		t.Error("Expecting a valid portfolio value.")
	}

	actualValue = btutil.Round2dp(value.Float64)
	expectedValue = 1626.67
	if actualValue != expectedValue {
		t.Errorf("Unexpected portfolio value: wanted %0.2f, got %0.2f", expectedValue, actualValue)
	}
}

func TestPortfolioUnitsWeight(t *testing.T) {
	p, err := NewPortfolio("XXX", "AUD")
	if err != nil {
		t.Errorf("Error in NewPortfolio - %s", err)
	}

	stock1, err1 := NewStock("ZZA AU", "AUD")
	stock2, err2 := NewStock("ZZU US", "USD")
	if err := btutil.AnyValidError(err1, err2); err != nil {
		t.Errorf("Error in NewStock - %s", err)
	}

	aud, err1 := NewCash("AUD")
	usd, err2 := NewCash("USD")
	if err := btutil.AnyValidError(err1, err2); err != nil {
		t.Errorf("Error in NewCash - %s", err)
	}

	// this portfolio is currently empty
	if p.NumPositions() != 0 {
		t.Errorf("Expecting zero positions for an initialised portfolio.")
	}
	if p.GetUnits(&stock1) != 0.0 {
		t.Error("Expecting zero units in stock1.")
	}
	if p.GetUnits(&aud) != 0.0 {
		t.Error("Expecting zero units in aud cash.")
	}

	// add 100 shares of each stock and $100 for each currency
	p.ModifyPositions(&stock1, 100)
	p.ModifyPositions(&stock2, 100)
	p.ModifyPositions(&aud, 100)
	p.ModifyPositions(&usd, 100)

	stock1.SetPrice(Price{Float64: 1.5, Valid: true})
	stock2.SetPrice(Price{Float64: 2.5, Valid: true})

	fxRates := FxRates{}
	audusd, err := NewFxRate("AUDUSD", Price{Float64: 0.75, Valid: true})
	if err != nil {
		t.Errorf("Error in NewFxRate - %s", err)
	}

	fxRates.Register(&audusd)
	p.SetFxRates(&fxRates)

	// this portfolio has a base currency of AUD
	// stock value = (1.5 * 100) + (2.5 * 100) / 0.75 = 483.3333
	// cash value = 100 + 100 / 0.75 = 233.3333
	// total value = AUD 716.6667
	value, err := p.GetValue()
	if err != nil {
		t.Errorf("Error in portfolio.GetValue - %s", err)
	}
	if !value.Valid {
		t.Error("Expecting a valid portfolio value.")
	}

	actualValue := btutil.Round2dp(value.Float64)
	expectedValue := 716.67
	if actualValue != expectedValue {
		t.Errorf("Unexpected portfolio value: wanted %0.2f, got %0.2f", expectedValue, actualValue)
	}

	// check the units
	for _, targetAsset := range []IAssetReadOnly{&stock1, &stock2, &aud, &usd} {
		if p.GetUnits(targetAsset) != 100 {
			t.Errorf("'%s': expecting to hold 100 units", targetAsset.GetTicker())
		}
	}

	// and the weights
	// stock1 weight = (1.5 * 100) / 716.6667 = 0.2093
	// stock2 weight = ((2.5 * 100) / 0.75) / 716.6667 = 0.4651
	// aud weight = 100 / 716.6667 = 0.1395
	// usd weight = (100 / 0.75) / 716.6667 = 0.1860
	// sum of weights = 100%
	wStock1, err1 := p.GetWeight(&stock1)
	wStock2, err2 := p.GetWeight(&stock2)
	wAud, err3 := p.GetWeight(&aud)
	wUsd, err4 := p.GetWeight(&usd)
	if err := btutil.AnyValidError(err1, err3, err3, err4); err != nil {
		t.Errorf("Error in portfolio.GetWeight - %s", err)
	}

	if btutil.Round4dp(wStock1.Float64) != 0.2093 {
		t.Errorf("'%s' expecting a portfolio weight of 0.2093, got %0.4f", stock1.GetTicker(), wStock1.Float64)
	}
	if btutil.Round4dp(wStock2.Float64) != 0.4651 {
		t.Errorf("'%s' expecting a portfolio weight of 0.4651, got %0.4f", stock2.GetTicker(), wStock2.Float64)
	}
	if btutil.Round4dp(wAud.Float64) != 0.1395 {
		t.Errorf("'%s' expecting a portfolio weight of 0.4651, got %0.4f", aud.GetTicker(), wAud.Float64)
	}
	if btutil.Round4dp(wUsd.Float64) != 0.1860 {
		t.Errorf("'%s' expecting a portfolio weight of 0.1860, got %0.4f", usd.GetTicker(), wUsd.Float64)
	}
}

func TestPortfolioSnapshots(t *testing.T) {
	p, err := NewPortfolio("XXX", "AUD")
	if err != nil {
		t.Errorf("Error in NewPortfolio - %s", err)
	}

	stock, err1 := NewStock("ZZA AU", "AUD")
	cash, err2 := NewCash("AUD")
	if err := btutil.AnyValidError(err1, err2); err != nil {
		t.Errorf("Asset error - %s", err)
	}

	// add 100 units of each
	p.ModifyPositions(&stock, 100)
	p.ModifyPositions(&cash, 100)
	stock.SetPrice(Price{Float64: 1.5, Valid: true})

	// portfolio value is 1.50 * 100 + 100 = 250 AUD
	// wStock = 150 / 250 = 60%
	// wCash = 100 / 250 = 40%
	t1 := time.Date(2020, time.December, 14, 0, 0, 0, 0, time.UTC)
	p.TakeSnapshot(t1)

	stock.SetPrice(Price{Float64: 3.0, Valid: true})
	// portfolio value is 3.00 * 100 + 100 = 400 AUD
	// wStock = 300 / 400 = 75%
	// wCash = 100 / 400 = 25%
	t2 := time.Date(2020, time.December, 15, 0, 0, 0, 0, time.UTC)
	p.TakeSnapshot(t2)

	history := p.GetHistory()
	snap1, _ := history[t1]
	snap2, _ := history[t2]

	// check the first snapshot
	if !snap1.GetTime().Equal(t1) {
		t.Error("Unexpected time for snap1")
	}
	if snap1.GetValue().Float64 != 250 {
		t.Errorf("Unexpected value: wanted 250.00, got %0.2f", snap1.GetValue().Float64)
	}
	weights := snap1.GetWeights()
	wStock, _ := weights[&stock]
	wCash, _ := weights[&cash]
	if wStock.Float64 != 0.6 {
		t.Errorf("Unexpected stock weight: wanted 0.60, got %0.2f", wStock.Float64)
	}
	if wCash.Float64 != 0.4 {
		t.Errorf("Unexpected cash weight: wanted 0.40, got %0.2f", wCash.Float64)
	}

	// and the second
	if !snap2.GetTime().Equal(t2) {
		t.Error("Unexpected time for snap2")
	}
	if snap2.GetValue().Float64 != 400 {
		t.Errorf("Unexpected value: wanted 400.00, got %0.2f", snap2.GetValue().Float64)
	}
	weights = snap2.GetWeights()
	wStock, _ = weights[&stock]
	wCash, _ = weights[&cash]
	if wStock.Float64 != 0.75 {
		t.Errorf("Unexpected stock weight: wanted 0.75, got %0.2f", wStock.Float64)
	}
	if wCash.Float64 != 0.25 {
		t.Errorf("Unexpected cash weight: wanted 0.25, got %0.2f", wCash.Float64)
	}
}

func TestPortfolioNoFxRate(t *testing.T) {
	p, err1 := NewPortfolio("XXX", "AUD")
	stock, err2 := NewStock("ZZA US", "USD")
	audusd, err3 := NewFxRate("AUDUSD", Price{Float64: 0.75, Valid: true})
	if err := btutil.AnyValidError(err1, err2, err3); err != nil {
		t.Errorf("Error in asset creation - %s", err)
	}

	stock.SetPrice(Price{Float64: 3.0, Valid: true})
	fxRates := FxRates{}
	fxRates.Register(&audusd)

	// transfer 100 shares of stock
	// note the portfolio fx rates is currently empty
	p.ModifyPositions(&stock, 100)

	value, err := p.GetValue()
	if err != nil {
		t.Errorf("Error in portfolio.GetValue - %s", err)
	}
	if value != nullPrice {
		t.Error("Expecting a null price with no FX rate for valuation.")
	}

	// set the fx rates so we can get a valuation
	p.SetFxRates(&fxRates)
	value, err = p.GetValue()
	if err != nil {
		t.Errorf("Error in portfolio.GetValue - %s", err)
	}
	if value.Float64 != 400 { // 300 / 0.75 = 400
		t.Errorf("Expecting a valuation of $400, got %0.2f", value.Float64)
	}
}

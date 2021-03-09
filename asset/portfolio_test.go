package asset

import (
	"gobacktrader/btutil"
	"testing"
)

func TestNewPortfolio(t *testing.T) {
	p := NewPortfolio("XXX", "AUD")
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
	p := NewPortfolio("XXX", "AUD")

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

func TestPortfolioValuationCurrency(t *testing.T) {
	p1 := NewPortfolio("XXX", "AUD")
	p2 := NewPortfolio("YYY", "USD")

	// create a stock
	stock, err1 := NewStock("ZZB AU", "AUD")
	cash, err2 := NewCash("AUD")
	for _, err := range []error{err1, err2} {
		if err != nil {
			t.Errorf("Error in portfolio asset creation - %s", err)
		}
	}

	// Give both p1 and p2 200 shares of stock and $100 cash.
	p1.ModifyPositions(&stock, 200)
	p1.ModifyPositions(&cash, 100)
	p2.ModifyPositions(&stock, 200)
	p2.ModifyPositions(&cash, 100)

	stock.SetPrice(Price{Float64: 2.5, Valid: true})

	fxRates := FxRates{}
	audusd := NewFxRate("AUDUSD", Price{Float64: 0.75, Valid: true})
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
	p := NewPortfolio("XXX", "AUD")
	code := p.GetCode()
	if code != "XXX" {
		t.Errorf("Unexpected portfolio code: wanted 'XXX', got '%s'", code)
	}

	stock1, err1 := NewStock("ZZA AU", "AUD")
	stock2, err2 := NewStock("ZZU US", "USD")
	stock3, err3 := NewStock("ZZG AU", "GBP")
	for _, err := range []error{err1, err2, err3} {
		if err != nil {
			t.Errorf("Error in NewStock - %s", err)
		}
	}

	aud, err1 := NewCash("AUD")
	usd, err2 := NewCash("USD")
	gbp, err3 := NewCash("GBP")
	for _, err := range []error{err1, err2, err3} {
		if err != nil {
			t.Errorf("Error in NewCash - %s", err)
		}
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
	audusd := NewFxRate("AUDUSD", Price{Float64: 0.75, Valid: true})
	gbpaud := NewFxRate("GBPAUD", Price{Float64: 1.80, Valid: true})
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

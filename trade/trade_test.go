package trade

import (
	"gobacktrader/asset"
	"gobacktrader/btutil"
	"gobacktrader/compliance"
	"testing"
)

func TestTradeInit(t *testing.T) {
	portfolio, err1 := asset.NewPortfolio("XXX", "AUD")
	stock, err2 := asset.NewStock("ZZB AU", "AUD")
	aud, err3 := asset.NewCash("AUD")
	if err := btutil.AnyValidError(err1, err2, err3); err != nil {
		t.Errorf("Error in asset init - %s", err)
	}

	trade := Trade{
		portfolio:   portfolio,
		targetAsset: stock,
		units:       100.0,
	}

	if trade.GetPortfolio() != portfolio {
		t.Error("Unexpected portfolio")
	}
	if trade.GetAsset() != stock {
		t.Error("Unexpected asset")
	}
	if trade.GetUnits() != 100 {
		t.Error("Unexpected units")
	}
	if cash, _ := trade.GetBaseCurrencyCash(); cash != aud {
		t.Error("Unexpected base currency cash")
	}
}

func TestGetLocalCurrencyValue(t *testing.T) {
	portfolio, err1 := asset.NewPortfolio("XXX", "AUD")
	stock, err2 := asset.NewStock("ZZB AU", "AUD")
	if err := btutil.AnyValidError(err1, err2); err != nil {
		t.Errorf("Error in asset init - %s", err)
	}

	trade1 := Trade{
		portfolio:   portfolio,
		targetAsset: stock,
		units:       +100.0,
	}

	trade2 := Trade{
		portfolio:   portfolio,
		targetAsset: stock,
		units:       -100.0,
	}

	var value asset.Price
	// without a price the trade value should be invalid
	value = trade1.GetLocalCurrencyValue()
	if value.Valid {
		t.Error("Trade value should be invalid")
	}
	value = trade2.GetLocalCurrencyValue()
	if value.Valid {
		t.Error("Trade value should be invalid")
	}

	// check the trade value
	// the absolute value should be returned
	stock.SetPrice(asset.Price{Float64: 2.50, Valid: true})
	value = trade1.GetLocalCurrencyValue()
	if !value.Valid {
		t.Error("Trade should have a valid value")
	}
	if value.Float64 != +250 {
		t.Errorf("Unexpected trade value: wanted 250, got %0.2f", value.Float64)
	}

	value = trade2.GetLocalCurrencyValue()
	if !value.Valid {
		t.Error("Trade should have a valid value")
	}
	if value.Float64 != +250 {
		t.Errorf("Unexpected trade value: wanted 250, got %0.2f", value.Float64)
	}
}

func TestGetLocalCurrencyConsideration(t *testing.T) {
	portfolio, err1 := asset.NewPortfolio("XXX", "AUD")
	stock, err2 := asset.NewStock("ZZB AU", "AUD")
	if err := btutil.AnyValidError(err1, err2); err != nil {
		t.Errorf("Error in asset init - %s", err)
	}

	buyTrade := NewTrade(portfolio, stock, +100.0)
	sellTrade := NewTrade(portfolio, stock, -100.0)

	// consideration should be invalid where we don't
	// have a price.
	var consideration asset.Price
	consideration = buyTrade.GetLocalCurrencyConsideration()
	if consideration.Valid {
		t.Error("Consideration should be invalid")
	}
	consideration = sellTrade.GetLocalCurrencyConsideration()
	if consideration.Valid {
		t.Error("Consideration should be invalid")
	}

	// with a price the consideration should be valid
	// -ve for buys as cash is going out
	// +ve for sells as cash is coming in
	stock.SetPrice(asset.Price{Float64: 2.50, Valid: true})
	consideration = buyTrade.GetLocalCurrencyConsideration()
	if !consideration.Valid {
		t.Error("We should have a valid consideration")
	}
	if consideration.Float64 != -250 {
		t.Errorf("Unexpected consideration: wanted -250, got %0.2f", consideration.Float64)
	}

	consideration = sellTrade.GetLocalCurrencyConsideration()
	if !consideration.Valid {
		t.Error("We should have a valid consideration")
	}
	if consideration.Float64 != +250 {
		t.Errorf("Unexpected consideration: wanted +250, got %0.2f", consideration.Float64)
	}
}

func TestTradePassesCompliance(t *testing.T) {
	portfolio, err1 := asset.NewPortfolio("XXX", "AUD")
	stock, err2 := asset.NewStock("ZZB AU", "AUD")
	if err := btutil.AnyValidError(err1, err2); err != nil {
		t.Errorf("Error in asset init - %s", err)
	}

	// with no compliance rules in place the portfolio should pass.
	passesCompliance, err := portfolio.PassesCompliance()
	if err != nil {
		t.Errorf("Error in portfolio.PassesComplinace() - %s", err)
	}
	if !passesCompliance {
		t.Error("Expecting a portfolio with no compliance rules to pass")
	}

	// no more than 100% of the portfolio in this stock
	stockLimit := compliance.NewWeightLimit(stock, 1.0)
	portfolio.AddComplianceRule(stockLimit)

	// without a stock price or portfolio value we cannot calculate weights
	// PassesCompliance should return an error

	trade := NewTrade(portfolio, stock, +1000)
	_, err = trade.PassesCompliance()
	errStr := btutil.GetErrorString(err)
	if errStr != "portfolio has no assigned executing broker" {
		t.Errorf("Unexpected error string - %s", err)
	}
}

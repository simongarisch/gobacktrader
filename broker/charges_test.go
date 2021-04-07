package broker

import (
	"gobacktrader/asset"
	"gobacktrader/btutil"
	"gobacktrader/trade"
	"testing"
)

func TestNoCharges(t *testing.T) {
	portfolio, err1 := asset.NewPortfolio("XXX", "AUD")
	stock, err2 := asset.NewStock("ZZB AU", "AUD")
	cash, err3 := asset.NewCash("AUD")
	if err := btutil.AnyValidError(err1, err2, err3); err != nil {
		t.Errorf("Error in asset init - %s", err)
	}

	portfolio.Transfer(cash, 1000.0)
	buyTrade := trade.NewTrade(portfolio, stock, +100.0)
	sellTrade := trade.NewTrade(portfolio, stock, -100.0)

	charges := NewNoCharges()
	if portfolio.GetUnits(cash) != 1000 {
		t.Error("Expecting 1K AUD cash")
	}

	err := charges.Charge(buyTrade) // no charge, so cash should be static
	if err != nil {
		t.Errorf("Error in NoCharges{}.Charge() - %s", err)
	}
	if portfolio.GetUnits(cash) != 1000 {
		t.Error("Expecting 1K AUD cash")
	}

	err = charges.Charge(sellTrade) // no charge, so cash should be static
	if err != nil {
		t.Errorf("Error in NoCharges{}.Charge() - %s", err)
	}
	if portfolio.GetUnits(cash) != 1000 {
		t.Error("Expecting 1K AUD cash")
	}
}

func TestFixedRatePlusPercentageChargesBasic(t *testing.T) {
	portfolio, err1 := asset.NewPortfolio("XXX", "AUD")
	stock, err2 := asset.NewStock("ZZB AU", "AUD")
	cash, err3 := asset.NewCash("AUD")
	if err := btutil.AnyValidError(err1, err2, err3); err != nil {
		t.Errorf("Error in asset init - %s", err)
	}

	portfolio.Transfer(cash, 1000.0)
	buyTrade := trade.NewTrade(portfolio, stock, +100.0)
	sellTrade := trade.NewTrade(portfolio, stock, -100.0)

	// we should be starting with 1K AUD cash
	if portfolio.GetUnits(cash) != 1000 {
		t.Error("Expecting 1K AUD cash")
	}

	// apply a fixed rate of $20 plus 1% of the trade
	charges, err := NewFixedRatePlusPercentageCharges(20, 0.01, "AUD")
	if err != nil {
		t.Errorf("Error in NewFixedRatePlusPercentageCharges - %s", err)
	}

	// where an asset has no valid price this should return an error
	err1 = charges.Charge(buyTrade)
	err2 = charges.Charge(sellTrade)
	for _, err := range []error{err1, err2} {
		errStr := btutil.GetErrorString(err)
		if errStr != "cannot apply charges to a trade with invalid value" {
			t.Errorf("Unexpected error string - %s", err)
		}
	}

	stock.SetPrice(asset.Price{Float64: 2.50, Valid: true})

	// charges per trade will be 20 + 2.50*100*1% = $22.5
	// 1000 - 22.5 = 977.5
	// 1000 - 22.5 * 2 = 955
	err = charges.Charge(buyTrade)
	if err != nil {
		t.Errorf("Error in Charge - %s", err)
	}
	portfolioAUD := portfolio.GetUnits(cash)
	if portfolioAUD != 977.5 {
		t.Errorf("Unexpected AUD cash - %0.2f", portfolioAUD)
	}

	err = charges.Charge(sellTrade)
	if err != nil {
		t.Errorf("Error in Charge - %s", err)
	}
	portfolioAUD = portfolio.GetUnits(cash)
	if portfolioAUD != 955.0 {
		t.Errorf("Unexpected AUD cash - %0.2f", portfolioAUD)
	}
}

func TestFixedRatePlusPercentageChargesUSD(t *testing.T) {
	portfolio, err1 := asset.NewPortfolio("XXX", "AUD")
	stock, err2 := asset.NewStock("ZZB AU", "AUD")
	aud, err3 := asset.NewCash("AUD")
	usd, err4 := asset.NewCash("USD")
	if err := btutil.AnyValidError(err1, err2, err3, err4); err != nil {
		t.Errorf("Error in asset init - %s", err)
	}

	portfolio.Transfer(aud, 1000.0)
	buyTrade := trade.NewTrade(portfolio, stock, +100.0)
	sellTrade := trade.NewTrade(portfolio, stock, -100.0)

	// we should be starting with 1K AUD cash
	if portfolio.GetUnits(aud) != 1000 {
		t.Error("Expecting 1K AUD cash")
	}

	// apply a fixed rate of $20 plus 1% of the trade (in USD)
	charges, err := NewFixedRatePlusPercentageCharges(20, 0.01, "USD")
	if err != nil {
		t.Errorf("Error in NewFixedRatePlusPercentageCharges - %s", err)
	}

	// where an asset has no valid price this should return an error
	err1 = charges.Charge(buyTrade)
	err2 = charges.Charge(sellTrade)
	for _, err := range []error{err1, err2} {
		errStr := btutil.GetErrorString(err)
		if errStr != "cannot apply charges to a trade with invalid value" {
			t.Errorf("Unexpected error string - %s", err)
		}
	}

	stock.SetPrice(asset.Price{Float64: 2.50, Valid: true})

	fxRates := asset.NewFxRates()
	audusd, err := asset.NewFxRate("AUDUSD", asset.Price{Float64: 0.75, Valid: true})
	if err != nil {
		t.Errorf("Error in asset.NewFxRate - %s", err)
	}
	fxRates.Register(audusd)
	portfolio.SetFxRates(fxRates)

	// charges per trade (in USD) will be 20 + (2.50*100*1%)*0.75 = USD 21.875
	// once 21.875
	// twice 43.75
	err = charges.Charge(buyTrade)
	if err != nil {
		t.Errorf("Error in Charge - %s", err)
	}
	if portfolio.GetUnits(aud) != 1000 {
		t.Errorf("Unexpected AUD cash - %0.2f", portfolio.GetUnits(aud))
	}
	if portfolio.GetUnits(usd) != -21.875 {
		t.Errorf("Unexpected USD cash - %0.2f", portfolio.GetUnits(usd))
	}

	err = charges.Charge(sellTrade)
	if err != nil {
		t.Errorf("Error in Charge - %s", err)
	}
	if portfolio.GetUnits(aud) != 1000 {
		t.Errorf("Unexpected AUD cash - %0.2f", portfolio.GetUnits(aud))
	}
	if portfolio.GetUnits(usd) != -43.75 {
		t.Errorf("Unexpected USD cash - %0.2f", portfolio.GetUnits(usd))
	}
}

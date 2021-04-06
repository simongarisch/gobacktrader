package backtest

import (
	"gobacktrader/asset"
	"gobacktrader/btutil"
	"testing"
)

func TestRegisterPortfolio(t *testing.T) {
	backtest := NewBacktest()

	portfolio, err := asset.NewPortfolio("XXX", "AUD")
	if err != nil {
		t.Errorf("Error in NewPortfolio - %s", err)
	}

	stock, err := asset.NewStock("ZZB AU", "AUD")
	if err != nil {
		t.Errorf("Error in NewStock - %s", err)
	}

	if backtest.HasPortfolio(portfolio) {
		t.Error("Backtest should not have portfolio registered.")
	}
	if backtest.HasAsset(stock) {
		t.Error("Backtest should not have asset registered.")
	}

	err1 := backtest.RegisterPortfolio(portfolio)
	err2 := backtest.RegisterAsset(stock)
	if err := btutil.AnyValidError(err1, err2); err != nil {
		t.Errorf("Error in backtest.Register - %s", err)
	}

	if !backtest.HasPortfolio(portfolio) {
		t.Error("Backtest should have portfolio registered.")
	}
	if !backtest.HasAsset(stock) {
		t.Error("Backtest should have asset registered.")
	}

	// we should be able to register the same asset and portfolio
	// again without issue
	err1 = backtest.RegisterPortfolio(portfolio)
	err2 = backtest.RegisterAsset(stock)
	if err := btutil.AnyValidError(err1, err2); err != nil {
		t.Errorf("Error in backtest.Register - %s", err)
	}

	// try to register a different portfolio and asset with the same codes
	portfolio2, err := asset.NewPortfolio("XXX", "USD")
	if err != nil {
		t.Errorf("Error in NewPortfolio - %s", err)
	}

	stock2, err := asset.NewStock("ZZB AU", "USD")
	if err != nil {
		t.Errorf("Error in NewStock - %s", err)
	}

	err1 = backtest.RegisterPortfolio(portfolio2)
	err2 = backtest.RegisterAsset(stock2)
	errStr1 := btutil.GetErrorString(err1)
	errStr2 := btutil.GetErrorString(err2)
	if errStr1 != "portfolio code 'XXX' is already in use and needs to be unique" {
		t.Errorf("Unexpected error string '%s'", errStr1)
	}
	if errStr2 != "asset ticker 'ZZB AU' is already in use and needs to be unique" {
		t.Errorf("Unexpected error string '%s'", errStr2)
	}
}

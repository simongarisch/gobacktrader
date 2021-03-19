package backtest

import (
	"gobacktrader/asset"
	"gobacktrader/btutil"
	"testing"
)

func TestRegisterPortfolio(t *testing.T) {
	backtest := NewBacktest()

	p, err := asset.NewPortfolio("XXX", "AUD")
	if err != nil {
		t.Errorf("Error in NewPortfolio - %s", err)
	}

	a, err := asset.NewStock("ZZB AU", "AUD")
	if err != nil {
		t.Errorf("Error in NewStock - %s", err)
	}

	if backtest.HasPortfolio(&p) {
		t.Error("Backtest should not have portfolio registered.")
	}
	if backtest.HasAsset(&a) {
		t.Error("Backtest should not have asset registered.")
	}

	err1 := backtest.RegisterPortfolio(&p)
	err2 := backtest.RegisterAsset(&a)
	if err := btutil.AnyValidError(err1, err2); err != nil {
		t.Errorf("Error in backtest.Register - %s", err)
	}

	if !backtest.HasPortfolio(&p) {
		t.Error("Backtest should have portfolio registered.")
	}
	if !backtest.HasAsset(&a) {
		t.Error("Backtest should have asset registered.")
	}

	// we should be able to register the same asset and portfolio
	// again without issue
	err1 = backtest.RegisterPortfolio(&p)
	err2 = backtest.RegisterAsset(&a)
	if err := btutil.AnyValidError(err1, err2); err != nil {
		t.Errorf("Error in backtest.Register - %s", err)
	}

	// try to register a different portfolio and asset with the same codes
	p2, err := asset.NewPortfolio("XXX", "USD")
	if err != nil {
		t.Errorf("Error in NewPortfolio - %s", err)
	}

	a2, err := asset.NewStock("ZZB AU", "USD")
	if err != nil {
		t.Errorf("Error in NewStock - %s", err)
	}

	err1 = backtest.RegisterPortfolio(&p2)
	err2 = backtest.RegisterAsset(&a2)
	if err1.Error() != "portfolio code 'XXX' is already in use and needs to be unique" {
		t.Errorf("Unexpected error string '%s'", err1.Error())
	}
	if err2.Error() != "asset ticker 'ZZB AU' is already in use and needs to be unique" {
		t.Errorf("Unexpected error string '%s'", err2.Error())
	}
}

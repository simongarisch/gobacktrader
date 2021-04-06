package compliance

import (
	"gobacktrader/asset"
	"gobacktrader/btutil"
	"testing"
)

func TestUnitLimit(t *testing.T) {
	portfolio, err1 := asset.NewPortfolio("XXX", "AUD")
	cash, err2 := asset.NewCash("AUD")
	stock, err3 := asset.NewStock("ZZB AU", "AUD")
	if err := btutil.AnyValidError(err1, err2, err3); err != nil {
		t.Errorf("Error in asset creation - %s", err)
	}

	// transfer 100 shares of stock and 100 AUD to the portfolio
	portfolio.Transfer(stock, 100)
	portfolio.Transfer(cash, 100)

	cashLimit := NewUnitLimit(cash, 100)
	stockLimit := NewUnitLimit(stock, 99)

	if cashLimit.GetAsset() != cash {
		t.Error("Unexpected asset")
	}
	if stockLimit.GetAsset() != stock {
		t.Error("Unexpected asset")
	}
	if cashLimit.GetLimit() != 100 {
		t.Error("Unexpected cash limit")
	}
	if stockLimit.GetLimit() != 99 {
		t.Error("Unexpected stock limit")
	}

	cashPasses, err1 := cashLimit.Passes(portfolio)
	stockPasses, err2 := stockLimit.Passes(portfolio)
	if err := btutil.AnyValidError(err1, err2); err != nil {
		t.Errorf("Error in UnitLimit{}.Passes() - %s", err)
	}

	if cashPasses != true {
		t.Error("Expecting cash to pass the unit limit")
	}
	if stockPasses != false {
		t.Error("Expecting stock to fail the unit limit")
	}
}

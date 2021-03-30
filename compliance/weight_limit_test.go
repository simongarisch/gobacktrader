package compliance

import (
	"gobacktrader/asset"
	"gobacktrader/btutil"
	"testing"
)

func TestWeightLimit(t *testing.T) {
	portfolio, err1 := asset.NewPortfolio("XXX", "AUD")
	cash, err2 := asset.NewCash("AUD")
	stock, err3 := asset.NewStock("ZZB AU", "AUD")
	if err := btutil.AnyValidError(err1, err2, err3); err != nil {
		t.Errorf("Error in asset creation - %s", err)
	}

	cashLimit := NewWeightLimit(&cash, 0.5)
	stockLimit := NewWeightLimit(&stock, 0.6)

	if cashLimit.GetAsset() != &cash {
		t.Error("Unexpected asset")
	}
	if stockLimit.GetAsset() != &stock {
		t.Error("Unexpected asset")
	}
	if cashLimit.GetLimit() != 0.5 {
		t.Error("Unexpected cash weight limit")
	}
	if stockLimit.GetLimit() != 0.6 {
		t.Error("Unexpected stock weight limit")
	}

	// the portfolio will currently have a zero value
	// so we cannot calculate portfolio weights.
	cashPasses, err := cashLimit.Passes(&portfolio)
	if cashPasses != false {
		t.Error("Expecting cash weight limit not to pass")
	}
	errStr := btutil.GetErrorString(err)
	if errStr != "cannot calculate weights for portfolio with zero value" {
		t.Errorf("Unexpected error string, got '%s'", errStr)
	}

	// transfer 100 shares of stock and 100 AUD to the portfolio
	portfolio.ModifyPositions(&stock, 100)
	portfolio.ModifyPositions(&cash, 100)

	// we have transferred assets to the portfolio,
	// but we still haven't set the price of our stock
	cashPasses, err = cashLimit.Passes(&portfolio) // weight not valid
	if cashPasses != false {
		t.Error("Expecting cash weight limit not to pass")
	}
	if err != nil {
		t.Errorf("Error in WeightLimit{}.Passes() - %s", err)
	}

	// set the price of our stock to 1.50 and check limits
	// cash value = 100, stock value = 100 * 1.50 = 150
	// cash weight = 100 / 250 = 40%
	// stock weight = 150 / 250 = 60%
	// they are both right on their limits and should pass
	stock.SetPrice(asset.Price{Float64: 1.50, Valid: true})
	cashPasses, err1 = cashLimit.Passes(&portfolio)
	stockPasses, err2 := stockLimit.Passes(&portfolio)
	if err := btutil.AnyValidError(err1, err2); err != nil {
		t.Errorf("Error in WeightLimit{}.Passes() - %s", err)
	}
	if !cashPasses {
		t.Errorf("Expecting the cash weight limit to pass")
	}
	if !stockPasses {
		t.Errorf("Expecting the stock weight limit to pass")
	}

	// given we are right on our limits adding one share of stock
	// to the portfolio should cause the stock weight limit to fail.
	portfolio.ModifyPositions(&stock, 1)
	stockPasses, err = stockLimit.Passes(&portfolio)
	if err != nil {
		t.Errorf("Error in WeightLimit{}.Passes() - %s", err)
	}
	if stockPasses {
		t.Error("Expecting the stock weight limit to fail")
	}
}

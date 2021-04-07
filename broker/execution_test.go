package broker

import (
	"gobacktrader/asset"
	"gobacktrader/btutil"
	"gobacktrader/trade"
	"testing"
)

func TestFillAtLast(t *testing.T) {
	portfolio, err1 := asset.NewPortfolio("XXX", "AUD")
	stock, err2 := asset.NewStock("ZZB AU", "AUD")
	cash, err3 := asset.NewCash("AUD")
	if err := btutil.AnyValidError(err1, err2, err3); err != nil {
		t.Errorf("Error in asset init - %s", err)
	}

	portfolio.Transfer(cash, 1000.0)
	buyTrade := trade.NewTrade(portfolio, stock, +100.0)
	sellTrade := trade.NewTrade(portfolio, stock, -100.0)

	fillAtLast := FillAtLast{}
	fillAtLastWithSlippage := FillAtLastWithSlippage{slippage: 0.02}

	// we need a valid price to calculate trade consideration
	err1 = fillAtLast.Execute(buyTrade)
	err2 = fillAtLastWithSlippage.Execute(buyTrade)
	for _, err := range []error{err1, err2} {
		errStr := btutil.GetErrorString(err)
		if errStr != "'ZZB AU' cannot execute a trade with invalid consideration" {
			t.Errorf("Unexpected error string '%s'", errStr)
		}
	}

	stock.SetPrice(asset.Price{Float64: 2.50, Valid: true})
	// execute the buy trade
	err := fillAtLast.Execute(buyTrade)
	if err != nil {
		t.Errorf("Error in FillAtLast{}.Execute() - %s", err)
	}
	if portfolio.GetUnits(stock) != 100 {
		t.Error("Unexpected stock position")
	}
	if portfolio.GetUnits(cash) != 750 {
		t.Error("Unexpected cash position")
	}

	// and the sell trade with some slippage
	// trade is for $250 worth of stock
	// 2% slippage = 0.02 * 250 = $5
	err = fillAtLastWithSlippage.Execute(sellTrade)
	if err != nil {
		t.Errorf("Error in FillAtLastWithSlippage{}.Execute() - %s", err)
	}
	if portfolio.GetUnits(stock) != 0 { // position should be sold
		t.Error("Unexpected stock position")
	}
	if portfolio.GetUnits(cash) != 995 { // the orig. 1K we started with less $5 slippage
		t.Error("Unexpected cash position")
	}
}

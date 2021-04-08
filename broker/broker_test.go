package broker

import (
	"gobacktrader/asset"
	"gobacktrader/btutil"
	"gobacktrader/trade"
	"testing"
)

func TestBroker(t *testing.T) {
	portfolio, err1 := asset.NewPortfolio("XXX", "AUD")
	stock, err2 := asset.NewStock("ZZB AU", "AUD")
	cash, err3 := asset.NewCash("AUD")
	if err := btutil.AnyValidError(err1, err2, err3); err != nil {
		t.Errorf("Error in asset init - %s", err)
	}

	portfolio.Transfer(cash, 1000.0)
	stock.SetPrice(asset.Price{Float64: 2.50, Valid: true})
	trade := trade.NewTrade(portfolio, stock, +100.0)

	broker := NewBroker(
		NewNoCharges(),
		NewFillAtLastWithSlippage(0.01),
	)

	broker.Execute(trade)

	// we should now have 100 shares in the portfolio
	// cash should be 1000 - (2.50 * 1.01) * 100 = 747.50
	if portfolio.GetUnits(stock) != 100 {
		t.Error("Unexpected stock position")
	}
	if portfolio.GetUnits(cash) != 747.50 {
		t.Error("Unexpected AUD cash position")
	}
}

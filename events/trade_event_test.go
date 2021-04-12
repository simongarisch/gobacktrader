package events

import (
	"gobacktrader/asset"
	"gobacktrader/broker"
	"gobacktrader/btutil"
	"gobacktrader/trade"
	"testing"
	"time"
)

func TestTradeEventBasic(t *testing.T) {
	portfolio, err1 := asset.NewPortfolio("XXX", "AUD")
	stock, err2 := asset.NewStock("ZZB AU", "AUD")
	aud, err3 := asset.NewCash("AUD")
	if err := btutil.AnyValidError(err1, err2, err3); err != nil {
		t.Errorf("Error in asset init - %s", err)
	}

	executingBroker := broker.NewBroker(
		broker.NewNoCharges(),
		broker.NewFillAtLast(),
	)

	portfolio.SetBroker(executingBroker)
	portfolio.Transfer(aud, 1000)
	stock.SetPrice(asset.Price{Float64: 2.50, Valid: true})
	trade := trade.NewTrade(portfolio, stock, 100)
	tradeTime := time.Date(2020, time.December, 14, 0, 0, 0, 0, time.UTC)

	event := NewTradeEvent(trade, tradeTime)
	if !event.GetTime().Equal(tradeTime) {
		t.Error("Unexpected trade time")
	}
	if event.IsProcessed() {
		t.Error("Event should not yet be processed")
	}

	err := event.Process()
	if err != nil {
		t.Errorf("Error in event.Process() - %s", err)
	}

	if portfolio.GetUnits(stock) != 100 {
		t.Error("Expecting the portfolio to have 100 shares of stock")
	}
	if portfolio.GetUnits(aud) != 750 {
		t.Error("Expecting the portfolio to hold 750 AUD")
	}
}

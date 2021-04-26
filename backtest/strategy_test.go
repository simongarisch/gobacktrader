package backtest

import (
	"gobacktrader/asset"
	"gobacktrader/btutil"
	"gobacktrader/trade"
	"testing"
)

type myStrategy struct{}

func (s myStrategy) GenerateTrades() ([]*trade.Trade, error) {
	return []*trade.Trade{nil, nil}, nil
}

func TestStrategy(t *testing.T) {
	var strategy IStrategy

	strategy = myStrategy{}
	trades, _ := strategy.GenerateTrades()
	if len(trades) != 2 {
		t.Error("Expecting two trades to be returned")
	}
	if trades[0] != nil {
		t.Error("Unexpected first trade to be nil")
	}
	if trades[1] != nil {
		t.Error("Unexpected second trade to be nil")
	}
}

func TestStrategySwap(t *testing.T) {
	portfolio, err1 := asset.NewPortfolio("XXX", "AUD")
	stock, err2 := asset.NewStock("ZZB AU", "AUD")
	if err := btutil.AnyValidError(err1, err2); err != nil {
		t.Errorf("Error in asset init - %s", err)
	}

	buyTrade := trade.NewTrade(portfolio, stock, +100.0)
	sellTrade := trade.NewTrade(portfolio, stock, -100.0)

	// a strategy that generates a buy trade
	generateBuyTrade := func() ([]*trade.Trade, error) {
		return []*trade.Trade{buyTrade}, nil
	}

	// a strategy that generates a sell trade
	generateSellTrade := func() ([]*trade.Trade, error) {
		return []*trade.Trade{sellTrade}, nil
	}

	strategy := NewStrategy(generateBuyTrade)
	trades, _ := strategy.GenerateTrades()
	if len(trades) != 1 {
		t.Fatal("Expecting one trade to be generated")
	}
	if trades[0] != buyTrade {
		t.Error("Unexpected trade generated")
	}

	// swap out this strategy
	strategy.SetGenerateTrades(generateSellTrade)
	trades, _ = strategy.GenerateTrades()
	if len(trades) != 1 {
		t.Fatal("Expecting one trade to be generated")
	}
	if trades[0] != sellTrade {
		t.Error("Unexpected trade generated")
	}
}

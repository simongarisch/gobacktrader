package backtest

import (
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

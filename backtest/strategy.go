package backtest

import "gobacktrader/trade"

// Strategy defines the interface for our trading strategy.
// This must have a GenerateTrades method that returns a slice
// of pointer to Trade.
type Strategy interface {
	GenerateTrades() []*trade.Trade
}

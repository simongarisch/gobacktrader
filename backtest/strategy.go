package backtest

import "gobacktrader/trade"

// IStrategy defines the interface for our trading strategy.
// This must have a GenerateTrades method that returns a slice
// of pointer to Trade.
type IStrategy interface {
	GenerateTrades() ([]*trade.Trade, error)
}

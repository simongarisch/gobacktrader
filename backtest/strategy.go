package backtest

import "gobacktrader/trade"

// IStrategy defines the interface for our trading strategy.
// This must have a GenerateTrades method that returns a slice
// of pointer to Trade.
type IStrategy interface {
	GenerateTrades() ([]*trade.Trade, error)
}

type generateTradesFunc func() ([]*trade.Trade, error)

// Strategy has a function field to generate trades.
type Strategy struct {
	generateTradesFunc
}

// NewStrategy returns a new strategy instance.
func NewStrategy(f generateTradesFunc) *Strategy {
	return &Strategy{generateTradesFunc: f}
}

// SetGenerateTrades sets our generate trades strategy function.
func (s *Strategy) SetGenerateTrades(f generateTradesFunc) {
	s.generateTradesFunc = f
}

// GenerateTrades returns a slice of trades to execute.
func (s *Strategy) GenerateTrades() ([]*trade.Trade, error) {
	return s.generateTradesFunc()
}

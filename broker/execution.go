package broker

import "gobacktrader/trade"

// FillAtLast executes a trade at the last available price.
type FillAtLast struct{}

// Execute executes a specific trade.
func (e FillAtLast) Execute(trade trade.Trade) error {
	return nil
}

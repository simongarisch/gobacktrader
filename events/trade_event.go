package events

import (
	"gobacktrader/trade"
	"time"
)

// TradeEvent defines a trade at a specific time.
type TradeEvent struct {
	BaseEvent
	trade *trade.Trade
}

// NewTradeEvent returns a new instance of TradeEvent.
func NewTradeEvent(trade *trade.Trade, eventTime time.Time) TradeEvent {
	return TradeEvent{
		BaseEvent: BaseEvent{eventTime: eventTime, processed: false},
		trade:     trade,
	}
}

// GetTrade returns the trade for this event.
func (e *TradeEvent) GetTrade() *trade.Trade {
	return e.trade
}

// Process with execute the trade if it passes compliance.
func (e *TradeEvent) Process() error {
	_, err := e.trade.Execute()
	if err != nil {
		return err
	}
	e.processed = true
	return nil
}

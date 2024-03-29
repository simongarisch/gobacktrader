package broker

import (
	"gobacktrader/asset"
	"gobacktrader/btutil"
)

// ChargesStrategy defines the interface for broker charges.
type ChargesStrategy interface {
	Charge(asset.ITrade) error
}

// ExecutionStrategy defines the interface for broker execution.
type ExecutionStrategy interface {
	Execute(asset.ITrade) error
}

// Broker defines an executing broker with associated charges.
type Broker struct {
	charges   ChargesStrategy
	execution ExecutionStrategy
}

// NewBroker returns a new Broker instance.
func NewBroker(charges ChargesStrategy, execution ExecutionStrategy) *Broker {
	return &Broker{
		charges:   charges,
		execution: execution,
	}
}

// Execute will use our broker instance to execute a trade.
func (b *Broker) Execute(trade asset.ITrade) error {
	err1 := b.execution.Execute(trade)
	err2 := b.charges.Charge(trade)
	return btutil.AnyValidError(err1, err2)
}

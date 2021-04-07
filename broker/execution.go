package broker

import (
	"fmt"
	"gobacktrader/btutil"
	"gobacktrader/trade"
)

// FillAtLast executes a trade at the last available price.
type FillAtLast struct{}

// Execute executes a specific trade.
func (e FillAtLast) Execute(trade trade.Trade) error {
	consideration := trade.GetLocalCurrencyConsideration()
	if !consideration.Valid {
		tradeTicker := trade.GetAsset().GetTicker()
		return fmt.Errorf("'%s' cannot execute a trade with invalid consideration", tradeTicker)
	}

	portfolio, asset, units := trade.GetPortfolio(), trade.GetAsset(), trade.GetUnits()
	portfolio.Trade(asset, units, &consideration.Float64)
	return nil
}

// FillAtLastWithSlippage executes a trade with some perecentage
// slippage to the last available price.
type FillAtLastWithSlippage struct {
	slippage float64
}

// Execute executes a specific trade with slippage.
func (e FillAtLastWithSlippage) Execute(trade trade.Trade) error {
	consideration := trade.GetLocalCurrencyConsideration()
	if !consideration.Valid {
		tradeTicker := trade.GetAsset().GetTicker()
		return fmt.Errorf("'%s' cannot execute a trade with invalid consideration", tradeTicker)
	}

	portfolio, asset, units := trade.GetPortfolio(), trade.GetAsset(), trade.GetUnits()
	// slippage on buy trades will result in more consideration being paid
	// slippage on sell trades will result in less consideration being received
	considerationFloat := consideration.Float64
	if btutil.Sgn(units) == +1.0 { // we are buyers and will have to pay more
		considerationFloat *= (1 + e.slippage)
	}
	if btutil.Sgn(units) == -1.0 { // we are sellers and will receive less
		considerationFloat *= (1 - e.slippage)
	}
	portfolio.Trade(asset, units, &considerationFloat)
	return nil
}

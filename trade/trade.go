package trade

import (
	"errors"
	"gobacktrader/asset"
	"gobacktrader/btutil"
	"math"
)

// Trade defines a trade in some asset for a given portfolio.
type Trade struct {
	portfolio   *asset.Portfolio
	targetAsset asset.IAssetReadOnly
	units       float64
}

// NewTrade returns a new Trade instance.
func NewTrade(portfolio *asset.Portfolio, targetAsset asset.IAssetReadOnly, units float64) *Trade {
	return &Trade{
		portfolio:   portfolio,
		targetAsset: targetAsset,
		units:       units,
	}
}

// GetPortfolio returns the target portfolio.
func (t *Trade) GetPortfolio() *asset.Portfolio {
	return t.portfolio
}

// GetAsset returns the target asset.
func (t *Trade) GetAsset() asset.IAssetReadOnly {
	return t.targetAsset
}

// GetUnits returns the units to be traded.
func (t *Trade) GetUnits() float64 {
	return t.units
}

// GetBaseCurrencyCash returns the cash base currency for this trade.
func (t *Trade) GetBaseCurrencyCash() (asset.IAssetReadOnly, error) {
	return asset.NewCash(t.targetAsset.GetBaseCurrency())
}

// GetLocalCurrencyValue returns the trade value.
func (t *Trade) GetLocalCurrencyValue() asset.Price {
	assetValue := t.targetAsset.GetValue()
	if !assetValue.Valid {
		return asset.Price{Float64: 0.0, Valid: false}
	}

	tradeValue := assetValue.Float64 * math.Abs(t.units)
	return asset.Price{Float64: tradeValue, Valid: true}
}

// GetLocalCurrencyConsideration returns the cash transferred
// for this trade in local currency.
func (t *Trade) GetLocalCurrencyConsideration() asset.Price {
	tradeValue := t.GetLocalCurrencyValue()
	if !tradeValue.Valid {
		return asset.Price{Float64: 0.0, Valid: false}
	}

	// cash goes out for buys and comes in for sells
	consideration := tradeValue.Float64 * btutil.Sgn(t.GetUnits()) * -1.0

	return asset.Price{Float64: consideration, Valid: true}
}

// PassesCompliance returns true if compliance passes, false otherwise.
func (t *Trade) PassesCompliance() (bool, error) {
	portfolio := t.GetPortfolio()
	if portfolio.NumComplianceRules() == 0 {
		return true, nil // no rules are in place
	}

	portfolioCopy, err := portfolio.Copy()
	if err != nil {
		return false, err
	}

	if portfolioCopy.GetBroker() == nil {
		return false, errors.New("portfolio has no assigned executing broker")
	}

	// use this portfolio copy to mock execute the trade
	// and check whether compliance passes after execution.
	err1 := portfolioCopy.GetBroker().Execute(t)
	passes, err2 := portfolioCopy.PassesCompliance()
	if err := btutil.AnyValidError(err1, err2); err != nil {
		return false, err
	}

	return passes, nil
}

// Execute will execute this trade if compliance passes.
// Returns a boolean representing whether the trade has
// been executed along with an error.
func (t *Trade) Execute() (bool, error) {
	// check if compliance passes
	passes, err := t.PassesCompliance()
	if err != nil {
		return false, err
	}

	if !passes {
		return false, nil // do not execute this trade
	}

	// then pass to the broker for execution
	executingBroker := t.GetPortfolio().GetBroker()
	if executingBroker == nil {
		return false, errors.New("portfolio requires an executing broker to call trade.Execute()")
	}

	err = executingBroker.Execute(t)
	if err != nil {
		return false, err
	}

	return true, nil
}

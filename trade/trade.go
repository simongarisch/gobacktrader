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
func NewTrade(portfolio *asset.Portfolio, targetAsset asset.IAssetReadOnly, units float64) Trade {
	return Trade{
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
func (t Trade) PassesCompliance() (bool, error) {
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

	return true, nil
}

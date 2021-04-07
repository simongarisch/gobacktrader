package trade

import (
	"gobacktrader/asset"
)

// Trade defines a trade in some asset for a given portfolio.
type Trade struct {
	portfolio   *asset.Portfolio
	targetAsset asset.IAssetReadOnly
	units       float64
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

	tradeValue := assetValue.Float64 * t.units
	return asset.Price{Float64: tradeValue, Valid: true}
}

// PassesCompliance returns true if compliance passes, false otherwise.
func (t Trade) PassesCompliance() (bool, error) {
	return true, nil
}

package compliance

import "gobacktrader/asset"

// UnitLimit places unit limits on our portfolio holdings.
type UnitLimit struct {
	portfolio   *asset.Portfolio
	targetAsset asset.IAssetReadOnly
	limit       float64
}

// NewUnitLimit returns a new instance of UnitLimit.
func NewUnitLimit(portfolio *asset.Portfolio, targetAsset asset.IAssetReadOnly, limit float64) UnitLimit {
	return UnitLimit{
		portfolio:   portfolio,
		targetAsset: targetAsset,
		limit:       limit,
	}
}

// GetPortfolio returns the portfolio for which this limit is applied.
func (r *UnitLimit) GetPortfolio() *asset.Portfolio {
	return r.portfolio
}

// GetAsset returns the asset for which this limit is applied.
func (r *UnitLimit) GetAsset() asset.IAssetReadOnly {
	return r.targetAsset
}

// GetLimit returns the unit limit.
func (r *UnitLimit) GetLimit() float64 {
	return r.limit
}

// Passes returns true if the UnitLimit rule passes, false otherwise.
func (r *UnitLimit) Passes() (bool, error) {
	units := r.portfolio.GetUnits(r.targetAsset)
	if units > r.limit {
		return false, nil
	}
	return true, nil
}

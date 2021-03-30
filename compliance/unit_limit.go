package compliance

import "gobacktrader/asset"

// UnitLimit places unit limits on our portfolio holdings.
type UnitLimit struct {
	targetAsset asset.IAssetReadOnly
	limit       float64
}

// NewUnitLimit returns a new instance of UnitLimit.
func NewUnitLimit(targetAsset asset.IAssetReadOnly, limit float64) UnitLimit {
	unitLimit := UnitLimit{
		targetAsset: targetAsset,
		limit:       limit,
	}
	return unitLimit
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
func (r *UnitLimit) Passes(portfolio *asset.Portfolio) (bool, error) {
	units := portfolio.GetUnits(r.targetAsset)
	if units > r.limit {
		return false, nil
	}
	return true, nil
}

package compliance

import (
	"gobacktrader/asset"
	"math"
)

// WeightLimit places weight limits on our portfolio holdings.
type WeightLimit struct {
	targetAsset asset.IAssetReadOnly
	limit       float64
}

// NewWeightLimit returns a new instance of WeightLimit
func NewWeightLimit(targetAsset asset.IAssetReadOnly, limit float64) WeightLimit {
	weightLimit := WeightLimit{
		targetAsset: targetAsset,
		limit:       limit,
	}
	return weightLimit
}

// GetAsset returns the asset for which this limit is applied.
func (r *WeightLimit) GetAsset() asset.IAssetReadOnly {
	return r.targetAsset
}

// GetLimit returns the unit limit.
func (r *WeightLimit) GetLimit() float64 {
	return r.limit
}

// Passes returns true if the WeightLimit rule passes, false otherwise.
func (r *WeightLimit) Passes(portfolio *asset.Portfolio) (bool, error) {
	weight, err := portfolio.GetWeight(r.targetAsset)
	if err != nil {
		return false, err
	}

	if !weight.Valid { // prices or fx rates not initialised
		return false, nil
	}

	if math.Abs(weight.Float64) > math.Abs(r.limit) {
		return false, nil
	}
	return true, nil
}

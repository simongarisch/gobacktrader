package asset

import (
	"database/sql"
)

// Weight represents an asset weight for a given portfolio.
type Weight sql.NullFloat64

var nullWeight = Weight{Float64: 0.0, Valid: false}

// Portfolio consists of a collection of positions.
type Portfolio struct {
	code         string
	baseCurrency string
	positions    map[IAssetReadOnly]*Position
	fxRates      *FxRates
}

// NewPortfolio returns a new instance of Portfolio.
func NewPortfolio(code string, baseCurrency string) Portfolio {
	positions := make(map[IAssetReadOnly]*Position)
	return Portfolio{
		code:         code,
		baseCurrency: baseCurrency,
		positions:    positions,
	}
}

// GetCode returns our portfolio code
func (p *Portfolio) GetCode() string {
	return p.code
}

// GetBaseCurrency returns the base currency for our portfolio.
func (p *Portfolio) GetBaseCurrency() string {
	return p.baseCurrency
}

// NumPositions returns the number of portfolio positions.
func (p *Portfolio) NumPositions() int {
	return len(p.positions)
}

// HasAsset returns a boolean of true if a portfolio
// contains some asset, false otherwise.
func (p *Portfolio) HasAsset(a IAssetReadOnly) bool {
	_, ok := p.positions[a]
	return ok
}

// GetUnits returns the units held for a given asset.
func (p *Portfolio) GetUnits(a IAssetReadOnly) float64 {
	position, ok := p.positions[a]
	if !ok {
		return 0.0
	}
	return position.GetUnits()
}

// GetWeight returns the portfolio weight in a given asset.
func (p *Portfolio) GetWeight(a IAssetReadOnly) (Weight, error) {
	_, positionWeights, err := p.GetValueWeights()
	if err != nil {
		return nullWeight, err
	}

	weight, ok := positionWeights[a]
	if !ok {
		return nullWeight, nil
	}

	return weight, nil
}

// ModifyPositions allows us to increment and decrement positions
// in the portfolio.
func (p *Portfolio) ModifyPositions(a IAssetReadOnly, units float64) {
	// if the asset is already held then modify its position
	if p.HasAsset(a) {
		p.positions[a].Increment(units)
	}
	// otherwise create a new position for this asset
	p.positions[a] = &Position{a, units}
}

// GetValue returns our portfolio value.
func (p Portfolio) GetValue() (Price, error) {
	portfolioValue, _, err := p.GetValueWeights()
	return portfolioValue, err
}

// GetValueWeights returns the portfolio value along with all position weights.
func (p *Portfolio) GetValueWeights() (Price, map[IAssetReadOnly]Weight, error) {
	var totalValueFloat float64
	portfolioValue := nullPrice
	valid := true
	positionValues := make(map[IAssetReadOnly]Price)
	positionWeights := make(map[IAssetReadOnly]Weight)

	for _, position := range p.positions {
		asset := position.GetAsset()
		value := position.GetValue()
		if !value.Valid { // no value, so value and weight invalid
			valid = false
			positionValues[asset] = nullPrice
			positionWeights[asset] = nullWeight
			continue
		}

		assetBaseCurrency := position.GetBaseCurrency()
		pair := assetBaseCurrency + p.baseCurrency
		fxRate, ok, err := p.fxRates.GetRate(pair)
		if err != nil {
			return portfolioValue, positionWeights, err
		}
		if !ok { // no fx rate, so value and weight invalid
			valid = false
			positionValues[asset] = nullPrice
			positionWeights[asset] = nullWeight
			continue
		}

		positionValueBaseCurrency := value.Float64 * fxRate
		positionValues[asset] = Price{Float64: positionValueBaseCurrency, Valid: true}
		totalValueFloat += positionValueBaseCurrency
	}

	if valid { // all assets have a valid value which we can use to derive portfolio value
		portfolioValue = Price{Float64: totalValueFloat, Valid: true}
		for asset, price := range positionValues {
			positionWeights[asset] = Weight{Float64: price.Float64 / totalValueFloat, Valid: true}
		}
	}

	return portfolioValue, positionWeights, nil
}

// SetFxRates sets the FX rates object to be used for this portfolio.
func (p *Portfolio) SetFxRates(fxRates *FxRates) {
	p.fxRates = fxRates
}

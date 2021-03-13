package asset

import (
	"database/sql"
	"fmt"
)

// Weight represents an asset weight for a given portfolio.
type Weight sql.NullFloat64

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

func (p *Portfolio) GetWeight(a IAssetReadOnly) (Weight, error) {
	invalidWeight := Weight{Float64: 0.0, Valid: false}
	position, ok := p.positions[a]
	if !ok {
		return Weight{Float64: 0.0, Valid: true}, nil
	}

	// we need a valid position value to continue
	positionValue := position.GetValue()
	if !positionValue.Valid {
		return invalidWeight, nil
	}

	// value this position in the portfolio base currency
	pair := position.GetBaseCurrency() + p.GetBaseCurrency()
	fxRate, ok, err := p.fxRates.GetRate(pair)
	if err != nil {
		return invalidWeight, err
	}
	if !ok {
		return invalidWeight, nil // no valid FX rate initialised.
	}

	// collect the portfolio value
	positionBaseCurrencyValue := positionValue.Float64 * fxRate
	portfolioValue, err := p.GetValue()
	if err != nil {
		return invalidWeight, err
	}
	if !portfolioValue.Valid {
		return invalidWeight, nil
	}

	if portfolioValue.Float64 == 0.0 {
		return invalidWeight, fmt.Errorf("'%s' unable to collect weights with a zero portfolio value", p.GetCode())
	}

	weight := positionBaseCurrencyValue / portfolioValue.Float64
	return Weight{Float64: weight, Valid: true}, nil
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
	totalValue := 0.0
	valid := true
	for _, position := range p.positions {
		value := position.GetValue()
		if !value.Valid {
			valid = false
			break
		}

		assetBaseCurrency := position.GetBaseCurrency()
		pair := assetBaseCurrency + p.baseCurrency
		fxRate, ok, err := p.fxRates.GetRate(pair)
		if err != nil {
			return Price{}, err
		}
		if !ok {
			valid = false
			break
		}

		totalValue += value.Float64 * fxRate
	}

	return Price{Float64: totalValue, Valid: valid}, nil
}

// SetFxRates sets the FX rates object to be used for this portfolio.
func (p *Portfolio) SetFxRates(fxRates *FxRates) {
	p.fxRates = fxRates
}

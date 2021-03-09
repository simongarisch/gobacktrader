package asset

// Position represents a holding in some asset.
type Position struct {
	asset IAssetReadOnly
	units float64
}

// NewPosition creates a new asset position.
func NewPosition(asset IAssetReadOnly, units float64) Position {
	return Position{asset: asset, units: units}
}

// Increment increments our position units.
func (p *Position) Increment(units float64) {
	p.units += units
}

// Decrement decrements our position units.
func (p *Position) Decrement(units float64) {
	p.units -= units
}

// GetUnits returns the position units.
func (p Position) GetUnits() float64 {
	return p.units
}

// GetAsset returns the position asset.
func (p *Position) GetAsset() IAssetReadOnly {
	return p.asset
}

// GetValue returns the position value.
func (p Position) GetValue() Price {
	assetValue := p.asset.GetValue()
	if !assetValue.Valid {
		return nullPrice
	}

	return Price{Float64: assetValue.Float64 * p.units, Valid: true}
}

// GetTicker returns the position ticker.
func (p Position) GetTicker() string {
	return p.asset.GetTicker()
}

// GetBaseCurrency returns the position base currency.
func (p Position) GetBaseCurrency() string {
	return p.asset.GetBaseCurrency()
}

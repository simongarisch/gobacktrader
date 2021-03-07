package asset

var unitPrice = Price{Float64: 1.0, Valid: true}

// Cash represents a cash asset.
type Cash struct {
	Asset
}

// NewCash returns a new cash asset.
func NewCash(currency string) error {
	return nil
}

// GetPrice returns the price for our cash asset.
func (c Cash) GetPrice() Price {
	return unitPrice
}

// GetValue returns the value for our cash asset.
func (c Cash) GetValue() Price {
	return unitPrice
}

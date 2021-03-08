// Package asset focuses on asset valuation and portfolio composition.
package asset

import (
	"database/sql"
	"gobacktrader/btutil"
)

// Price is the unit of measurement for asset price and value.
type Price sql.NullFloat64

var (
	nullPrice         = Price{Float64: 0.0, Valid: false}
	defaultMultiplier = 1.0
)

// Asset defines a generic asset type with a ticker, price and value.
type Asset struct {
	ticker       string
	baseCurrency string
	multiplier   float64
	price        Price
	value        Price
}

// IAsset defines an asset interface.
type IAsset interface {
	GetTicker() string
	GetBaseCurrency() string
	GetMultiplier() float64
	GetPrice() Price
	SetPrice(Price)
	GetValue() Price
	Revalue()
}

// NewAsset creates a new asset instance with a
// default multiplier.
func NewAsset(ticker string, baseCurrency string) Asset {
	ticker = btutil.CleanString(ticker)
	return Asset{
		ticker:       ticker,
		baseCurrency: baseCurrency,
		multiplier:   defaultMultiplier,
	}
}

// NewAssetWithMultiplier create a new asset with
// a non-default multiplier.
func NewAssetWithMultiplier(ticker string, baseCurrency string, multiplier float64) Asset {
	ticker = btutil.CleanString(ticker)
	return Asset{
		ticker:       ticker,
		baseCurrency: baseCurrency,
		multiplier:   multiplier,
	}
}

// GetTicker returns the asset's ticker code.
func (a Asset) GetTicker() string {
	return a.ticker
}

// GetBaseCurrency returns the asset's base currency code.
func (a Asset) GetBaseCurrency() string {
	return a.baseCurrency
}

// GetMultiplier returns the asset's multiplier
func (a Asset) GetMultiplier() float64 {
	return a.multiplier
}

// GetPrice returns the asset's price.
func (a Asset) GetPrice() Price {
	return a.price
}

// SetPrice sets the asset's price.
// The Revalue method is automatically called after setting price.
func (a *Asset) SetPrice(price Price) {
	a.price = price
	a.Revalue()
}

// GetValue returns the asset's value.
func (a Asset) GetValue() Price {
	return a.value
}

// Revalue revalues our asset.
func (a *Asset) Revalue() {
	if !a.price.Valid {
		a.value = nullPrice
		return
	}

	priceFloat := a.price.Float64
	multiplier := a.GetMultiplier()
	a.value = Price{Float64: priceFloat * multiplier, Valid: true}
}

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

// IAssetReadOnly defines the interface for read only assets.
// Given these interface methods take a pointer receiver only
// pointers to asset can satisfy this interface.
type IAssetReadOnly interface {
	GetTicker() string
	GetBaseCurrency() string
	GetValue() Price
}

// IAssetWriteOnly defines the interface for write only assets.
// SetPrice takes a pointer receiver so only pointers to asset
// can satisfy this interface.
type IAssetWriteOnly interface {
	SetPrice(price Price)
}

// NewAsset creates a new asset instance with a
// default multiplier.
func NewAsset(ticker string, baseCurrency string) (Asset, error) {
	ticker = btutil.CleanString(ticker)
	baseCurrency, err := ValidateCurrency(baseCurrency)
	asset := Asset{
		ticker:       ticker,
		baseCurrency: baseCurrency,
		multiplier:   defaultMultiplier,
	}
	return asset, err
}

// NewAssetWithMultiplier create a new asset with
// a non-default multiplier.
func NewAssetWithMultiplier(ticker string, baseCurrency string, multiplier float64) (Asset, error) {
	ticker = btutil.CleanString(ticker)
	baseCurrency, err := ValidateCurrency(baseCurrency)
	asset := Asset{
		ticker:       ticker,
		baseCurrency: baseCurrency,
		multiplier:   multiplier,
	}
	return asset, err
}

// GetTicker returns the asset's ticker code.
func (a *Asset) GetTicker() string {
	return a.ticker
}

// GetBaseCurrency returns the asset's base currency code.
func (a *Asset) GetBaseCurrency() string {
	return a.baseCurrency
}

// GetMultiplier returns the asset's multiplier
func (a *Asset) GetMultiplier() float64 {
	return a.multiplier
}

// GetPrice returns the asset's price.
func (a *Asset) GetPrice() Price {
	return a.price
}

// SetPrice sets the asset's price.
// The Revalue method is automatically called after setting price.
func (a *Asset) SetPrice(price Price) {
	a.price = price
	a.Revalue()
}

// GetValue returns the asset's value.
func (a *Asset) GetValue() Price {
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

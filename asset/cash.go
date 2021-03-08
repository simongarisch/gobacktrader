package asset

import (
	"fmt"
	"gobacktrader/btutil"
)

var unitPrice = Price{Float64: 1.0, Valid: true}

// Cash represents a cash asset.
type Cash struct {
	currency string
}

// NewCash returns a new cash asset.
func NewCash(currency string) (Cash, error) {
	currency = btutil.CleanString(currency)
	cash := Cash{currency: currency}
	if len(currency) != 3 {
		return cash, fmt.Errorf("'%s' is an invalid currency code", currency)
	}

	return cash, nil
}

// GetCurrency returns the cash currency.
func (c Cash) GetCurrency() string {
	return c.currency
}

// GetTicker returns the cash ticker.
func (c Cash) GetTicker() string {
	return c.currency
}

// GetPrice returns the price for our cash asset.
func (c Cash) GetPrice() Price {
	return unitPrice
}

// GetValue returns the value for our cash asset.
func (c Cash) GetValue() Price {
	return unitPrice
}

// GetBaseCurrency returns the base currency code.
func (c Cash) GetBaseCurrency() string {
	return c.currency
}

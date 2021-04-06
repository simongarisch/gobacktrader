package asset

// Cash represents a cash asset.
type Cash struct {
	priceHistory
	currency string
}

type cashObjects struct {
	pool map[string]*Cash
}

var cashObjectsSingleton = cashObjects{}

func (c *cashObjects) getCash(currency string) (*Cash, error) {
	if c.pool == nil {
		c.pool = make(map[string]*Cash)
	}

	_, ok := c.pool[currency]
	if !ok {
		cash, err := cashFactory(currency)
		if err != nil {
			return &cash, err
		}
		c.pool[currency] = &cash
	}

	return c.pool[currency], nil
}

func cashFactory(currency string) (Cash, error) {
	currency, err := ValidateCurrency(currency)
	cash := Cash{currency: currency}
	cash.price = unitPrice
	return cash, err
}

// NewCash returns a new cash asset.
func NewCash(currency string) (*Cash, error) {
	return cashObjectsSingleton.getCash(currency)
}

// GetCurrency returns the cash currency.
func (c Cash) GetCurrency() string {
	return c.currency
}

// GetTicker returns the cash ticker.
func (c Cash) GetTicker() string {
	return c.currency
}

// GetValue returns the value for our cash asset.
func (c Cash) GetValue() Price {
	return unitPrice
}

// GetBaseCurrency returns the base currency code.
func (c Cash) GetBaseCurrency() string {
	return c.currency
}

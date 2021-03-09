package asset

// NewStock returns a new stock asset.
func NewStock(ticker string, baseCurrency string) (Asset, error) {
	return NewAssetWithMultiplier(ticker, baseCurrency, 1.0)
}

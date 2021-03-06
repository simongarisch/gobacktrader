package asset

// NewStock returns a new stock asset.
func NewStock(ticker string) Asset {
	return NewAssetWithMultiplier(ticker, 1.0)
}

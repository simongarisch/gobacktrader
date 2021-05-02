package datasources

import (
	"gobacktrader/asset"
	"gobacktrader/events"
	"time"
)

// IAssetPriceQuery defines the asset price query interface.
// These must have the ability to generate asset price events.
type IAssetPriceQuery interface {
	GetURL() string
	GenerateEvents() ([]events.AssetPriceEvent, error)
}

// Query is the base struct for data queries
type Query struct {
	targetAsset asset.IAssetReadOnly
	startDate   time.Time
	endDate     time.Time
	apiKey      string
	ticker      string
}

// NewQuery returns a new instance of Query.
func NewQuery(targetAsset asset.IAssetReadOnly, startDate time.Time, endDate time.Time) Query {
	return Query{
		targetAsset: targetAsset,
		startDate:   startDate,
		endDate:     endDate,
		apiKey:      "demo",
	}
}

// GetAsset returns the query asset instance.
func (q Query) GetAsset() asset.IAssetReadOnly {
	return q.targetAsset
}

// GetStartDate returns the query start date.
func (q Query) GetStartDate() time.Time {
	return q.startDate
}

// GetEndDate returns the query end date.
func (q Query) GetEndDate() time.Time {
	return q.endDate
}

// GetAPIKey gets the query api key.
func (q Query) GetAPIKey() string {
	return q.apiKey
}

// SetAPIKey sets the query api key.
func (q *Query) SetAPIKey(apiKey string) *Query {
	q.apiKey = apiKey
	return q
}

// GetTicker returns the ticker used for our query.
func (q Query) GetTicker() string {
	ticker := q.ticker
	if ticker != "" {
		return ticker
	}
	return q.targetAsset.GetTicker()
}

// SetTicker sets the ticker used for our query and returns the query instance.
// The ticker used for a specific query may differ from the asset ticker,
// so this provides a work around.
func (q *Query) SetTicker(ticker string) *Query {
	q.ticker = ticker
	return q
}

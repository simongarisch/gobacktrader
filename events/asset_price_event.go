package events

import (
	"time"
	"gobacktrader/asset"
)

// AssetPriceEvent defines price events for generic assets.
type AssetPriceEvent struct {
	targetAsset asset.IAssetWriteOnly
	eventTime time.Time
	price asset.Price
	processed bool
}

// NewAssetPriceEvent returns a new instance of an unprocessed
// AssetPriceEvent object.
func NewAssetPriceEvent(targetAsset asset.IAssetWriteOnly, eventTime time.Time, price asset.Price) AssetPriceEvent {
	return AssetPriceEvent{
		targetAsset: targetAsset,
		eventTime: eventTime,
		price: price,
		processed: false,
	}
}

// GetTime returns the event time.
func (e AssetPriceEvent) GetTime() time.Time {
	return e.eventTime
}

// IsProcessed returns true if an event has been processed, false otherwise.
func (e AssetPriceEvent) IsProcessed() bool {
	return e.processed
}

// Process will action the asset price event.
func (e *AssetPriceEvent) Process() {
	e.targetAsset.SetPrice(e.price)
	e.processed = true
}

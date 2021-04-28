package events

import (
	"gobacktrader/asset"
	"time"
)

// AssetPriceEvent defines price events for generic assets.
type AssetPriceEvent struct {
	BaseEvent
	targetAsset asset.IAssetWriteOnly
	price       asset.Price
}

// NewAssetPriceEvent returns a new instance of an unprocessed
// AssetPriceEvent object.
func NewAssetPriceEvent(targetAsset asset.IAssetWriteOnly, eventTime time.Time, price asset.Price) AssetPriceEvent {
	return AssetPriceEvent{
		BaseEvent:   BaseEvent{eventTime: eventTime, processed: false},
		targetAsset: targetAsset,
		price:       price,
	}
}

// Process will action the asset price event.
func (e *AssetPriceEvent) Process() error {
	e.targetAsset.SetPrice(e.price)
	e.processed = true
	return nil
}

// GetPrice returns the event price.
func (e AssetPriceEvent) GetPrice() asset.Price {
	return e.price
}

package events

import (
	"gobacktrader/asset"
	"time"
)

// FxRateEvent attaches a time to an FX rate change.
type FxRateEvent struct {
	AssetPriceEvent
}

// NewFxRateEvent returns a new instance of FxRateEvent.
func NewFxRateEvent(rate *asset.FxRate, eventTime time.Time, price asset.Price) AssetPriceEvent {
	return NewAssetPriceEvent(rate, eventTime, price)
}

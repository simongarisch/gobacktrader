package events

import (
	"time"
	"gobacktrader/asset"
)

type FxRateEvent struct {
	AssetPriceEvent
}

func NewFxRateEvent(rate *asset.FxRate, eventTime time.Time, price asset.Price) AssetPriceEvent {
	return NewAssetPriceEvent(rate, eventTime, price)
}

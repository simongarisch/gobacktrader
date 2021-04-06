package events

import (
	"gobacktrader/asset"
	"testing"
	"time"
)

func TestAssetPriceEvent(t *testing.T) {
	stock, err := asset.NewStock("ZZB AU", "AUD")
	if err != nil {
		t.Errorf("Error in asset.NewStock - %s", err)
	}
	eventTime := time.Date(2021, time.March, 13, 0, 0, 0, 0, time.UTC)
	price := asset.Price{Float64: 3.00, Valid: true}

	// create a new asset price event
	assetPriceEvent := NewAssetPriceEvent(stock, eventTime, price)
	if !assetPriceEvent.GetTime().Equal(eventTime) {
		t.Error("Unexpected event time.")
	}
	if assetPriceEvent.IsProcessed() == true {
		t.Error("This event should be new and unprocessed.")
	}

	// process this event
	assetPriceEvent.Process()
	value := stock.GetValue()
	if !value.Valid {
		t.Error("Expecting a valid stock value.")
	}
	if value.Float64 != 3.0 {
		t.Errorf("Expecting a value of $3, got %0.2f", value.Float64)
	}
	if assetPriceEvent.IsProcessed() == false {
		t.Error("Expecting this event to be processed.")
	}
}

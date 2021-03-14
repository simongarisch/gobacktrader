package events

import (
	"testing"
	"time"
	"gobacktrader/asset"
)

func TestFxRateEvent(t *testing.T) {
	audusd, err := asset.NewFxRate("AUDUSD", asset.Price{Float64: 0.75, Valid: true})
	if err != nil {
		t.Errorf("Error in NewFxRate - %s", err)
	}

	eventTime := time.Date(2021, time.March, 13, 0, 0, 0, 0, time.UTC)
	newRate := asset.Price{Float64: 0.80, Valid: true}

	// create a new FX rate event
	fxRateEvent := NewFxRateEvent(&audusd, eventTime, newRate)
	if !fxRateEvent.GetTime().Equal(eventTime) {
		t.Error("Unexpected event time.")
	}
	if fxRateEvent.IsProcessed() == true {
		t.Error("This event should be new and unprocessed.")
	}

	// process this event
	fxRateEvent.Process()
	rate := audusd.GetRate()
	if !rate.Valid {
		t.Error("Expecting a valid rate.")
	}
	if rate.Float64 != 0.8 {
		t.Errorf("Expecting a rate of 0.8, got %0.2f", rate.Float64)
	}
}
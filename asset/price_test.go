package asset

import (
	"testing"
	"time"
)

func TestPriceSnapshot(t *testing.T) {
	stock, err := NewStock("ZZB AU", "AUD")
	if err != nil {
		t.Errorf("Error in NewStock - %s", err)
	}

	// set the stock price
	price := Price{Float64: 2.5, Valid: true}
	stock.SetPrice(price)

	// create our new snapshot and test
	timestamp := time.Date(2020, time.December, 14, 0, 0, 0, 0, time.UTC)
	snap := newPriceSnapshot(timestamp, stock)

	if !snap.GetTime().Equal(timestamp) {
		t.Error("Unexpected timestamp")
	}
	if snap.GetPrice().Float64 != 2.50 {
		t.Error("Unexpected price")
	}
}

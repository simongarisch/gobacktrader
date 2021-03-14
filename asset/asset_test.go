package asset

import (
	"testing"
	"time"
)

func TestNewAsset(t *testing.T) {
	asset, err := NewAsset("zzb au", "AUD")
	if err != nil {
		t.Errorf("Error in NewAsset - %s", err)
	}
	if asset.GetTicker() != "ZZB AU" {
		t.Error("Unexpected ticker")
	}
	if asset.GetBaseCurrency() != "AUD" {
		t.Error("Unexpected base currency")
	}
	if asset.GetPrice().Valid {
		t.Error("Expecting an uninitialised price")
	}

	multiplier := asset.GetMultiplier()
	if multiplier != defaultMultiplier {
		t.Errorf("Unexpected multiplier: wanted %0.2f, got %0.2f", defaultMultiplier, multiplier)
	}
}

func TestNewAssetWithMultiplier(t *testing.T) {
	asset, err := NewAssetWithMultiplier("ZZB AU", "AUD", 100.0)
	if err != nil {
		t.Errorf("Error in NewAssetWithMultiplier - %s", err)
	}
	if asset.GetTicker() != "ZZB AU" {
		t.Error("Unexpected ticker")
	}
	if asset.GetMultiplier() != 100.0 {
		t.Errorf("Unexpected multiplier")
	}
	if asset.GetPrice().Valid {
		t.Error("Expecting an uninitialised price")
	}
}

func TestGetSetPrice(t *testing.T) {
	asset, err := NewAsset("ZZB AU", "AUD")
	if err != nil {
		t.Errorf("Error in NewAsset - %s", err)
	}
	priceFloat := 2.75
	price := Price{Float64: priceFloat, Valid: true}
	asset.SetPrice(price)

	actualPrice := asset.GetPrice().Float64
	if actualPrice != priceFloat {
		t.Errorf("Unexpected price: wanted %0.2f, got %0.2f", priceFloat, actualPrice)
	}
}

func TestRevalue(t *testing.T) {
	asset, err := NewAssetWithMultiplier("ZZB AU", "AUD", 100.0)
	if err != nil {
		t.Errorf("Error in NewAssetWithMultiplier - %s", err)
	}
	price := Price{Float64: 2.0, Valid: true}
	asset.SetPrice(price)

	value := asset.GetValue()
	if !value.Valid {
		t.Error("We should have a valid value")
	}
	valueFloat := value.Float64
	expectedValueFloat := 200.0
	if valueFloat != expectedValueFloat {
		t.Errorf("Unexpected value: wanted %0.2f, got %0.2f", expectedValueFloat, valueFloat)
	}
}

func TestAssetHistory(t *testing.T) {
	stock, err := NewStock("ZZB AU", "AUD")
	if err != nil {
		t.Errorf("Error in NewStock - %s", err)
	}

	time1 := time.Date(2020, time.December, 14, 0, 0, 0, 0, time.UTC)
	time2 := time.Date(2020, time.December, 15, 0, 0, 0, 0, time.UTC)
	price1 := Price{Float64: 2.5, Valid: true}
	price2 := Price{Float64: 3.0, Valid: true}

	stock.SetPrice(price1)
	stock.TakeSnapshot(time1)
	stock.SetPrice(price2)
	stock.TakeSnapshot(time2)

	history := stock.GetHistory()
	
	// check our first snapshot
	snap1 := history[time1]
	if !snap1.GetTime().Equal(time1) {
		t.Error("snap1 - unexpected time.")
	}
	if snap1.GetPrice().Float64 != 2.5 {
		t.Error("snap1 - unexpected price.")
	}
	if snap1.GetValue().Float64 != 2.5 {
		t.Error("snap1 - unexpected value.")
	}

	// and our second snapshot
	snap2 := history[time2]
	if !snap2.GetTime().Equal(time2) {
		t.Error("snap2 - unexpected time.")
	}
	if snap2.GetPrice().Float64 != 3.0 {
		t.Error("snap2 - unexpected price.")
	}
	if snap2.GetValue().Float64 != 3.0 {
		t.Error("snap2 - unexpected value.")
	}
}
